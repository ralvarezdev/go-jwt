package server

import (
	"errors"
)

var (
	TokenHasExpiredError = errors.New("token has expired")
)
