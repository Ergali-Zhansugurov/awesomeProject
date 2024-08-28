package broker

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/logger"
	"awesomeProject/internal/message_broker/broker_models"
	"awesomeProject/internal/models"
	"awesomeProject/internal/rabbit/common_rabbit"
	"context"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"github.com/rabbitmq/amqp091-go"
	"github.com/wagslane/go-rabbitmq"
	"gopkg.in/yaml.v3"
	"log"
)

const cacheTopic = "cache"

type (
	RabbitBroker struct {
		Logger        *logger.Logger
		conn          *amqp091.Connection
		channel       *amqp091.Channel
		cache         *lru.TwoQueueCache
		clientID      string
		Consumer      rabbitmq.Consumer
		consumeHandle *cacheConsumeHandler
		producer      *rabbitmq.Publisher
		queue         *amqp091.Queue
	}

	cacheConsumeHandler struct {
		cache *lru.TwoQueueCache
	}
)

func NewRabbitBroker(cache *lru.TwoQueueCache, clientID string) broker_models.CacheBroker {
	return &RabbitBroker{
		cache:         cache,
		clientID:      clientID,
		consumeHandle: &cacheConsumeHandler{cache: cache},
	}
}

func (c *RabbitBroker) Connect(ctx context.Context, cfg config.MainConfig) error {
	common_rabbit.CreateQueue(cfg, "q")

	producer, err := common_rabbit.InitProducer(cfg)
	if err != nil {
		c.Logger.Errorf("error", err)
	}
	c.producer = producer
	c.Add(1)
	c.Add(2)
	c.Add(3)
	// Start consuming messages

	q := "q"
	consumer, err := common_rabbit.InitConsumer(cfg, q)
	if err != nil {
		c.Logger.Errorf("error", err)
	}
	c.Consumer = consumer

	err = consumer.StartConsuming(
		RabbitHAndler,
		"my_queue",
		[]string{"routing_key", "routing_key_2"},
		rabbitmq.WithConsumeOptionsConcurrency(10),
		rabbitmq.WithConsumeOptionsQueueDurable,
		rabbitmq.WithConsumeOptionsQuorum,
		rabbitmq.WithConsumeOptionsBindingExchangeName("events"),
		rabbitmq.WithConsumeOptionsBindingExchangeKind("topic"),
		rabbitmq.WithConsumeOptionsBindingExchangeDurable,
		rabbitmq.WithConsumeOptionsConsumerName("RabbitConsumer"),
	)

	return nil

}

func (c *RabbitBroker) Close() error {
	if err := c.channel.Close(); err != nil {
		return err
	}
	if err := c.conn.Close(); err != nil {
		return err
	}
	return nil
}

func (c *RabbitBroker) Add(key interface{}) error {
	msg := &models.Msg{
		Command: models.CommandAdd,
		Key:     key,
	}
	bytes, err := yaml.Marshal(msg)
	if err != nil {
		c.Logger.Errorf("PushMessageToQueue", "RabbitMQ", fmt.Sprintf("serialization error. %s", err.Error()), "widget-library-service")
		return nil
	}
	return c.producer.Publish(
		bytes,
		[]string{"routing_key"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsExchange("events"),
	)
}

func (c *RabbitBroker) Purge() error {
	msg := &models.Msg{
		Command: models.CommandPurge,
	}
	bytes, err := yaml.Marshal(msg)
	if err != nil {
		return err
	}
	return c.producer.Publish(
		bytes,
		[]string{"routing_key"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsExchange("events"),
	)
}
func RabbitHAndler(d rabbitmq.Delivery) rabbitmq.Action {
	messageBody := string(d.Body)
	log.Printf("Consumed message: %s", messageBody)
	switch models.RabbitCommand(messageBody) {
	case models.CommandPurge:
		// Логика для команды PURGE
		return rabbitmq.Ack
	case models.CommandAdd:
		// Логика для команды ADD
		return rabbitmq.Ack
	default:
		log.Printf("Unknown command: %s", messageBody)
		return rabbitmq.NackDiscard
	}
}
