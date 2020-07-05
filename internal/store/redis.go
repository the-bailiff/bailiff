package store

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/gorilla/sessions"
	"github.com/rbcervilla/redisstore"

	"github.com/the-bailiff/bailiff/internal/config"
)

func provideRedisStore(
	client *redis.Client,
	opts options,
	cfg config.Config,
) (sessions.Store, error) {
	store, err := redisstore.NewRedisStore(client)
	if err != nil {
		return nil, err
	}

	store.KeyPrefix(cfg.Store.Redis.KeyPrefix)
	store.Options(opts.Options)

	return store, nil
}

func provideRedisOptions(cfg config.Config) *redis.Options {
	v := &cfg.Store.Redis.Options
	fmt.Printf("%+v\n", v)
	return v
}
