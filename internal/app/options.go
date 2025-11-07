package app

import (
	"log"

	"github.com/mygallery/mygallery/internal/config"
)

type options struct {
	configPath     string
	configOverride *config.Config
	logger         *log.Logger
}

// Option configures application bootstrap behaviour.
type Option func(*options)

func defaultOptions() options {
	return options{
		configPath: "config.yaml",
		logger:     log.Default(),
	}
}

// WithConfigPath overrides the default config file path.
func WithConfigPath(path string) Option {
	return func(o *options) {
		o.configPath = path
	}
}

// WithLogger sets a custom logger.
func WithLogger(logger *log.Logger) Option {
	if logger == nil {
		return func(o *options) {}
	}
	return func(o *options) {
		o.logger = logger
	}
}

// WithConfig injects a pre-loaded configuration instance.
func WithConfig(cfg *config.Config) Option {
	return func(o *options) {
		o.configOverride = cfg
	}
}
