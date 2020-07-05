package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const appName = "bailiff"
const cfgFileName = "config"

func setup(v *viper.Viper) {
	setupFile(v)
	setupEnvVars(v)
	setupDefaults(v)
}

func setupFile(v *viper.Viper) {
	v.SetConfigName(cfgFileName)
	v.AddConfigPath(fmt.Sprintf("/etc/%s/", appName))
	v.AddConfigPath(fmt.Sprintf("$HOME/.%s", appName))
	v.AddConfigPath(".")
}

func setupEnvVars(v *viper.Viper) {
	v.SetEnvPrefix(appName)
}
