package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultBaseURL    = "https://fapi.uk"
	AltBaseURL        = "https://l2.fapi.uk"
	DefaultTimeout    = 30 * time.Second
	DefaultMaxRetries = 3
	DefaultRateLimit  = 5.0 // QPS
)

// Config holds the configuration for the uTools API client.
type Config struct {
	// BaseURL is the API base URL. Default: https://fapi.uk
	BaseURL string

	// APIKey is the uTools API key for authentication.
	APIKey string

	// AuthToken is the Twitter auth_token, required by some endpoints
	// (e.g. HomeTimeline, Notifications).
	AuthToken string

	// Timeout is the HTTP request timeout.
	Timeout time.Duration

	// MaxRetries is the maximum number of retries on rate limit / transient errors.
	MaxRetries int

	// RateLimit is the maximum requests per second (QPS).
	RateLimit float64
}

// LoadFromFile creates a Config by reading a config.ini file.
// The INI file format supports [xcatch] section with keys:
//
//	api_key, auth_token, base_url, timeout_sec, max_retries, rate_limit
func LoadFromFile(path string) (*Config, error) {
	kvs, err := parseINI(path, "xcatch")
	if err != nil {
		return nil, fmt.Errorf("config: load %s: %w", path, err)
	}

	cfg := &Config{
		BaseURL:    DefaultBaseURL,
		Timeout:    DefaultTimeout,
		MaxRetries: DefaultMaxRetries,
		RateLimit:  DefaultRateLimit,
	}

	if v, ok := kvs["api_key"]; ok {
		cfg.APIKey = v
	}
	if v, ok := kvs["auth_token"]; ok {
		cfg.AuthToken = v
	}
	if v, ok := kvs["base_url"]; ok && v != "" {
		cfg.BaseURL = v
	}
	if v, ok := kvs["timeout_sec"]; ok {
		if sec, err := strconv.Atoi(v); err == nil && sec > 0 {
			cfg.Timeout = time.Duration(sec) * time.Second
		}
	}
	if v, ok := kvs["max_retries"]; ok {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			cfg.MaxRetries = n
		}
	}
	if v, ok := kvs["rate_limit"]; ok {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f > 0 {
			cfg.RateLimit = f
		}
	}

	return cfg, nil
}

// Load loads configuration with the following priority (highest to lowest):
//  1. config.ini file (if it exists at the given path)
//  2. Environment variables (as fallback for any missing fields)
//
// If path is empty, it defaults to "config.ini" in the current directory.
func Load(path string) *Config {
	if path == "" {
		path = "config.ini"
	}

	// Try loading from file first
	cfg, err := LoadFromFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("config warning: failed to parse %s, falling back to defaults/env: %v", path, err)
		}
		// File not found or parse error, start from defaults
		cfg = &Config{
			BaseURL:    DefaultBaseURL,
			Timeout:    DefaultTimeout,
			MaxRetries: DefaultMaxRetries,
			RateLimit:  DefaultRateLimit,
		}
	}

	// Environment variables override file values (if set)
	if v := os.Getenv("XCATCH_API_KEY"); v != "" {
		cfg.APIKey = v
	}
	if v := os.Getenv("XCATCH_AUTH_TOKEN"); v != "" {
		cfg.AuthToken = v
	}
	if v := os.Getenv("XCATCH_BASE_URL"); v != "" {
		cfg.BaseURL = v
	}
	if v := os.Getenv("XCATCH_TIMEOUT_SEC"); v != "" {
		if sec, err := strconv.Atoi(v); err == nil && sec > 0 {
			cfg.Timeout = time.Duration(sec) * time.Second
		}
	}
	if v := os.Getenv("XCATCH_MAX_RETRIES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			cfg.MaxRetries = n
		}
	}
	if v := os.Getenv("XCATCH_RATE_LIMIT"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f > 0 {
			cfg.RateLimit = f
		}
	}

	return cfg
}

// parseINI reads an INI file and returns key-value pairs for the given section.
// If section is empty, it reads keys before any section header.
func parseINI(path, section string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make(map[string]string)
	currentSection := ""
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		// Section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.TrimSpace(line[1 : len(line)-1])
			continue
		}

		// Key = Value
		if currentSection == section {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				result[key] = val
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// Validate checks that required fields are set.
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return ErrMissingAPIKey
	}
	if c.BaseURL == "" {
		c.BaseURL = DefaultBaseURL
	}
	if c.Timeout <= 0 {
		c.Timeout = DefaultTimeout
	}
	if c.MaxRetries < 0 {
		c.MaxRetries = DefaultMaxRetries
	}
	if c.RateLimit <= 0 {
		c.RateLimit = DefaultRateLimit
	}
	return nil
}
