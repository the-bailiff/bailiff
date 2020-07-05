//+build wireinject

package store

import (
	"github.com/go-redis/redis"
	"github.com/google/wire"
	"github.com/gorilla/sessions"
	"github.com/rbcervilla/redisstore"

	"github.com/the-bailiff/bailiff/internal/config"
)

func redisStoreImpl() (sessions.Store, error) {
	wire.Build(
		config.ProvideConfig,
		provideOptions,
		provideRedisStore,
		provideRedisOptions,
		redis.NewClient,
	)
	return &redisstore.RedisStore{}, nil
}
