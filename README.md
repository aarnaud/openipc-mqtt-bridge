# OpenIPC MQTT Bridge

## Features:

- [X] MQTT integration for Home Assistant
- [X] Play on speaker
- [ ] Adjust Majestic config 


## Config:

Exemple `openipc-mqtt-bridge.yaml`

```yaml
AUDIO_FILES_PATH: "/PATH_TO_YOUR_AUDIO_FILES"
MQTT_BROKER_HOST: "192.168.1.245" # require if MQTT_ENABLED
MQTT_BROKER_PORT: 1883
MQTT_CLIENT_ID: "openipc"
MQTT_BASE_TOPIC: "openipc"
MQTT_USERNAME: "openipc"
MQTT_PASSWORD: "CHANGEME"
```


## Example with mosquitto 

```bash
mosquitto_pub -h 192.168.1.123 -u openipc -P CHANGEME -t "openipc/192.168.1.234/playonspeaker"  -m "siren.wav"
```

## Example with Home Assistant

```yaml
        - service: mqtt.publish
          data:
            topic: openipc/192.168.1.234/playonspeaker
            payload: "siren.wav"
```