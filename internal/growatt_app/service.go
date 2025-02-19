package growatt_app

import (
	"fmt"
	"log/slog"
	"noah-mqtt/internal/endpoint"
	"noah-mqtt/internal/misc"
	"noah-mqtt/pkg/models"
	"os"
	"time"
)

type Options struct {
	ServerUrl                     string
	Username                      string
	Password                      string
	PollingInterval               time.Duration
	BatteryDetailsPollingInterval time.Duration
	ParameterPollingInterval      time.Duration
}
type GrowattAppService struct {
	opts      Options
	client    *Client
	devices   []models.NoahDevicePayload
	endpoints []endpoint.Endpoint
	loggedIn  bool
}

func NewGrowattAppService(options Options) *GrowattAppService {
	return &GrowattAppService{
		opts:     options,
		client:   newClient(options.ServerUrl, options.Username, options.Password),
		loggedIn: false,
	}
}

func (g *GrowattAppService) Login() error {
	slog.Info("logging in to growatt (app)...")

	if err := g.client.Login(); err != nil {
		return err
	}
	g.loggedIn = true
	return nil
}

func (g *GrowattAppService) StartPolling() {
	g.enumerateDevices()
	go g.poll()
}

func (g *GrowattAppService) fetchDevices() []models.NoahDevicePayload {
	slog.Info("fetching plant list")
	list, err := g.client.GetPlantList()
	if err != nil {
		slog.Error("could not get plant list", slog.String("error", err.Error()))
		misc.Panic(err)
	}

	var devices []models.NoahDevicePayload

	for _, plant := range list.PlantList {
		slog.Info("fetch plant details", slog.Int("plantId", plant.ID))
		if info, err := g.client.GetNoahPlantInfo(fmt.Sprintf("%d", plant.ID)); err != nil {
			slog.Error("could not get plant info", slog.Int("plantId", plant.ID), slog.String("error", err.Error()))
		} else {
			if len(info.Obj.DeviceSn) > 0 {
				devices = append(devices, models.NoahDevicePayload{
					PlantId:   plant.ID,
					Serial:    info.Obj.DeviceSn,
					Batteries: nil,
				})
				slog.Info("found device sn", slog.String("deviceSn", info.Obj.DeviceSn), slog.Int("plantId", plant.ID))
			}
		}
	}

	if len(devices) == 0 {
		slog.Info("no noah devices found")
		<-time.After(60 * time.Second)
		os.Exit(0)
	}

	return devices
}

func (g *GrowattAppService) enumerateDevices() {
	devices := g.fetchDevices()

	for i, device := range devices {
		if data, err := g.client.GetNoahInfo(device.Serial); err != nil {
			slog.Error("could not noah status", slog.String("error", err.Error()), slog.String("serialNumber", device.Serial))
		} else {
			batCount := len(data.Obj.Noah.BatSns)
			var batteries []models.NoahDeviceBatteryPayload
			for i := 0; i < batCount; i++ {
				batteries = append(batteries, models.NoahDeviceBatteryPayload{
					Alias: fmt.Sprintf("BAT%d", i),
				})
			}

			devices[i].Model = data.Obj.Noah.Model
			devices[i].Version = data.Obj.Noah.Version
			devices[i].Alias = data.Obj.Noah.Alias
			devices[i].Batteries = batteries
		}
	}

	g.devices = devices

	for _, e := range g.endpoints {
		e.SetDevices(devices)
	}
}

func (g *GrowattAppService) AddEndpoint(e endpoint.Endpoint) {
	g.endpoints = append(g.endpoints, e)
	e.SetParameterApplier(g)
}

func (g *GrowattAppService) ensureParameterLogin() bool {
	if !g.loggedIn {
		if err := g.Login(); err != nil {
			slog.Error("could not login to growatt account (app)", slog.String("error", err.Error()))
			return false
		}
	}
	return true
}

func (g *GrowattAppService) SetOutputPowerW(device models.NoahDevicePayload, power float64) bool {
	slog.Info("trying to set default power (app)", slog.String("device", device.Serial), slog.Int("power", int(power)))
	if !g.ensureParameterLogin() {
		slog.Error("unable to set default power (app)", slog.String("device", device.Serial))
		return false
	}
	if err := g.client.SetDefaultPower(device.Serial, power); err != nil {
		slog.Error("unable to set default power (app)", slog.String("error", err.Error()), slog.String("device", device.Serial))
		return false
	} else {
		go g.pollParameterData(device)
		slog.Info("set default power (app)", slog.String("device", device.Serial), slog.Int("power", int(power)))
		return true
	}
}
func (g *GrowattAppService) SetChargingLimit(device models.NoahDevicePayload, limit float64) bool {
	slog.Info("trying to set charging limit (app)", slog.String("device", device.Serial), slog.Float64("limit", limit))
	if !g.ensureParameterLogin() {
		slog.Error("unable to set charging limit (app)", slog.String("device", device.Serial))
		return false
	}
	if data, err := g.client.GetNoahInfo(device.Serial); err != nil {
		slog.Error("unable to get parameter status (app)", slog.String("error", err.Error()))
		return false
	} else {
		dl := misc.ParseFloat(data.Obj.Noah.ChargingSocLowLimit)

		slog.Info("trying to set charging limit (app)", slog.String("device", device.Serial), slog.Float64("chargingLimit", limit), slog.Float64("dischargeLimit", dl))
		if err := g.client.SetSocLimit(device.Serial, limit, dl); err != nil {
			slog.Error("unable to set charging limit (app)", slog.String("error", err.Error()))
			return false
		} else {
			go g.pollParameterData(device)
			slog.Info("set charging limit (app)", slog.String("device", device.Serial), slog.Float64("chargingLimit", limit), slog.Float64("dischargeLimit", dl))
			return true
		}
	}
}

func (g *GrowattAppService) SetDischargeLimit(device models.NoahDevicePayload, limit float64) bool {
	slog.Info("trying to set discharge limit (app)", slog.String("device", device.Serial), slog.Float64("limit", limit))
	if !g.ensureParameterLogin() {
		slog.Error("unable to set discharge limit (app)", slog.String("device", device.Serial))
		return false
	}
	if data, err := g.client.GetNoahInfo(device.Serial); err != nil {
		slog.Error("unable to get parameter status (app)", slog.String("error", err.Error()))
		return false
	} else {
		cl := misc.ParseFloat(data.Obj.Noah.ChargingSocHighLimit)

		slog.Info("trying to set discharge limit (app)", slog.String("device", device.Serial), slog.Float64("chargingLimit", cl), slog.Float64("dischargeLimit", limit))
		if err := g.client.SetSocLimit(device.Serial, cl, limit); err != nil {
			slog.Error("unable to set discharge limit (app)", slog.String("error", err.Error()))
			return false
		} else {
			slog.Info("set discharge limit (app)", slog.String("device", device.Serial), slog.Float64("chargingLimit", cl), slog.Float64("dischargeLimit", limit))
			go g.pollParameterData(device)
			return true
		}
	}
}

func (g *GrowattAppService) poll() {
	slog.Info("start polling growatt (app)",
		slog.Int("interval", int(g.opts.PollingInterval/time.Second)),
		slog.Int("battery-details-interval", int(g.opts.BatteryDetailsPollingInterval/time.Second)),
		slog.Int("parameter-interval", int(g.opts.ParameterPollingInterval/time.Second)))

	go func() {
		for {
			for _, device := range g.devices {
				g.pollStatus(device)
			}
			<-time.After(g.opts.PollingInterval)
		}
	}()

	go func() {
		for {
			for _, device := range g.devices {
				g.pollBatteryDetails(device)
			}
			<-time.After(g.opts.BatteryDetailsPollingInterval)
		}
	}()

	go func() {
		for {
			for _, device := range g.devices {
				g.pollParameterData(device)
			}
			<-time.After(g.opts.ParameterPollingInterval)
		}
	}()
}
