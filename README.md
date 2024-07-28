# noah-mqtt
![License](https://img.shields.io/github/license/mtrossbach/noah-mqtt) ![GitHub last commit](https://img.shields.io/github/last-commit/mtrossbach/noah-mqtt) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mtrossbach/noah-mqtt)

Polls and publishes data and metrics of the Growatt Noah battery for balcony power plants to MQTT for Home Assistant

UNDER CONSTRUCTION!

## Configuration

Use the following environment variables to configure noah-mqtt:

| Environment Variable | Description                                     | Default       |
|:---------------------|:------------------------------------------------|:--------------|
| POLLING_INTERVAL     | Interval between new data is fetched in seconds | 10            |
| GROWATT_USERNAME     | Username for your Growatt account (required)    | -             |
| GROWATT_PASSWORD     | Password for your Growatt account (required)    | -             |
| MQTT_HOST            | MQTT broker host (required)                     | -             |
| MQTT_PORT            | MQTT broker port                                | 1883          |
| MQTT_CLIENT_ID       | MQTT client id                                  | noah-mqtt     |
| MQTT_USERNAME        | MQTT username                                   | -             |
| MQTT_PASSWORD        | MQTT password                                   | -             |
| MQTT_TOPIC_PREFIX    | MQTT base topic                                 | noah2mqtt     |
| HOMEASSISTANT_TOPIC_PREFIX    | HomeAssistant base topic                        | homeassistant |


