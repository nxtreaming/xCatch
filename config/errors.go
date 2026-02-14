package config

import "errors"

var (
	ErrMissingAPIKey = errors.New("config: XCATCH_API_KEY is required")
)
