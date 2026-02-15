//go:build integration

package utools

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xCatch/xcatch/config"
)

func discoverConfigPath() (string, error) {
	if p := os.Getenv("XCATCH_CONFIG_PATH"); p != "" {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
		return "", fmt.Errorf("XCATCH_CONFIG_PATH not found: %s", p)
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	candidates := []string{
		filepath.Join(wd, "config.ini"),
		filepath.Join(wd, "..", "config.ini"),
		filepath.Join(wd, "..", "..", "config.ini"),
	}

	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("config.ini not found from working dir %s", wd)
}

func requireIntegrationClient(t *testing.T) *Client {
	t.Helper()
	if os.Getenv("XCATCH_RUN_INTEGRATION") != "1" {
		t.Skip("set XCATCH_RUN_INTEGRATION=1 to run real API integration tests")
	}

	configPath, err := discoverConfigPath()
	if err != nil {
		t.Skipf("integration config file not found: %v", err)
	}

	cfg := config.Load(configPath)
	if cfg.APIKey == "" {
		t.Skipf("missing XCATCH_API_KEY (or api_key/xcatch_api_key in %s)", configPath)
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("new integration client: %v", err)
	}
	return client
}

func requireNonEmptyJSON(t *testing.T, method string, raw json.RawMessage) {
	t.Helper()
	if len(raw) == 0 {
		t.Fatalf("%s returned empty payload", method)
	}
	if !json.Valid(raw) {
		t.Fatalf("%s returned invalid json: %s", method, string(raw))
	}
	t.Logf("%s sample response: %s", method, Truncate(string(raw), 240))
}

func lookupConfigValue(path, section, key string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	targetSection := strings.ToLower(strings.TrimSpace(section))
	targetKey := strings.ToLower(strings.TrimSpace(key))
	currentSection := ""

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.ToLower(strings.TrimSpace(line[1 : len(line)-1]))
			continue
		}

		if currentSection != targetSection {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		k := strings.ToLower(strings.TrimSpace(parts[0]))
		if k == targetKey {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}

func integrationTestValue(t *testing.T, envKey string) string {
	t.Helper()
	if v := os.Getenv(envKey); v != "" {
		return v
	}

	configPath, err := discoverConfigPath()
	if err != nil {
		t.Logf("config path discovery failed for %s: %v", envKey, err)
		return ""
	}

	v, err := lookupConfigValue(configPath, "xcatch", envKey)
	if err != nil {
		t.Logf("config lookup failed for %s: %v", envKey, err)
		return ""
	}
	return v
}

func requireIntegrationJSON(t *testing.T, name string, call func() (json.RawMessage, error)) {
	t.Helper()
	raw, err := call()
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode >= 500 {
			t.Skipf("%s skipped due upstream %d: %s", name, apiErr.StatusCode, Truncate(apiErr.Message, 200))
		}
		t.Fatalf("%s error: %v", name, err)
	}
	requireNonEmptyJSON(t, name, raw)
}

func TestUserIntegration_RealAPI(t *testing.T) {
	client := requireIntegrationClient(t)
	userID := integrationTestValue(t, "XCATCH_TEST_USER_ID")
	if userID == "" {
		t.Skip("missing XCATCH_TEST_USER_ID (env or config.ini [xcatch])")
	}
	screenName := integrationTestValue(t, "XCATCH_TEST_SCREEN_NAME")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("GetUserByID", func(t *testing.T) {
		requireIntegrationJSON(t, "GetUserByID", func() (json.RawMessage, error) {
			return client.GetUserByID(ctx, userID)
		})
	})

	t.Run("GetUsersByIDs", func(t *testing.T) {
		requireIntegrationJSON(t, "GetUsersByIDs", func() (json.RawMessage, error) {
			return client.GetUsersByIDs(ctx, []string{userID})
		})
	})

	t.Run("GetUsernameChanges", func(t *testing.T) {
		requireIntegrationJSON(t, "GetUsernameChanges", func() (json.RawMessage, error) {
			return client.GetUsernameChanges(ctx, userID)
		})
	})

	t.Run("LookupUser by userID", func(t *testing.T) {
		requireIntegrationJSON(t, "LookupUser(by userID)", func() (json.RawMessage, error) {
			return client.LookupUser(ctx, "", userID)
		})
	})

	t.Run("GetUserByIDV2", func(t *testing.T) {
		requireIntegrationJSON(t, "GetUserByIDV2", func() (json.RawMessage, error) {
			return client.GetUserByIDV2(ctx, userID)
		})
	})

	t.Run("GetUsersByIDsV2", func(t *testing.T) {
		requireIntegrationJSON(t, "GetUsersByIDsV2", func() (json.RawMessage, error) {
			return client.GetUsersByIDsV2(ctx, []string{userID})
		})
	})

	t.Run("GetUserByScreenName", func(t *testing.T) {
		if screenName == "" {
			t.Skip("missing XCATCH_TEST_SCREEN_NAME (env or config.ini [xcatch])")
		}
		requireIntegrationJSON(t, "GetUserByScreenName", func() (json.RawMessage, error) {
			return client.GetUserByScreenName(ctx, screenName)
		})
	})

	t.Run("LookupUser by screenName", func(t *testing.T) {
		if screenName == "" {
			t.Skip("missing XCATCH_TEST_SCREEN_NAME (env or config.ini [xcatch])")
		}
		requireIntegrationJSON(t, "LookupUser(by screenName)", func() (json.RawMessage, error) {
			return client.LookupUser(ctx, screenName, "")
		})
	})

	t.Run("GetUserByScreenNameV2", func(t *testing.T) {
		if screenName == "" {
			t.Skip("missing XCATCH_TEST_SCREEN_NAME (env or config.ini [xcatch])")
		}
		requireIntegrationJSON(t, "GetUserByScreenNameV2", func() (json.RawMessage, error) {
			return client.GetUserByScreenNameV2(ctx, screenName)
		})
	})

	t.Run("GetAccountAnalytics", func(t *testing.T) {
		if client.authToken == "" {
			t.Skip("missing auth token; set XCATCH_AUTH_TOKEN or auth_token in config.ini")
		}
		requireIntegrationJSON(t, "GetAccountAnalytics", func() (json.RawMessage, error) {
			return client.GetAccountAnalytics(ctx)
		})
	})
}
