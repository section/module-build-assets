package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"syscall"

	metrics "github.com/section-io/module-metrics"
	"github.com/section-io/proxymessage"
)

var (
	proxyMessageClient *proxymessage.Client
)

type jsonBusMessage struct {
	MessageType string `json:"message_type"`
}

// redirect openresty access logs to stdout and error logs to stderr
// export module metrics
func ensureOpenRestyLogsToStdio() {
	ensureLogDirectory()

	log.Printf("Redirecting access.log to STDOUT and error.log to STDERR.\n")

	pid := os.Getpid()

	fifoFilePath := path.Join(logPath, "access.log")
	// sets up the Section metrics module to export metrics to prometheus
	err := metrics.SetupModule(fifoFilePath, os.Stdout, os.Stderr)
	if err != nil {
		panic(err)
	}

	mustSymlink(fmt.Sprintf("/proc/%d/fd/%d", pid, syscall.Stderr),
		path.Join(logPath, "error.log"))
}

// update the module configuration
// reads and updates the configuration from the user's Section configuration, updates error_log.conf, resolver.conf and reloads nginx
func handleUpdateEnvironment(message *jsonBusMessage) {
	prepareConfiguration()
	if err := reloadOpenRestyConfiguration(); err != nil {
		log.Printf("OpenResty configuration reload error: %#v\n", err)
	}
}

// handle incoming updateenvironment messages
func handleNextMessage() {
	encodedMessage := <-proxyMessageClient.InboundMessageChannel

	var message jsonBusMessage
	err := json.Unmarshal([]byte(encodedMessage), &message)
	if err == nil {
		// check if the message type is correct and handle it
		if message.MessageType == "updateenvironment" {
			handleUpdateEnvironment(&message)
			return
		}

		log.Printf("unsupported message type '%s'", message.MessageType)
	} else {
		log.Printf("json unmarshal error: %v", err)
	}

	log.Printf("message could not be processed: %v", encodedMessage)
}

// starts the nginx server
func startOpenResty() {
	cmd := exec.Command(openRestyBinaryPath, "-g", "daemon off;")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Panicf("OpenResty error: %#v", err)
	}
}

// default entrypoint for the module
func runKubernetesEntrypoint() {

	// set up the redis client to listen to incoming udpateEnvironment messages
	proxyMessageClient = proxymessage.NewClientFromEnvVars()

	// set up openresty configuration and start openresty
	ensureOpenRestyLogsToStdio()
	prepareConfiguration()
	startOpenResty()

	log.Printf("Handling incoming messages on key '%s'.", proxyMessageClient.GetListKey())

	// starts handling updateEnvironment messages
	for {
		handleNextMessage()
	}
}
