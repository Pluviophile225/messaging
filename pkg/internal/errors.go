package internal

import "fmt"

// CertificateErr represents an error associated with interacting with a Certificate.
type CertificateErr struct {
	description string
}

func (ce CertificateErr) Error() string {
	return fmt.Sprintf("Unable to process certificate properties: %s", ce.description)
}

// NewCertificateErr constructs a new CertificateErr
func NewCertificateErr(message string) CertificateErr {
	return CertificateErr{description: message}
}

// BrokerURLErr represents an error associated parsing a broker's URL.
type BrokerURLErr struct {
	description string
}

func (bue BrokerURLErr) Error() string {
	return fmt.Sprintf("Unable to process broker URL: %s", bue.description)
}

// NewBrokerURLErr constructs a new BrokerURLErr
func NewBrokerURLErr(description string) BrokerURLErr {
	return BrokerURLErr{description: description}
}

type PublishHostURLErr struct {
	description string
}

func (p PublishHostURLErr) Error() string {
	return fmt.Sprintf("Unable to use PublishHost URL: %s", p.description)
}

func NewPublishHostURLErr(message string) PublishHostURLErr {
	return PublishHostURLErr{description: message}
}

type SubscribeHostURLErr struct {
	description string
}

func (p SubscribeHostURLErr) Error() string {
	return fmt.Sprintf("Unable to use SubscribeHost URL: %s", p.description)
}

func NewSubscribeHostURLErr(message string) SubscribeHostURLErr {
	return SubscribeHostURLErr{description: message}
}

type MissingConfigurationErr struct {
	missingConfiguration string
	description          string
}

func (mce MissingConfigurationErr) Error() string {
	return fmt.Sprintf("Missing configuration '%s' : %s", mce.missingConfiguration, mce.description)
}

func NewMissingConfigurationErr(missingConfiguration string, message string) MissingConfigurationErr {
	return MissingConfigurationErr{
		missingConfiguration: missingConfiguration,
		description:          message,
	}
}

type InvalidTopicErr struct {
	topic       string
	description string
}

func (ite InvalidTopicErr) Error() string {
	return fmt.Sprintf("Invalid topic '%s': %s", ite.topic, ite.description)
}

func NewInvalidTopicErr(topic string, description string) InvalidTopicErr {
	return InvalidTopicErr{
		topic:       topic,
		description: description,
	}
}
