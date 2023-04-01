package redis

import (
	"messaging/pkg/internal"
	"messaging/pkg/types"
)

// OptionalClientConfiguration contains additional configuration properties which can be provided via the
// MessageBus.Optional's field.
type OptionalClientConfiguration struct {
	Password string
}

// NewClientConfiguration creates a OptionalClientConfiguration based on the configuration properties provided.
func NewClientConfiguration(config types.MessageBusConfig) (OptionalClientConfiguration, error) {
	redisConfig := OptionalClientConfiguration{}
	err := internal.Load(config.Optional, &redisConfig)
	if err != nil {
		return OptionalClientConfiguration{}, err
	}
	return redisConfig, nil
}
