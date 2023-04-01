package redis

import (
	"fmt"
	"strings"
)

// DisconnectErr represents errors which occur when attempting to disconnect from a Redis server.
type DisconnectErr struct {
	// disconnectErrors contains the descriptive error messages that occur while attempting to disconnect one or more
	// underlying clients.
	disconnectErrors []string
}

// Error constructs an appropriate error message based on the error descriptions provided.
func (d DisconnectErr) Error() string {
	return fmt.Sprintf("Unable to disconnect client(s): %s", strings.Join(d.disconnectErrors, ","))
}

// NewDisconnectErr created a new DisconnectErr
func NewDisconnectErr(disconnectErrors []string) DisconnectErr {
	return DisconnectErr{
		disconnectErrors: disconnectErrors,
	}
}
