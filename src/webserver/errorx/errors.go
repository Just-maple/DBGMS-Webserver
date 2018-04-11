package errorx

import (
	"errors"
)

var ErrAuthFailed = errors.New("auth failed")
var ErrMethodInvalid = errors.New("invalid method")
var ErrIdInvalid = errors.New("invalid Id")
