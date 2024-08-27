package broker_models

import (
	"awesomeProject/internal/config"
	"context"
)

type Broker interface {
	Connect(ctx context.Context) error
	Close() error
	Cache() CacheBroker
}
type SubBroker interface {
	Connect(ctx context.Context, cfg config.MainConfig) error
	Close() error
}
type CacheBroker interface {
	SubBroker
	Add(key interface{}) error
	Purge() error
}
