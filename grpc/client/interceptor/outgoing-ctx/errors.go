package outgoing_ctx

import (
	"errors"
)

var (
	FailedToGetOutgoingContextError = errors.New(
		"failed to get outgoing context",
	)
)
