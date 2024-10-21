package http

import (
	"awesomeProject/internal/message_broker/broker_models"
	"awesomeProject/internal/store"
	lru "github.com/hashicorp/golang-lru"
)

type BaseResource struct {
	store  store.Store
	broker broker_models.Broker
	cache  *lru.TwoQueueCache
}

func NewBaseResource(store store.Store, cache *lru.TwoQueueCache, broker broker_models.Broker) *BaseResource {
	return &BaseResource{store: store, cache: cache, broker: broker}
}
