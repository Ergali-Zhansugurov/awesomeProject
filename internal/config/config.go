package config

import (
	"awesomeProject/internal/logger"
	"github.com/ilyakaznacheev/cleanenv"

	"sync"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug" env-required:"true" `
	Listen  struct {
		Type   string `yaml:"type"  env:"TYPE" env-default:"port"`
		BindIP string `yaml:"bind_ip"  env:"BINDIP" env-default:"localhost"`
		Port   string `yaml:"port"  env:"PORT" env-default:"8080"`
	} `yaml:"listen"`
	MainConfig MainConfig    `yaml:"main_config" `
	Storage    StorageConfig `yaml:"storage" `
}
type MainConfig struct {
	RabbitmqUser      string `yaml:"rabbitmq_user"  env:"RABBITMQ_USER" env-default:"guest"`
	RabbitmqPass      string `yaml:"rabbitmq_pass"  env:"RABBITMQ_PASS" env-default:"guest"`
	RabbitmqIp        string `yaml:"rabbitmq_ip"  env:"RABBITMQ_IP" env-default:"localhost"`
	RabbitmqPort      string `yaml:"rabbitmq_port"   env:"RABBITMQ_PORT" env-default:"5672"`
	Rabbitmqvhost     string `yaml:"rabbitmq_vhost"  env:"RABBITMQ_VHOST" env-default:"/"`
	RabbitmqCertsPath string `yaml:"rabbitmq_certs_path"  env:"RABBITMQ_CERTS_PATH" env-default:"/path/to/certs"`
	RabbitmqCertName  string `yaml:"rabbitmq_cert_name"  env:"RABBITMQ_CERT_NAME" env-default:"cert.pem"`
}
type StorageConfig struct {
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"PORT" env-default:"5432"`
	Database string `yaml:"database" env:"DATABASE" env-default:"postgres"`
	Username string `yaml:"username"  env:"USERNAME" env-default:"postgres"`
	Password string `yaml:"password"  env:"PASSWORD" env-default:"Fencing.666"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logger.GetLogger()
		logger.Info("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
