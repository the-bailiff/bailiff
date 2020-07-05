package store

import (
	"errors"
	"strings"

	"github.com/gorilla/sessions"

	"github.com/the-bailiff/bailiff/internal/config"
)

const providerRedis = "REDIS"

func newProviderTypeError() error {
	allowedTypes := []string{providerRedis}
	t := "invalid Store Provider type, "
	t += "should be one of: "
	t += strings.Join(allowedTypes, ", ")
	return errors.New(t)
}

func ProvideBaseStore(cfg config.Config) (sessions.Store, error) {
	switch {
	case !cfg.Store.Redis.Empty():
		return redisStoreImpl()
	default:
		return nil, newProviderTypeError()
	}
}
