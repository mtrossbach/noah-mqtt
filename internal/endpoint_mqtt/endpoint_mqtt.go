package endpoint_mqtt

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"nexa-mqtt/internal/endpoint"
	"nexa-mqtt/internal/homeassistant"
	"nexa-mqtt/pkg/models"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Options struct {
	MqttClient  mqtt.Client
	TopicPrefix string
	HaClient    homeassistant.HaClient
}

type Endpoint struct {
	opts          Options
	devs          []models.NoahDevicePayload
	param_applier endpoint.ParameterApplier
	stateLock     sync.Mutex
	lastParameter models.ParameterPayload
	newParameter  models.ParameterPayload
	publishTimer  *time.Timer
}

func NewEndpoint(options Options) *Endpoint {
	return &Endpoint{
		opts:          options,
		lastParameter: models.EmptyParameterPayload(),
	}
}

func (e *Endpoint) SetParameterApplier(applier endpoint.ParameterApplier) {
	e.param_applier = applier
}

func (e *Endpoint) SetDevices(devices []models.NoahDevicePayload) {
	for _, dev := range e.devs {
		e.opts.MqttClient.Unsubscribe(parameterCommandTopic(e.opts.TopicPrefix, dev.Serial))
	}

	e.devs = devices

	for _, dev := range devices {
		e.opts.MqttClient.Subscribe(parameterCommandTopic(e.opts.TopicPrefix, dev.Serial), 0, e.parametersSubscription(dev))
	}

	var haDevices []homeassistant.DeviceInfo
	for _, dev := range devices {
		var bats []homeassistant.BatteryInfo
		for i, bat := range dev.Batteries {
			bats = append(bats, homeassistant.BatteryInfo{
				Alias:      bat.Alias,
				StateTopic: stateTopicBattery(e.opts.TopicPrefix, dev.Serial, i),
			})
		}
		haDevices = append(haDevices, homeassistant.DeviceInfo{
			SerialNumber:          dev.Serial,
			Model:                 dev.Model,
			Version:               dev.Version,
			Alias:                 dev.Alias,
			StateTopic:            deviceStateTopic(e.opts.TopicPrefix, dev.Serial),
			ParameterStateTopic:   parameterStateTopic(e.opts.TopicPrefix, dev.Serial),
			ParameterCommandTopic: parameterCommandTopic(e.opts.TopicPrefix, dev.Serial),
			Batteries:             bats,
		})
	}

	e.opts.HaClient.SetDevices(haDevices)
}

func (e *Endpoint) PublishDeviceStatus(device models.NoahDevicePayload, status models.DevicePayload) {
	if b, err := json.Marshal(status); err != nil {
		slog.Error("could not marshal device status data", slog.String("error", err.Error()))
	} else {
		e.opts.MqttClient.Publish(deviceStateTopic(e.opts.TopicPrefix, device.Serial), 0, false, string(b))
		slog.Debug("device status sent to mqtt", slog.String("data", string(b)), slog.String("device", device.Serial))
	}
}

func (e *Endpoint) PublishBatteryDetails(device models.NoahDevicePayload, details []models.BatteryPayload) {
	var logData []any
	for i, bat := range details {
		if b, err := json.Marshal(bat); err != nil {
			slog.Error("could not marshal battery data", slog.String("error", err.Error()))
		} else {
			e.opts.MqttClient.Publish(stateTopicBattery(e.opts.TopicPrefix, device.Serial, i), 0, false, string(b))
			logData = append(logData, slog.String(fmt.Sprintf("BAT%d", i), string(b)))
		}
	}
	logData = append(logData, slog.String("device", device.Serial))
	slog.Debug("battery data sent to mqtt", logData...)
}

func (e *Endpoint) PublishParameterData(device models.NoahDevicePayload, param models.ParameterPayload) {
	if b, err := json.Marshal(param); err != nil {
		slog.Error("could not marshal parameter data", slog.String("error", err.Error()), slog.String("device", device.Serial))
	} else {
		e.opts.MqttClient.Publish(parameterStateTopic(e.opts.TopicPrefix, device.Serial), 0, false, string(b))
		slog.Debug("parameter data sent to mqtt", slog.String("data", string(b)), slog.String("device", device.Serial))

		e.stateLock.Lock()
		defer e.stateLock.Unlock()

		e.lastParameter = param
	}
}

const debounceDelay = 500 * time.Millisecond

func (e *Endpoint) parametersSubscription(dev models.NoahDevicePayload) func(client mqtt.Client, message mqtt.Message) {
	return func(client mqtt.Client, message mqtt.Message) {
		if e.param_applier == nil {
			slog.Error("no parameter applier is set or support. parameter changes are not applied!")
			return
		}

		var payload models.ParameterPayload
		if err := json.Unmarshal(message.Payload(), &payload); err != nil {
			slog.Error("unable to unmarshal parameter command payload", slog.String("error", err.Error()))
			return
		}

		e.stateLock.Lock()
		defer e.stateLock.Unlock()

		e.newParameter.UpdateFrom(payload)

		if e.publishTimer != nil {
			e.publishTimer.Stop()
		}

		e.publishTimer = time.AfterFunc(debounceDelay, func() {
			e.debouncedParametersSubscription(dev)
		})
	}
}

func (e *Endpoint) debouncedParametersSubscription(dev models.NoahDevicePayload) {
	e.stateLock.Lock()
	defer e.stateLock.Unlock()

	e.lastParameter.UpdateFrom(e.newParameter)

	if e.newParameter.DefaultACCouplePower != nil || e.newParameter.DefaultMode != nil {
		e.param_applier.SetOutputPowerW(dev, *e.lastParameter.DefaultMode, *e.lastParameter.DefaultACCouplePower)
	}

	if e.newParameter.ChargingLimit != nil || e.newParameter.DischargeLimit != nil {
		e.param_applier.SetChargingLimits(dev, *e.lastParameter.ChargingLimit, *e.lastParameter.DischargeLimit)
	}

	e.newParameter = models.ParameterPayload{}
	e.publishTimer = nil

	go e.PublishParameterData(dev, e.lastParameter)
}
