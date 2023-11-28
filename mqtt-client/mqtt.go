package mqtt_client

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"openipc-mqtt-bridge/utils"
	"path"
	"time"
)

type Client struct {
	config                           *utils.ConfigMQTT
	instance                         mqtt.Client
	topicAvailability                string
	topicSetVol                      string
	topicSetGain                     string
	topicPlayOnSpeaker               string
	onConnectWatchTopicPlayOnSpeaker chan bool
}

func NewMQTT(config *utils.Config) *Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.MQTT.BrokerHost, config.MQTT.BrokerPort))
	opts.SetClientID(config.MQTT.ClientID)
	opts.SetUsername(config.MQTT.Username)
	opts.SetPassword(config.MQTT.Password)

	onConnectWatchTopicPlayOnSpeaker := make(chan bool, 1)
	opts.OnConnect = func(client mqtt.Client) {
		log.Info().Msg("MQTT Connected")
		onConnectWatchTopicPlayOnSpeaker <- true
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Err(err).Msgf("MQTT broker connection lost")
	}

	opts.ConnectRetryInterval = 5 * time.Second

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return &Client{
		config:                           config.MQTT,
		instance:                         client,
		topicAvailability:                path.Join(config.MQTT.BaseTopic, "available"),
		topicPlayOnSpeaker:               path.Join(config.MQTT.BaseTopic, "+", "playonspeaker"),
		onConnectWatchTopicPlayOnSpeaker: onConnectWatchTopicPlayOnSpeaker,
	}
}

func (c *Client) WatchTopicPlayOnSpeaker(callback mqtt.MessageHandler) {
	for {
		// wait for connection
		<-c.onConnectWatchTopicPlayOnSpeaker
		// https://www.home-assistant.io/integrations/button.mqtt/
		token := c.instance.Subscribe(c.topicPlayOnSpeaker, 1, callback)
		token.WaitTimeout(5 * time.Second)
		if !token.WaitTimeout(2 * time.Second) {
			log.Warn().Msgf("timeout to subscribe to topic %s", c.topicPlayOnSpeaker)
		}
		if token.Error() != nil {
			log.Error().Err(token.Error()).Msgf("failed to subscribe to topic %s", c.topicPlayOnSpeaker)
		}
		log.Info().Msgf("Subscribed to topic: %s", c.topicPlayOnSpeaker)
	}
}

func (c *Client) PublishAvailability() {
	// https://www.home-assistant.io/integrations/button.mqtt/
	log.Debug().Msgf("PublishAvailability to topic: %s", c.topicAvailability)
	token := c.instance.Publish(c.topicAvailability, 0, false, "online")
	if !token.WaitTimeout(2 * time.Second) {
		log.Warn().Msgf("timeout to publish availability to topic %s", c.topicAvailability)
	}
	if token.Error() != nil {
		log.Error().Err(token.Error()).Msgf("failed to publish availability to topic %s", c.topicAvailability)
	}
}
