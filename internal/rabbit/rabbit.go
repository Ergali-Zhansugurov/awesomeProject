package rabbit

import (
	configs "awesomeProject/internal/config"
	"awesomeProject/internal/logger"
	"awesomeProject/internal/rabbit/common_rabbit"
	"fmt"
	"github.com/wagslane/go-rabbitmq"
	"gopkg.in/yaml.v3"
)

type RabbitMQ struct {
	Logger   *logger.Logger
	Producer *rabbitmq.Publisher
	Consumer rabbitmq.Consumer
}

func InitProducer() (*rabbitmq.Publisher, error) {
	cfg := configs.GetConfig()
	//if cfg.AppProfile == "prod" {
	//  return commonRabbit.InitTlsProducer(cfg.MainConfig)
	//}
	return common_rabbit.InitProducer(cfg.MainConfig)
}

func InitConsumer(q string) (rabbitmq.Consumer, error) {
	cfg := configs.GetConfig()
	//if cfg.AppProfile == "prod" {
	//  return commonRabbit.InitTlsConsumer(cfg.MainConfig)
	//}
	return common_rabbit.InitConsumer(cfg.MainConfig, q)
}

func NewRabbitMQ(l *logger.Logger, p *rabbitmq.Publisher, c rabbitmq.Consumer) *RabbitMQ {
	return &RabbitMQ{Logger: l, Producer: p, Consumer: c}
}

func (r *RabbitMQ) PushMessageToQueue(message interface{}, queueName string, headers rabbitmq.Table) error {
	bytes, err := yaml.Marshal(message)
	if err != nil {
		r.Logger.Errorf("PushMessageToQueue", "RabbitMQ", fmt.Sprintf("serialization error. %s", err.Error()), "widget-library-service")
		return nil
	}
	return r.Producer.Publish(bytes, []string{queueName}, func(options *rabbitmq.PublishOptions) {
		options.Headers = headers
	})
}
