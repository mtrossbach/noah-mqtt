# noah-mqtt
![License](https://img.shields.io/github/license/mtrossbach/noah-mqtt) ![GitHub last commit](https://img.shields.io/github/last-commit/mtrossbach/noah-mqtt) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mtrossbach/noah-mqtt)

Polls and publishes data and metrics of the Growatt Noah battery for balcony power plants to MQTT for Home Assistant

UNDER CONSTRUCTION!

# Run with docker

To run the latest version make sure you have `Docker` installed. Execute the following command:

```
docker run --name noah-mqtt -e GROWATT_USERNAME=myusername -e GROWATT_PASSWORD=mypassword -e MQTT_HOST=localhost -e MQTT_PORT=1883 ghcr.io/mtrossbach/noah-mqtt:latest
```

Provide values for `GROWATT_USERNAME`, `GROWATT_PASSWORD`, `MQTT_HOST`, `MQTT_HOST`. The application will then connect to your MQTT broker and fetch all metrics and data for all your Noah devices in your account.

# Build and Run

To build the application, make sure to have a current version of `Go` installed on your machine.

```
go build -o noah-mqtt cmd/noah-mqtt/main.go
```

Don't forget to specify your configuration using the environment variables when you run the application.

```
GROWATT_USERNAME=username GROWATT_PASSWORD=mypassword MQTT_HOST=localhost MQTT_PORT=1883 ./noah-mqtt
```

Provide values for `GROWATT_USERNAME`, `GROWATT_PASSWORD`, `MQTT_HOST`, `MQTT_HOST`. The application will then connect to your MQTT broker and fetch all metrics and data for all your Noah devices in your account.


# Configuration

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


