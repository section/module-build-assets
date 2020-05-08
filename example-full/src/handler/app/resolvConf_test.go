package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDnsNameservers(t *testing.T) {
	nameservers, err := getDNSNameservers("testdata/resolv.conf")
	assert.Nil(t, err, "getDnsNameservers")

	assert.Equal(t, 3, len(nameservers), "len(nameservers)")

	assert.Equal(t, "110.232.116.126", nameservers[0], "nameservers[0]")
	assert.Equal(t, "110.232.119.126", nameservers[1], "nameservers[1]")
	assert.Equal(t, "10.19.240.10", nameservers[2], "nameservers[2]")
}
