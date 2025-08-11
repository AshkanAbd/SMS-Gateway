package main

import (
	"github.com/AshkanAbd/arvancloud_sms_gateway/config"

	pkgCfg "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/config"
	pkgLog "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/logger"
)

var Config config.AppConfig

func init() {
	pkgLog.SetLogLevel("Trace")
	Config = *pkgCfg.Load[config.AppConfig]()
	pkgLog.SetLogLevel(Config.LogLevel)

	pkgLog.Info("Config is %#v", Config)
}

func main() {

}
