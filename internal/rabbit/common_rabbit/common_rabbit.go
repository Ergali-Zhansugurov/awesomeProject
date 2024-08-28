package common_rabbit

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/logger"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wagslane/go-rabbitmq"
	"io/ioutil"
)

const (
	MQ_URL     = "amqp://%v:%v@%v:%v/%v"
	MQ_URL_SSL = "amqps://%v:%v@%v:%v/%v"
)

func MqTlsConfig(certPath string) (*tls.Config, error) {
	caCert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(caCert)

	return &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            rootCAs,
	}, nil
}

func InitConsumer(cfg config.MainConfig, q string) (consumer rabbitmq.Consumer, err error) {
	url := fmt.Sprintf(MQ_URL, cfg.RabbitmqUser, cfg.RabbitmqPass, cfg.RabbitmqIp, cfg.RabbitmqPort, cfg.Rabbitmqvhost)
	return rabbitmq.NewConsumer(url, rabbitmq.Config{},
		rabbitmq.WithConsumerOptionsLogging,
	)
}

func InitProducer(cfg config.MainConfig) (publisher *rabbitmq.Publisher, err error) {
	url := fmt.Sprintf(MQ_URL, cfg.RabbitmqUser, cfg.RabbitmqPass, cfg.RabbitmqIp, cfg.RabbitmqPort, cfg.Rabbitmqvhost)
	return rabbitmq.NewPublisher(url, rabbitmq.Config{})
}

func CreateQueue(cfg config.MainConfig, queueName string) (err error) {
	url := fmt.Sprintf(MQ_URL, cfg.RabbitmqUser, cfg.RabbitmqPass, cfg.RabbitmqIp, cfg.RabbitmqPort, cfg.Rabbitmqvhost)

	connectRabbitMQ, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	defer connectRabbitMQ.Close()

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		return err
	}
	defer channelRabbitMQ.Close()

	_, err = channelRabbitMQ.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // auto delete
		false,     // exclusive
		false,     // no wait
		nil,       // arguments
	)
	return err
}

func DeleteQueue(cfg config.MainConfig, queueName string) (err error) {
	url := fmt.Sprintf(MQ_URL, cfg.RabbitmqUser, cfg.RabbitmqPass, cfg.RabbitmqIp, cfg.RabbitmqPort, cfg.Rabbitmqvhost)

	connectRabbitMQ, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	defer connectRabbitMQ.Close()

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		return err
	}
	defer channelRabbitMQ.Close()

	msgCnt, err := channelRabbitMQ.QueueDelete(queueName, false, false, false)
	if err != nil {
		return err
	}
	logger.GetLogger().Infof("delete rabbit queue", "rmq", fmt.Sprintf("in queue %s, %d message was purges.", queueName, msgCnt), context.Background())
	return err
}
