package broker

import (
	"awesomeProject/internal/message_broker/broker_models"
	"context"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
)

type Broker struct {
	broker        []string
	id            string
	Chache_Broker broker_models.CacheBroker
	cache         *lru.TwoQueueCache
}

func NewBroker(brokers []string, cache *lru.TwoQueueCache, id string) broker_models.Broker {
	broker := Broker{broker: brokers, cache: cache, id: id}
	return &broker
}

func (b Broker) Connect(ctx context.Context) error {
	Brokers := []broker_models.SubBroker{b.Chache_Broker}

	fmt.Println("cache broker created")
	for i, Broker := range Brokers {
		fmt.Println("brokers ranging to connect")
		if err := Broker.Connect(ctx, b.broker[i]); err != nil {
			return err
		}
	}
	return nil
}

func (b Broker) Close() error {
	Brokers := []broker_models.SubBroker{b.Chache_Broker}
	for _, Broker := range Brokers {
		if err := Broker.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (b *Broker) Cache() broker_models.CacheBroker {
	if b.Chache_Broker == nil {
		b.Chache_Broker = NewCacheBroker(b.cache, b.id)
	}
	return b.Chache_Broker
}
