package types

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMsgQueueURL(t *testing.T) {

	port := 5570
	publishHost := HostInfo{
		Host:     "*",
		Port:     port,
		Protocol: "tcp",
	}
	subscribeHost := HostInfo{
		Host:     "localhost",
		Port:     port,
		Protocol: "tcp",
	}
	portString := strconv.Itoa(port)

	url := publishHost.GetHostURL()

	if !assert.Equal(t, url, "tcp://*:"+portString, "Failed to create correct URL for publisher") {
		t.Fatal()
	}

	url = subscribeHost.GetHostURL()

	if !assert.Equal(t, url, "tcp://localhost:"+portString, "Failed to create correct URL for subscriber") {
		t.Fatal()
	}
}

func TestIsHostInfoEmpty(t *testing.T) {

	port := 5570
	notEmptyHost := HostInfo{
		Host:     "*",
		Port:     port,
		Protocol: "tcp",
	}

	emptyHost := HostInfo{
		Host:     "",
		Protocol: "",
	}

	if !assert.False(t, notEmptyHost.IsHostInfoEmpty(), "Failed to return expected value") {
		t.Fatal()
	}

	if !assert.True(t, emptyHost.IsHostInfoEmpty(), "Failed to return expected value") {
		t.Fatal()
	}
}
