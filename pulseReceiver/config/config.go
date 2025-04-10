package config

import "time"

var cfg *Config

type (
	Config struct {
		App         `yaml:"app"`
		Server      `yaml:"server"`
		Log         `yaml:"logger"`
		Integration `yaml:"integration"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	Log struct {
		Level        string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
		ReportCaller bool   `yaml:"report-caller" default:"false"`
	}

	Server struct {
		Port         string        `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		ReadTimeout  time.Duration `yaml:"readTimeout" default:"10s"`
		WriteTimeout time.Duration `yaml:"writeTimeout" default:"10s"`
	}

	Integration struct {
		OpenTelemetry `env-required:"true" yaml:"otlp"`
	}

	OpenTelemetry struct {
		URL      string        `env-required:"true" yaml:"url" env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
		Protocol string        `env-required:"true" yaml:"protocol" env:"OTEL_EXPORTER_OTLP_PROTOCOL"`
		Timeout  time.Duration `env-required:"true" yaml:"timeout" env:"OTEL_EXPORTER_OTLP_TIMEOUT"`
	}
)

func ExportConfig(config *Config) {
	cfg = config
}

func GetConfig() *Config {
	return cfg
}
