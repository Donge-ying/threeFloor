package config

import (
	"github.com/FishGoddess/logit"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	Log *logit.Logger
)

func init() {
	Log = logit.NewLogger()
	log.SetOutput(os.Stdout)
	viper.SetConfigFile("app.toml")
	err := viper.ReadInConfig()
	if err != nil {
		msg := "Fatal error: Failed to read configuration file: " + err.Error()
		Log.Error(msg)
		log.Fatal(msg)
	}
}
