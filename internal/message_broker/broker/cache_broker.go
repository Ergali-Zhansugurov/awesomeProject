/*
package broker

import (

	"context"
	"encoding/json"
	"log"

	"manga/internal/message_broker/broker_models"
	"manga/internal/models"

	"github.com/IBM/sarama"
	lru "github.com/hashicorp/golang-lru"

)
const cacheTopic = "cache"
type (

	CacheBroker struct {
		syncProducer  sarama.SyncProducer
		consumerGroup sarama.ConsumerGroup

		consumeHandler *cacheConsumeHandler
		clientID       string
	}
	cacheConsumeHandler struct {
		cache *lru.TwoQueueCache
		ready chan bool
	}

)

	func NewCacheBroker(cache *lru.TwoQueueCache, clientID string) broker_models.CacheBroker {
		return &CacheBroker{
			clientID: clientID,
			consumeHandler: &cacheConsumeHandler{cache: cache, ready: make(chan bool)},
		}
	}

	func (c CacheBroker) Connect(ctx context.Context, brokers []string) error {
		producerConfig := sarama.NewConfig()
		producerConfig.Producer.RequiredAcks = sarama.WaitForAll
		producerConfig.Producer.Retry.Max = 10
		producerConfig.Producer.Return.Successes = true

		syncProducer, err := sarama.NewSyncProducer(brokers, producerConfig)
		if err != nil {
			return err
		}
		c.syncProducer = syncProducer

		consumerConfig := sarama.NewConfig()
		consumerConfig.Consumer.Return.Errors = true
		consumerGroup, err := sarama.NewConsumerGroup(brokers, c.clientID, consumerConfig)
		if err != nil {
			return err
		}
		c.consumerGroup = consumerGroup

		go func() {
			for {
				if err := c.consumerGroup.Consume(ctx, []string{"cacheTopic"}, c.consumeHandler); err != nil {
					log.Panicf("Error from consumer: %v", err)
				}
				if ctx.Err() != nil {
					return
				}
				c.consumeHandler.ready = make(chan bool)
			}
		}()
		<-c.consumeHandler.ready

		return nil
	}

	func (c *CacheBroker) Close() error {
		if err := c.syncProducer.Close(); err != nil {
			return err
		}

		if err := c.consumerGroup.Close(); err != nil {
			return err
		}

		return nil
	}

	func (c *CacheBroker) Remove(key interface{}) error {
		msg := &models.CacheMsg{
			Command: models.CacheCommandRemove,
			Key:     key,
		}

		msgRaw, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		_, _, err = c.syncProducer.SendMessage(&sarama.ProducerMessage{
			Topic: cacheTopic,
			Value: sarama.StringEncoder(msgRaw),
		})
		if err != nil {
			return err
		}

		return nil
	}

	func (c *CacheBroker)Add(key interface{})error{
		msg := &models.CacheMsg{
			Command: models.CacheCommandAdd,
			Key: key,
		}

		msgRaw, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		_, _, err = c.syncProducer.SendMessage(&sarama.ProducerMessage{
			Topic: cacheTopic,
			Value: sarama.StringEncoder(msgRaw),
		})
		if err != nil {
			return err
		}

		return nil
	}

	func (c *CacheBroker) Purge() error {
		msg := &models.CacheMsg{
			Command: models.CacheCommandPurge,
		}

		msgRaw, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		_, _, err = c.syncProducer.SendMessage(&sarama.ProducerMessage{
			Topic: cacheTopic,
			Value: sarama.StringEncoder(msgRaw),
		})
		if err != nil {
			return err
		}

		return nil
	}

	func (c *cacheConsumeHandler) Setup(session sarama.ConsumerGroupSession) error {
		close(c.ready)
		return nil
	}

	func (c *cacheConsumeHandler) Cleanup(session sarama.ConsumerGroupSession) error {
		return nil
	}

	func (c *cacheConsumeHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
		for msg := range claim.Messages() {
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic)

			cacheMsg := new(models.CacheMsg)
			if err := json.Unmarshal(msg.Value, cacheMsg); err != nil {
				return err
			}

			switch cacheMsg.Command {
			case models.CacheCommandRemove:
				c.cache.Remove(cacheMsg.Key)
			case models.CacheCommandPurge:
				c.cache.Purge()
			}

			session.MarkMessage(msg, "")
		}

		return nil
	}
*/
package broker

import (
	"context"
	"encoding/json"
	"log"

	"awesomeProject/internal/message_broker/broker_models"
	"awesomeProject/internal/models"

	lru "github.com/hashicorp/golang-lru"
	"github.com/rabbitmq/amqp091-go"
)

const cacheTopic = "cache"

type (
	CacheBroker struct {
		conn          *amqp091.Connection
		channel       *amqp091.Channel
		cache         *lru.TwoQueueCache
		clientID      string
		consumeHandle *cacheConsumeHandler
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

func (c *CacheBroker) Connect(ctx context.Context, amqpURI string) error {

	conn, err := amqp091.Dial(amqpURI)
	if err != nil {
		return err
	}
	c.conn = conn

	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	c.channel = channel

	// Declare the exchange and queue
	err = c.channel.ExchangeDeclare(
		cacheTopic, // name
		"fanout",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return err
	}

	_, err = c.channel.QueueDeclare(
		cacheTopic, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return err
	}

	err = c.channel.QueueBind(
		cacheTopic, // queue name
		cacheTopic, // routing key
		cacheTopic, // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Start consuming messages
	msgs, err := c.channel.Consume(
		cacheTopic, // queue
		c.clientID, // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			log.Printf("Message received: value = %s", string(msg.Body))
			cacheMsg := new(models.CacheMsg)
			if err := json.Unmarshal(msg.Body, cacheMsg); err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				continue
			}

			switch cacheMsg.Command {
			case models.CacheCommandRemove:
				c.cache.Remove(cacheMsg.Key)
			case models.CacheCommandPurge:
				c.cache.Purge()
			}
		}
	}()

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

func (c *CacheBroker) Remove(key interface{}) error {
	msg := &models.CacheMsg{
		Command: models.CacheCommandRemove,
		Key:     key,
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

func (c *CacheBroker) Add(key interface{}) error {
	msg := &models.CacheMsg{
		Command: models.CacheCommandAdd,
		Key:     key,
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
