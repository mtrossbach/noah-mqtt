package config

import (
	"os"
	"strconv"
	"sync"
	"time"
)

type Config struct {
	Version         string
	PollingInterval time.Duration
	Growatt         Growatt
	Mqtt            Mqtt
}

type Growatt struct {
	Username     string
	Password     string
	SerialNumber string
	PlantId      string
}

type Mqtt struct {
	Host     string
	Port     int
	ClientId string
	Username string
	Password string
	Topic    string
}

var _config Config
var _once sync.Once

func Get() Config {
	_once.Do(func() {
		_config = Config{
			Version:         getEnv("VERSION", "local"),
			PollingInterval: time.Duration(s2i(getEnv("POLLING_INTERVAL", "10"))) * time.Second,
			Growatt: Growatt{
				Username:     getEnv("GROWATT_USERNAME", ""),
				Password:     getEnv("GROWATT_PASSWORD", ""),
				SerialNumber: getEnv("GROWATT_SERIAL_NUMBER", ""),
				PlantId:      getEnv("GROWATT_PLANT_ID", ""),
			},
			Mqtt: Mqtt{
				Host:     getEnv("MQTT_HOST", "localhost"),
				Port:     s2i(getEnv("MQTT_PORT", "1883")),
				ClientId: getEnv("MQTT_CLIENT_ID", "noah-mqtt"),
				Username: getEnv("MQTT_USERNAME", ""),
				Password: getEnv("MQTT_PASSWORD", ""),
				Topic:    getEnv("MQTT_TOPIC", "home/sensor/noah"),
			},
		}
	})
	return _config
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func s2i(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}