package config

import "context"

type Config struct {
	Model       string `json:"model,omitempty"`        // The Ollama model to use.
	Repo        string `json:"repo,omitempty"`         // Path to the repository to review.
	Endpoint    string `json:"endpoint,omitempty"`     // API endpoint for remote Ollama instance, includes Port.
	Debug       bool   `json:"debug,omitempty"`        // Enable debug output.
	Timeout     int    `json:"timeout,omitempty"`      // Request timeout in seconds (default: 30).
	AccessToken string `json:"access_token,omitempty"` // Access token for Ollama authentication.
	NoColor     bool   `json:"no_color,omitempty"`     // Disable color output.
	Dry         bool   `json:"dry,omitempty"`          // Does a dry-run, e.g., skips the creation of new data
}

type key string

var (
	configKey key = "prollama_config"
)

func FromContext(ctx context.Context) (*Config, bool) {
	c, ok := ctx.Value(configKey).(*Config)
	return c, ok
}

func NewContext(ctx context.Context, c *Config) context.Context {
	return context.WithValue(ctx, configKey, c)
}
