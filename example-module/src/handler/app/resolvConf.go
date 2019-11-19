package main

import (
	"bufio"
	"os"
	"regexp"

	"github.com/pkg/errors"
)

// Because dnsReadConfig in https://golang.org/src/net/dnsconfig_unix.go is private

const defaultResolvConfPath = "/etc/resolv.conf"

var (
	resolvConfNameserverPattern = regexp.MustCompile(`^nameserver +([0-9\.]+)`)
)

// get the dns names servers from the resolv.conf file at the given path
func getDNSNameservers(resolvConfPath string) ([]string, error) {
	if resolvConfPath == "" {
		resolvConfPath = defaultResolvConfPath
	}

	file, err := os.Open(resolvConfPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not open '%s'.", resolvConfPath)
	}
	defer mustClose(file)

	scanner := bufio.NewScanner(file)

	var servers []string
	for scanner.Scan() {
		line := scanner.Text()
		match := resolvConfNameserverPattern.FindStringSubmatch(line)
		if len(match) == 2 {
			servers = append(servers, match[1])
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, errors.Wrapf(err, "Could not read '%s'.", resolvConfPath)
	}

	if len(servers) == 0 {
		return nil, errors.Errorf("No nameservers found in '%s'.", resolvConfPath)
	}

	return servers, nil
}
