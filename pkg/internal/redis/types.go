package redis

import (
	"crypto/tls"
	"messaging/pkg/types"
)

// RedisClientCreator type alias for functions which create RedisClient implementation.
//
// This is mostly used for testing purposes so that we can easily inject mocks.
type RedisClientCreator func(redisServerURL string, password string, tlsConfig *tls.Config) (RedisClient, error)

// RedisClient provides functionality needed to read and send messages to/from Redis' Redis Pub/Sub functionality.
//
// The main reason for this interface is to abstract out the underlying client from Client so that it can be mocked and
// allow for easy unit testing. Since 'go-redis' does not leverage interfaces and has complicated entities it can become
// complex to test the operations without requiring a running Redis server.
type RedisClient interface {
	// Subscribe creates the subscription in Redis
	Subscribe(topic string)
	// Unsubscribe closes the subscription in Redis and removes it.
	Unsubscribe(topic string)
	// Send sends a message to the specified topic, aka Publish.
	Send(topic string, message types.MessageEnvelope) error
	// Receive blocking operation which receives the next message for the specified subscribed topic
	// This supports multi-level topic scheme with wild cards
	Receive(topic string) (*types.MessageEnvelope, error)
	// Close cleans up any entities which need to be deconstructed.
	Close() error
}
