package config

import "context"

type Config struct {
	Model       string `json:"model"`        // The Ollama model to use.
	Repo        string `json:"repo"`         // Path to the repository to review.
	Endpoint    string `json:"endpoint"`     // API endpoint for remote Ollama instance, includes Port.
	Debug       bool   `json:"debug"`        // Enable debug output.
	Timeout     int    `json:"timeout"`      // Request timeout in seconds (default: 30).
	AccessToken string `json:"access_token"` // Access token for Ollama authentication.
	NoColor     bool   `json:"no_color"`     // Disable color output.
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
