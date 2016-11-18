// Package config retrieves the service configuration through environment variables.
// The package is included using the infamous "dot import" to spare us the stuttering.
package config

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/hotolab/envconfig"
)

// Config is initialized at runtime and contains the whole application configuration.
var Config cfg

type cfg struct {
	// HTTP
	HTTPBind string `envconfig:"HTTP_BIND" default:"127.0.0.1"`
	HTTPPort int    `envconfig:"HTTP_PORT" default:"8080"`

	// Logging
	LogLevel string `envconfig:"LOG_LEVEL" default:"debug"`

	// Cors
	CorsOrigins []string `envconfig:"CORS_ORIGINS" default:"*"`
	CorsHeaders []string `envconfig:"CORS_HEADERS" default:"Origin,X-Requested-With,Content-Type,Accept,Accept-Language,Authorization"`
	CorsMethods []string `envconfig:"CORS_METHODS" default:"GET,OPTIONS"`

	// Cache
	CacheDuration time.Duration `envconfig:"CACHE_DURATION" default:"2m"`
}

// InitializeConfig loads the configuration from the environment.
func InitializeConfig() {
	err := envconfig.Process("", &Config)
	if err != nil {
		log.Fatal(err)
	}
}
