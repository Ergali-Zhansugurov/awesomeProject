package broker_models

import "context"

type Broker interface {
	Connect(ctx context.Context) error
	Close() error
	Cache() CacheBroker
}
type SubBroker interface {
	Connect(ctx context.Context, brokers string) error
	Close() error
}
type CacheBroker interface {
	SubBroker
	Add(key interface{}) error
	Remove(key interface{}) error
	Purge() error
}
