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
		Kafka         `env-required:"true" yaml:"kafka"`
		Postgres      `env-required:"true" yaml:"postgres"`
	}

	Kafka struct {
		URL        string `env-required:"true" yaml:"url"`
		GroupId    string `env-required:"true" yaml:"groupId"`
		WorkTask   int    `env-required:"true" yaml:"workTask"`
		BatchSize  int    `env-required:"true" yaml:"batchSize"`
		BatchLimit int    `env-required:"true" yaml:"batchLimit"`
		Poll       int    `env-required:"true" yaml:"poll"`
	}

	Postgres struct {
		Host         string `env-required:"true" yaml:"host" env:"HOST_DB"`
		PostgresPort int    `env-required:"true" yaml:"port" env:"PORT_DB"`
		Username     string `env-required:"true" yaml:"username" env:"USERNAME_DB"`
		Password     string `env-required:"true" yaml:"password" env:"PASSWORD_DB"`
		Dbname       string `yaml:"dbname"`
		Schema       string `yaml:"schema"`
		MaxIdleConns int    `yaml:"maxIdleConns"`
		MaxOpenConns int    `yaml:"MaxOpenConns"`
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
