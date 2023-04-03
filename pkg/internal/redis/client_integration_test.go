package redis

import (
	"fmt"
	"messaging/pkg/types"
	"net/url"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	RedisURLEnvName = "REDIS_SERVER_TEST"
	DefaultRedisURL = "redis://localhost:6379"
	TestStream      = "IntegrationTest"
)

func TestRedisClientIntegration(t *testing.T) {
	redisHostInfo := getRedisHostInfo(t)
	client, err := NewClient(types.MessageBusConfig{
		Broker: redisHostInfo,
	})

	require.NoError(t, err, "Failed to create Redis client")
	testMessage := types.MessageEnvelope{
		CorrelationID: "12345",
		Payload:       []byte("test-message"),
		ReceivedTopic: "IntegrationTest",
	}

	testCompleteWaitGroup := &sync.WaitGroup{}
	testCompleteWaitGroup.Add(1)
	createSub(t, &client, TestStream, testMessage, testCompleteWaitGroup)
	// TODO Add proper hook/wait
	time.Sleep(time.Duration(1) * time.Second)
	err = client.Publish(testMessage, TestStream)
	require.NoError(t, err, "Failed to publish test message")
	testCompleteWaitGroup.Wait()

	err = client.Disconnect()
	require.NoError(t, err)
}

// TestRedisUnsubscribeIntegration end-to-end test of subscribing and unsubscribing.
// Redis must be running with Device virtual publishing events.
func TestRedisUnsubscribeIntegration(t *testing.T) {
	redisHostInfo := getRedisHostInfo(t)
	client, err := NewClient(types.MessageBusConfig{
		Broker: redisHostInfo,
	})

	require.NoError(t, err, "Failed to create Redis client")

	messages := make(chan types.MessageEnvelope, 1)
	errs := make(chan error, 1)

	eventTopic := "openyurt/events/#"
	topics := []types.TopicChannel{
		{
			Topic:    eventTopic,
			Messages: messages,
		},
	}

	println("Subscribing to topic: " + eventTopic)
	err = client.Subscribe(topics, errs)
	require.NoError(t, err)
	require.Equal(t, 1, len(client.existingTopics))

	messageCount := 0

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-time.After(time.Second * 10):
				return

			case err = <-errs:
				require.Failf(t, "failed", "Unexpected error message received: %v", err)
				return

			case message := <-messages:
				println(fmt.Sprintf("Received message from topic: %v", message.ReceivedTopic))
				messageCount++
				if messageCount > 3 {
					println("Unsubscribing from topic: " + eventTopic)
					err = client.Unsubscribe(eventTopic)
					require.NoError(t, err)

					time.Sleep(time.Second)
					_, exists := client.existingTopics[eventTopic]
					assert.False(t, exists)
				}
			}
		}

	}()

	wg.Wait()
	assert.Greater(t, messageCount, 3)
	assert.Equal(t, 0, len(client.existingTopics))
}

// TestRedisRequestIntegration depends on Redis and Device Virtual to be running
//func TestRedisRequestIntegration(t *testing.T) {
//	redisHostInfo := getRedisHostInfo(t)
//	client, err := NewClient(types.MessageBusConfig{
//		Broker: redisHostInfo,
//	})
//	require.NoError(t, err, "Failed to create Redis client, Redis must be running")
//
//	dsClient := http.NewCommonClient("http://localhost:59910")
//	_, err = dsClient.Ping(context.Background())
//	require.NoError(t, err, "Device Virtual must be running")
//
//	deviceService := "device-virtual"
//	responseTopicPrefix := commonConstants.BuildTopic(commonConstants.DefaultBaseTopic, commonConstants.ResponseTopic, deviceService)
//	requestTopic := commonConstants.BuildTopic(commonConstants.DefaultBaseTopic, deviceService, commonConstants.ValidateDeviceSubscribeTopic)
//
//	// Device Virtual doesn't implement validation interface, so the SDK will always return ok response
//	expectedErrorCode := 0
//
//	// Send multiple Requests to make sure bug that failed after first successfully request is fixed.
//
//	requestId := uuid.NewString()
//	correlationId := uuid.NewString()
//	fmt.Printf("Sending first request to topic %s with requestId %s\n", requestTopic, requestId)
//	response, err := client.Request(types.MessageEnvelope{RequestID: requestId, CorrelationID: correlationId}, requestTopic, responseTopicPrefix, time.Second*5)
//	require.NoError(t, err)
//	assert.Equal(t, requestId, response.RequestID)
//	assert.Equal(t, expectedErrorCode, response.ErrorCode)
//
//	requestId = uuid.NewString()
//	correlationId = uuid.NewString()
//	fmt.Printf("Sending second request to topic %s with requestId %s\n", requestTopic, requestId)
//	response, err = client.Request(types.MessageEnvelope{RequestID: requestId, CorrelationID: correlationId}, requestTopic, responseTopicPrefix, time.Second*5)
//	require.NoError(t, err)
//	assert.Equal(t, requestId, response.RequestID)
//	assert.Equal(t, expectedErrorCode, response.ErrorCode)
//
//	correlationId = uuid.NewString()
//	requestId = uuid.NewString()
//	fmt.Printf("Sending third request to topic %s with requestId %s\n", requestTopic, requestId)
//	response, err = client.Request(types.MessageEnvelope{RequestID: requestId, CorrelationID: correlationId}, requestTopic, responseTopicPrefix, time.Second*5)
//	require.NoError(t, err)
//	assert.Equal(t, requestId, response.RequestID)
//	assert.Equal(t, expectedErrorCode, response.ErrorCode)
//}

func createSub(t *testing.T, client *Client, stream string, expectedMessage types.MessageEnvelope, doneWaitGroup *sync.WaitGroup) {
	msgError := make(chan error)
	messageChannel := make(chan types.MessageEnvelope)
	err := client.Subscribe([]types.TopicChannel{
		{
			Topic:    stream,
			Messages: messageChannel,
		},
	}, msgError)

	require.NoError(t, err, "Failed to create subscription")
	go func() {
		defer doneWaitGroup.Done()
		select {
		case message := <-msgError:
			assert.Nil(t, message, "Unexpected error message received: %v", message)

		case message := <-messageChannel:
			assert.Equal(t, expectedMessage, message)
		}
	}()
}

func getRedisHostInfo(t *testing.T) types.HostInfo {
	redisURLString := getRedisURL()
	redisURL, err := url.Parse(redisURLString)
	if err != nil {
		require.NoError(t, err, "Unable to parse the Redis URL %s: %v", redisURL, err)
		t.FailNow()
	}

	port, err := strconv.Atoi(redisURL.Port())
	if err != nil {
		require.NoError(t, err)
		t.FailNow()
	}

	host := redisURL.Hostname()
	scheme := redisURL.Scheme

	return types.HostInfo{
		Host:     host,
		Port:     port,
		Protocol: scheme,
	}
}

func getRedisURL() string {
	redisURLString := os.Getenv(RedisURLEnvName)
	if redisURLString == "" {
		return DefaultRedisURL
	}

	return redisURLString
}
