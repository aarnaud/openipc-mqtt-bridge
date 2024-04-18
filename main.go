package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	mqtt_client "openipc-mqtt-bridge/mqtt-client"
	"openipc-mqtt-bridge/utils"
	"os"
	"path"
	"regexp"
	"time"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	config := utils.GetConfig()
	hostRegex := regexp.MustCompile(fmt.Sprintf(`^%s/(.+?)/`, config.MQTT.BaseTopic))
	cli := mqtt_client.NewMQTT(config)

	go func() {
		for {
			cli.PublishAvailability()
			time.Sleep(time.Second * 30)
		}
	}()
	cli.WatchTopicPlayOnSpeaker(func(client mqtt.Client, message mqtt.Message) {
		matches := hostRegex.FindStringSubmatch(message.Topic())
		if matches == nil && len(matches) < 2 {
			log.Info().Msgf("failed to extract camera host from topic: %s", message.Topic())
			return
		}
		cameraHost := matches[1]
		cameraSpeakerUrl := fmt.Sprintf("http://%s/play_audio", cameraHost)
		file, err := os.Open(path.Join(config.AudioFilesPath, string(message.Payload())))
		if err != nil {
			log.Error().Msgf("failed to open %s", message.Payload())
			return
		}
		defer file.Close()
		httpCli := &http.Client{
			Timeout: time.Second * 30,
		}
		req, err := http.NewRequest(http.MethodPost, cameraSpeakerUrl, file)
		req.Header.Add("Content-Type", "binary/octet-stream")
		req.SetBasicAuth(config.CameraUser, config.CameraPassword)
		response, err := httpCli.Do(req)
		if err != nil {
			log.Error().Msgf("failed to send %s to %s with error %s", message.Payload(), cameraHost, err)
			return
		}
		if response.StatusCode != http.StatusOK {
			log.Error().Msgf("failed to send %s to %s with status %s",
				message.Payload(), cameraHost, response.StatusCode)
			return
		}
		log.Info().Msgf("play %s on %s", message.Payload(), cameraHost)
	})
}
