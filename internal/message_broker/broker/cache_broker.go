package broker

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/logger"
	"awesomeProject/internal/rabbit/common_rabbit"
	"context"
	"encoding/json"
	"fmt"
	"github.com/wagslane/go-rabbitmq"
	"gopkg.in/yaml.v3"
	"log"

	"awesomeProject/internal/message_broker/broker_models"
	"awesomeProject/internal/models"

	lru "github.com/hashicorp/golang-lru"
	"github.com/rabbitmq/amqp091-go"
)

const cacheTopic = "cache"

type (
	CacheBroker struct {
		Logger        *logger.Logger
		conn          *rabbitmq.Conn
		channel       *amqp091.Channel
		cache         *lru.TwoQueueCache
		clientID      string
		Consumer      *rabbitmq.Consumer
		consumeHandle *cacheConsumeHandler
		producer      *rabbitmq.Publisher
		queue         *amqp091.Queue
	}

	cacheConsumeHandler struct {
		cache *lru.TwoQueueCache
	}
)

func NewCacheBroker(cache *lru.TwoQueueCache, clientID string) broker_models.CacheBroker {
	return &CacheBroker{
		cache:         cache,
		clientID:      clientID,
		consumeHandle: &cacheConsumeHandler{cache: cache},
	}
}

func (c *CacheBroker) Connect(ctx context.Context, cfg config.MainConfig) error {
	common_rabbit.CreateQueue(cfg, "q")

	producer, err := common_rabbit.InitProducer(cfg)
	if err != nil {
		c.Logger.Errorf("error", err)
	}
	c.producer = producer
	// Start consuming messages
	q := "q"
	consumer, err := common_rabbit.InitConsumer(cfg, q)
	if err != nil {
		c.Logger.Errorf("error", err)
	}
	c.Consumer = consumer
	err = c.Consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
		log.Printf("consumed: %v", string(d.Body))
		// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
		return rabbitmq.Ack
	})
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (c *CacheBroker) Close() error {
	if err := c.channel.Close(); err != nil {
		return err
	}
	if err := c.conn.Close(); err != nil {
		return err
	}
	return nil
}

func (c *CacheBroker) Add(key interface{}) error {
	msg := &models.CacheMsg{
		Command: models.CacheCommandAdd,
		Key:     key,
	}
	headers := make(map[string]interface{})
	bytes, err := yaml.Marshal(msg)
	if err != nil {
		c.Logger.Errorf("PushMessageToQueue", "RabbitMQ", fmt.Sprintf("serialization error. %s", err.Error()), "widget-library-service")
		return nil
	}
	return c.producer.Publish(bytes, []string{c.queue.Name}, func(options *rabbitmq.PublishOptions) {
		options.Headers = headers
	})
}

func (c *CacheBroker) Purge() error {
	msg := &models.CacheMsg{
		Command: models.CacheCommandPurge,
	}
	msgRaw, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = c.channel.Publish(
		cacheTopic, // exchange
		cacheTopic, // routing key
		false,      // mandatory
		false,      // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        msgRaw,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
