package http

import (
	"awesomeProject/internal/message_broker/broker_models"
	"awesomeProject/internal/store"

	lru "github.com/hashicorp/golang-lru"
)

type ServerOption func(srv *Server)

func WithAddress(addres string) ServerOption {
	return func(srv *Server) {
		srv.Addres = addres
	}
}

func WithStore(store store.Store) ServerOption {
	return func(srv *Server) {
		srv.store = store
	}

}

func WithCache(cache *lru.TwoQueueCache) ServerOption {
	return func(srv *Server) {
		srv.cache = cache
	}
}

func WithBroker(broker broker_models.Broker) ServerOption {
	return func(srv *Server) {

		srv.broker = broker

	}
}
