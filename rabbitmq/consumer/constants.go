package consumer

import (
	"time"
)

var (
	// DefaultTickerInterval is the default interval for the ticker in the consumer
	DefaultTickerInterval = 1 * time.Second

	// DefaultTokensMessageConsumerChannelBufferSize is the default buffer size for the tokens message consumer channel
	DefaultTokensMessageConsumerChannelBufferSize = 100
)
