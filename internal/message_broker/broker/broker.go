package broker

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/message_broker/broker_models"
	"context"
	lru "github.com/hashicorp/golang-lru"
)

type Broker struct {
	configs       config.MainConfig
	id            string
	Chache_Broker broker_models.CacheBroker
	cache         *lru.TwoQueueCache
}

func NewBroker(cfg config.MainConfig, cache *lru.TwoQueueCache, id string) broker_models.Broker {
	broker := Broker{configs: cfg, cache: cache, id: id, Chache_Broker: NewCacheBroker(cache, "10")}
	return &broker
}

func (b Broker) Connect(ctx context.Context) error {
	Brokers := []broker_models.SubBroker{b.Chache_Broker}
	for _, Broker := range Brokers {
		if err := Broker.Connect(ctx, b.configs); err != nil {
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
