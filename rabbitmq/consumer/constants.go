package consumer

import (
	"time"
)

var (
	// DefaultTickerInterval is the default interval for the ticker in the consumer
	DefaultTickerInterval = 1 * time.Second

	// DefaultConsumerMessagesChannelBufferSize is the default buffer size for the consume messages channel
	DefaultConsumerMessagesChannelBufferSize = 100
)
