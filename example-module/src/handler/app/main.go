package main // import "section.io/module-handler"

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	openRestyBinaryPath = "/usr/bin/openresty"
	logPath             = "/var/log/nginx"
	configOutputPath    = "/var/openresty/"
)

var validLogLevels = [...]string{"emerg", "alert", "crit", "error", "warn",
	"notice", "info", "debug"}

const defaultLogLevel = "warn"

func ensureLogDirectory() {
	mustMkdirAll(logPath, 0755)
}

// reads the value of Ops Variables used by Section to team to fine tune module behaviour
func getOpsVariable(name string) string {

	const opsVariablesFile = "/.section.io/ops-variables"

	fullName := "SECTION_OPS_" + strings.ToUpper(name)

	data, err := ioutil.ReadFile(opsVariablesFile)
	if err != nil {
		// DEBUG log.Printf("Failed to read file '%s': %v", opsVariablesFile, err)
	} else {
		content := string(data)
		for _, line := range strings.Split(content, "\n") {
			if strings.HasPrefix(line, fullName+"=") {
				return line[len(fullName)+1:]
			}
		}
	}

	return ""
}

// creates the error_log.conf file used to the set the log level for the nginx error logs
// error_log.conf is included in the nginx.conf : https://github.com/section-io/module-build-assets/blob/50dee498b0f982a50fcfbaf48ada1ac233a4d509/example-module/src/proxy/nginx.conf#L20
func updateErrorLogConf() {
	confPath := path.Join(configOutputPath, "error_log.conf")

	logLevelOpsVariable := getOpsVariable("log_level")

	confLogLevel := defaultLogLevel
	for _, level := range validLogLevels {
		if level == logLevelOpsVariable {
			confLogLevel = level
			break
		}
	}

	if confLogLevel != defaultLogLevel {
		log.Printf("Using non-default error_log level: %s\n", confLogLevel)
	}

	confContent := fmt.Sprintf("error_log %s %s;\n",
		path.Join(logPath, "error.log"),
		confLogLevel)

	mustWriteFile(confPath, []byte(confContent), 0644)
}

// generates the resolver.conf file used by nginx as nginx does not read from /etc/resolver.conf
func generateNginxResolver() {
	confPath := path.Join(configOutputPath, "resolver.conf")

	nameservers, err := getDNSNameservers("")
	if err != nil {
		log.Panicf("Failed to get DNS nameservers: %#v", err)
	}

	// turn off looking up of ipv6 addresses for the resolver
	// http://nginx.org/en/docs/http/ngx_http_core_module.html#resolver
	confContent := fmt.Sprintf("resolver %s ipv6=off;\n", strings.Join(nameservers, " "))
	mustWriteFile(confPath, []byte(confContent), 0644)
}

// prepares openresty configuration
func prepareConfiguration() {
	mustMkdirAll(configOutputPath, 0755)

	updateErrorLogConf()
	generateNginxResolver()
	generateLuaFile()
}

func testOpenRestyConfiguration() error {
	command := exec.Command(openRestyBinaryPath, "-tq")

	_, err := command.Output()
	return err
}

func reloadOpenRestyConfiguration() error {
	command := exec.Command(openRestyBinaryPath, "-s", "reload")
	return command.Run()
}

func runValidate() {
	ensureLogDirectory()
	prepareConfiguration()

	if err := testOpenRestyConfiguration(); err != nil {
		writeStderrAndExit(err)
	}
}

func main() {
	defaultMode := path.Base(os.Args[0])

	mode := flag.String("mode", defaultMode, "Specify the mode")
	flag.Parse()
	// provides the entrypoints for the module
	// validate.sh - runs during git push to do module validation
	// default - starts the module on Section infrastructure
	switch *mode {
	case "validate.sh":
		runValidate()
	default:
		runKubernetesEntrypoint()
	}
}
