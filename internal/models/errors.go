package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: no matching record found")
var ErrIsExpired = errors.New("models: record has expired")
