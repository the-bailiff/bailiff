package config

import (
	"github.com/iamolegga/enviper"
	"github.com/spf13/viper"
)

func ProvideConfig() (c Config, err error) {
	v := viper.New()
	setup(v)
	err = enviper.New(v).Unmarshal(&c)
	return c, err
}
