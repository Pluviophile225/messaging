package types

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testRequestId     = "3ab0e022-464b-4bfe-bf7f-b0154093ddad"
	testCorrelationId = "fa1def22-96de-4d44-8811-00333438c8e3"
	testPayload       = `{"data" : "myData"}`
)

func TestNewMessageEnvelope(t *testing.T) {
	// lint:ignore SA1029 legacy
	// nolint:staticcheck // See golangci-lint #741
	ctx := context.WithValue(context.Background(), CorrelationID, testCorrelationId)
	// lint:ignore SA1029 legacy
	// nolint:staticcheck // See golangci-lint #741
	ctx = context.WithValue(ctx, ContentType, ContentTypeJSON)

	envelope := NewMessageEnvelope([]byte(testPayload), ctx)

	assert.Equal(t, ApiVersion, envelope.ApiVersion)
	assert.Equal(t, testCorrelationId, envelope.CorrelationID)
	assert.Equal(t, ContentTypeJSON, envelope.ContentType)
	assert.Equal(t, testPayload, string(envelope.Payload))
	assert.Empty(t, envelope.QueryParams)
}

func TestNewMessageEnvelopeEmpty(t *testing.T) {
	envelope := NewMessageEnvelope([]byte{}, context.Background())

	assert.Equal(t, ApiVersion, envelope.ApiVersion)
	assert.Empty(t, envelope.RequestID)
	assert.Empty(t, envelope.CorrelationID)
	assert.Empty(t, envelope.ContentType)
	assert.Empty(t, envelope.Payload)
	assert.Zero(t, envelope.ErrorCode)
	assert.Empty(t, envelope.QueryParams)
}

func TestNewMessageEnvelopeForRequest(t *testing.T) {
	expectedQueryParams := map[string]string{"foo": "bar"}
	emptyQueryParams := map[string]string{}

	tests := []struct {
		name          string
		queryParams   map[string]string
		expectedEmpty bool
	}{
		{"valid - normal queryParams map", expectedQueryParams, false},
		{"valid - empty queryParams map", emptyQueryParams, true},
		{"valid - nil queryParams", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envelope := NewMessageEnvelopeForRequest([]byte(testPayload), tt.queryParams)

			assert.NotEmpty(t, envelope.RequestID)
			assert.NotEmpty(t, envelope.CorrelationID)
			assert.Equal(t, ApiVersion, envelope.ApiVersion)
			assert.Equal(t, testPayload, string(envelope.Payload))
			assert.Equal(t, ContentTypeJSON, envelope.ContentType)
			assert.Equal(t, 0, envelope.ErrorCode)
			assert.NotNil(t, envelope.QueryParams)
			if tt.expectedEmpty {
				assert.Empty(t, envelope.QueryParams)
				return
			}

			assert.Equal(t, expectedQueryParams, envelope.QueryParams)
		})
	}
}

func TestNewMessageEnvelopeForResponse(t *testing.T) {
	invalidUUID := "123456"

	tests := []struct {
		name          string
		correlationId string
		requestId     string
		contentType   string
		expectedError bool
	}{
		{"valid", testCorrelationId, testRequestId, ContentTypeJSON, false},
		{"invalid - CorrelationID is not in UUID format", invalidUUID, testRequestId, ContentTypeJSON, true},
		{"invalid - invalid requestID is not in UUID format", testCorrelationId, invalidUUID, ContentTypeJSON, true},
		{"invalid - ContentType is empty", testCorrelationId, testRequestId, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envelope, err := NewMessageEnvelopeForResponse([]byte(testPayload), tt.requestId, tt.correlationId, tt.contentType)
			if tt.expectedError {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, testRequestId, envelope.RequestID)
			assert.Equal(t, testCorrelationId, envelope.CorrelationID)
			assert.Equal(t, ApiVersion, envelope.ApiVersion)
			assert.Equal(t, testPayload, string(envelope.Payload))
			assert.Equal(t, ContentTypeJSON, envelope.ContentType)
			assert.Equal(t, 0, envelope.ErrorCode)
			assert.NotNil(t, envelope.QueryParams)
		})
	}
}

func TestNewMessageEnvelopeFromJSON(t *testing.T) {
	invalidUUID := "123456"
	validEnvelope := testMessageEnvelope()
	validNoCorrelationIDEnvelope := validEnvelope
	validNoCorrelationIDEnvelope.CorrelationID = ""
	invalidApiVersionEnvelope := validEnvelope
	invalidApiVersionEnvelope.ApiVersion = "v2"
	invalidRequestIDEnvelope := validEnvelope
	invalidRequestIDEnvelope.RequestID = invalidUUID
	invalidCorrelationIDEnvelope := validEnvelope
	invalidCorrelationIDEnvelope.CorrelationID = invalidUUID
	invalidContentTypeEnvelope := validEnvelope
	invalidContentTypeEnvelope.ContentType = ""

	tests := []struct {
		name          string
		envelope      MessageEnvelope
		expectedError bool
	}{
		{"valid", validEnvelope, false},
		{"valid - CorrelationID is not set", validNoCorrelationIDEnvelope, false},
		{"invalid - API version not 'v1'", invalidApiVersionEnvelope, true},
		{"invalid - RequestID is not UUID format", invalidRequestIDEnvelope, true},
		{"invalid - CorrelationID is not UUID format", invalidCorrelationIDEnvelope, true},
		{"invalid - ContentType is not application/json", invalidContentTypeEnvelope, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := json.Marshal(tt.envelope)
			require.NoError(t, err)

			envelope, err := NewMessageEnvelopeFromJSON(payload)
			if tt.expectedError {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, testRequestId, envelope.RequestID)
			assert.NotEmpty(t, testCorrelationId, envelope.CorrelationID)
			assert.Equal(t, ApiVersion, envelope.ApiVersion)
			assert.Equal(t, testPayload, string(envelope.Payload))
			assert.Equal(t, ContentTypeJSON, envelope.ContentType)
			assert.Equal(t, 0, envelope.ErrorCode)
			assert.NotNil(t, envelope.QueryParams)
		})
	}
}

func TestNewMessageEnvelopeWithError(t *testing.T) {
	expectedPayload := `error: something failed`

	envelope := NewMessageEnvelopeWithError(testRequestId, expectedPayload)

	assert.NotEmpty(t, testCorrelationId, envelope.CorrelationID)
	assert.Equal(t, ApiVersion, envelope.ApiVersion)
	assert.Equal(t, testRequestId, envelope.RequestID)
	assert.Equal(t, 1, envelope.ErrorCode)
	assert.Equal(t, expectedPayload, string(envelope.Payload))
	assert.Equal(t, ContentTypeText, envelope.ContentType)
	assert.Empty(t, envelope.QueryParams)
}

func testMessageEnvelope() MessageEnvelope {
	return MessageEnvelope{
		CorrelationID: testCorrelationId,
		ApiVersion:    ApiVersion,
		RequestID:     testRequestId,
		ErrorCode:     0,
		Payload:       []byte(testPayload),
		ContentType:   ContentTypeJSON,
	}
}
