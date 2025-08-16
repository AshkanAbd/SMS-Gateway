package config

import (
	"strings"

	"github.com/spf13/viper"

	pkgLog "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/logger"
)

func Load[T any]() *T {
	viper.AddConfigPath(".")
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	t := new(T)

	if err := viper.ReadInConfig(); err != nil {
		pkgLog.Fatal(err, "failed to read config")
	} else {
		pkgLog.Debug("config file loaded successfully")
	}

	if err := viper.Unmarshal(&t); err != nil {
		pkgLog.Fatal(err, "failed to unmarshal config")
	} else {
		pkgLog.Debug("config file unmarshalled successfully")
	}

	return t
}
