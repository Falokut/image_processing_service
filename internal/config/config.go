package config

import (
	"sync"

	"github.com/Falokut/image_processing_service/pkg/jaeger"
	"github.com/Falokut/image_processing_service/pkg/metrics"
	logging "github.com/Falokut/online_cinema_ticket_office.loggerwrapper"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel        string `yaml:"log_level" env:"LOG_LEVEL"`
	HealthcheckPort string `yaml:"healthcheck_port"`
	EnableMetrics   bool   `yaml:"enable_metrics" env:"ENABLE_METRICS"`
	Listen          struct {
		Host            string `yaml:"host" env:"HOST"`
		Port            string `yaml:"port" env:"PORT"`
		Mode            string `yaml:"server_mode" env:"SERVER_MODE"` // support GRPC, REST, BOTH
		MaxRequestSize  int    `yaml:"max_request_size" env:"MAX_REQUEST_SIZE"`
		MaxResponseSize int    `yaml:"max_response_size" env:"MAX_RESPONSE_SIZE"`
	} `yaml:"listen"`

	PrometheusConfig struct {
		Name         string                      `yaml:"service_name" env:"PROMETHEUS_SERVICE_NAME"`
		ServerConfig metrics.MetricsServerConfig `yaml:"server_config"`
	} `yaml:"prometheus"`

	JaegerConfig jaeger.Config `yaml:"jaeger"`
}

const configsPath string = "configs/"

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		instance = &Config{}

		if err := cleanenv.ReadConfig(configsPath+"config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Fatal(help, " ", err)
		}
	})
	return instance
}
