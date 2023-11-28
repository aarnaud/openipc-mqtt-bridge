package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	AudioFilesPath string
	MQTT           *ConfigMQTT
}

type ConfigMQTT struct {
	BrokerHost string
	BrokerPort int
	ClientID   string
	BaseTopic  string
	Username   string
	Password   string
}

func GetConfig() *Config {
	// the env registry will look for env variables that start with "INTERCOM_".
	viper.SetEnvPrefix("OIPCSB")
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()                       // To get the value from the config file using key// viper package read .env
	viper.SetConfigName("openipc-mqtt-bridge") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.ReadInConfig()
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	viper.SetDefault("AUDIO_FILES_PATH", pwd)
	viper.SetDefault("MQTT_BROKER_PORT", 1883)
	viper.SetDefault("MQTT_CLIENT_ID", "openipc")
	viper.SetDefault("MQTT_BASE_TOPIC", "openipc")

	config := Config{
		AudioFilesPath: viper.GetString("AUDIO_FILES_PATH"),
		MQTT: &ConfigMQTT{
			BrokerHost: viper.GetString("MQTT_BROKER_HOST"),
			BrokerPort: viper.GetInt("MQTT_BROKER_PORT"),
			ClientID:   viper.GetString("MQTT_CLIENT_ID"),
			BaseTopic:  viper.GetString("MQTT_BASE_TOPIC"),
			Username:   viper.GetString("MQTT_USERNAME"),
			Password:   viper.GetString("MQTT_PASSWORD"),
		},
	}
	return &config
}
