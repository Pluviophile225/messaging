package types

// TopicChannel is the data structure for subscriber
type TopicChannel struct {
	// Topic for subscriber to filter on if any
	Topic string
	// Messages is the returned message channel for the subscriber
	Messages chan MessageEnvelope
}

// MessageBusConfig defines the messaging information need to connect to the message bus
// in a publish-subscribe pattern
type MessageBusConfig struct {
	// Broker contains the connection information for publishing and subscribing to the broker for MessageBus
	Broker HostInfo
	// Type indicates the message queue platform being used. e.g. "redis" for Redis Pub/Sub
	Type string
	// Optional contains all other properties of message bus that are specific to
	// certain concrete implementations like MQTT's Qos, for example.
	Optional map[string]string
}
