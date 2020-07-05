package config

import (
	"reflect"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

// Config

type Config struct {
	Port    string
	Rolling bool
	Proxy   string
	Cookie  Cookie
	Headers Headers
	Store   Store
}

func setupDefaults(v *viper.Viper) {
	v.SetDefault("Port", "80")
	v.SetDefault("Proxy", "8000")
	setCookieDefaults(v)
	setupHeadersDefaults(v)
	setupStoreDefaults(v)
}

// Cookie

type Cookie struct {
	Name     string
	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite string
}

func setCookieDefaults(v *viper.Viper) {
	v.SetDefault("Cookie.Name", "session")
}

// Headers

type Headers struct {
	ValuePrefix string
	SetPrefix   string
	DelPrefix   string
}

func setupHeadersDefaults(v *viper.Viper) {
	v.SetDefault("Headers.ValuePrefix", "X-Session-")
	v.SetDefault("Headers.SetPrefix", "X-Session-Set-")
	v.SetDefault("Headers.DelPrefix", "X-Session-Del-")
}

// Store

type Store struct {
	Redis StoreRedis
}

type StoreRedis struct {
	KeyPrefix     string
	redis.Options `mapstructure:",squash"`
}

func (s *StoreRedis) Empty() bool {
	return s.KeyPrefix == "" && redisOptionsEmpty(s.Options)
}

func redisOptionsEmpty(o redis.Options) bool {
	return reflect.DeepEqual(&o, &redis.Options{})
}

func setupStoreDefaults(v *viper.Viper) {
	v.SetDefault("Store.Redis.Addr", "localhost:6379")
}
