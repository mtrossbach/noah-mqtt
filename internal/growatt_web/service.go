package growatt_web

import (
	"fmt"
	"log/slog"
	"noah-mqtt/internal/endpoint"
	"noah-mqtt/internal/misc"
	"noah-mqtt/pkg/models"
	"time"
)

type Options struct {
	ServerUrl       string
	Username        string
	Password        string
	PollingInterval time.Duration
}
type GrowattService struct {
	opts      Options
	client    *Client
	devices   []models.NoahDevicePayload
	endpoints []endpoint.Endpoint
}

func NewGrowattService(options Options) *GrowattService {
	return &GrowattService{
		opts:   options,
		client: newClient(options.ServerUrl, options.Username, options.Password),
	}
}

func (g *GrowattService) Login() error {
	slog.Info("logging in to growatt (web)...")
	if err := g.client.Login(); err != nil {
		return err
	}
	return nil
}

func (g *GrowattService) StartPolling() {
	g.devices = g.enumerateDevices()
	for _, e := range g.endpoints {
		e.SetDevices(g.devices)
	}

	go g.poll()
}

func (g *GrowattService) AddEndpoint(e endpoint.Endpoint) {
	g.endpoints = append(g.endpoints, e)
}

func (g *GrowattService) enumerateDevices() []models.NoahDevicePayload {
	var enumeratedDevices []models.NoahDevicePayload

	plantList, err := g.client.GetPlantList()
	if err != nil {
		slog.Error("could not get plant list", slog.String("error", err.Error()))
		return enumeratedDevices
	}

	for _, plant := range plantList {
		if devices, err := g.client.GetNoahList(misc.S2i(plant.PlantId)); err != nil {
			slog.Error("could not get plant devices", slog.String("plantId", plant.PlantId), slog.String("error", err.Error()))
		} else {
			for _, dev := range devices.Datas {

				if history, err := g.client.GetNoahHistory(dev.Sn, "", ""); err != nil {
					slog.Error("could not get device history", slog.String("device", dev.Sn), slog.String("error", err.Error()))
				} else {
					if len(history.Obj.Datas) == 0 {
						slog.Error("could not get device history, data empty", slog.String("device", dev.Sn))
					} else {
						var batCount = history.Obj.Datas[0].BatteryPackageQuantity
						var batteries []models.NoahDeviceBatteryPayload
						for i := 0; i < batCount; i++ {
							batteries = append(batteries, models.NoahDeviceBatteryPayload{
								Alias: fmt.Sprintf("BAT%d", i),
							})
						}
						d := models.NoahDevicePayload{
							PlantId:   misc.S2i(dev.PlantID),
							Serial:    dev.Sn,
							Model:     dev.DeviceModel,
							Version:   dev.Version,
							Alias:     dev.Alias,
							Batteries: batteries,
						}

						enumeratedDevices = append(enumeratedDevices, d)
					}
				}

			}
		}
	}

	return enumeratedDevices
}

func (g *GrowattService) poll() {
	historyInterval := 3 * time.Minute

	slog.Info("start polling growatt (web)",
		slog.Int("interval", int(g.opts.PollingInterval/time.Second)),
		slog.Int("history-interval", int(historyInterval/time.Second)))

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
				g.pollHistory(device)
			}
			<-time.After(historyInterval)
		}
	}()
}

func (g *GrowattService) pollStatus(device models.NoahDevicePayload) {
	if status, err := g.client.GetNoahStatus(device.PlantId, device.Serial); err != nil {
		slog.Error("could not get device data", slog.String("error", err.Error()), slog.String("device", device.Serial))
	} else {
		if totals, err := g.client.GetNoahTotals(device.PlantId, device.Serial); err != nil {
			slog.Error("could not get device totals", slog.String("error", err.Error()), slog.String("device", device.Serial))
		} else {
			batteryPower := misc.ParseFloat(status.Obj.TotalBatteryPackChargingPower)

			chargePower := 0.0
			dischargePower := 0.0
			if batteryPower < 0 {
				dischargePower = -batteryPower
			} else {
				chargePower = batteryPower
			}

			payload := models.DevicePayload{
				OutputPower:           misc.ParseFloat(status.Obj.Pac),
				SolarPower:            misc.ParseFloat(status.Obj.Ppv),
				Soc:                   misc.ParseFloat(status.Obj.TotalBatteryPackSoc),
				ChargePower:           chargePower,
				DischargePower:        dischargePower,
				BatteryNum:            len(device.Batteries),
				GenerationTotalEnergy: misc.ParseFloat(totals.Obj.EacTotal),
				GenerationTodayEnergy: misc.ParseFloat(totals.Obj.EacToday),
				WorkMode:              models.WorkModeFromString(status.Obj.WorkMode),
				Status:                models.StatusFromString(status.Obj.Status),
			}

			for _, e := range g.endpoints {
				e.PublishDeviceStatus(device, payload)
			}

		}
	}
}

func (g *GrowattService) pollHistory(device models.NoahDevicePayload) {
	if details, err := g.client.GetNoahDetails(device.PlantId, device.Serial); err != nil {
		slog.Error("could not get device details data", slog.String("error", err.Error()))
	} else {
		if len(details.Datas) != 1 {
			slog.Error("could not get device details data", slog.String("device", device.Serial))
		} else {
			detailsData := details.Datas[0]
			cl := misc.ParseFloat(detailsData.ChargingSocHighLimit)
			dl := misc.ParseFloat(detailsData.ChargingSocLowLimit)
			op := misc.ParseFloat(detailsData.DefaultPower)
			paramPayload := models.ParameterPayload{
				ChargingLimit:  &cl,
				DischargeLimit: &dl,
				OutputPower:    &op,
			}

			for _, e := range g.endpoints {
				e.PublishParameterData(device, paramPayload)
			}
		}
	}

	if history, err := g.client.GetNoahHistory(device.Serial, "", ""); err != nil {
		slog.Error("could not get device history", slog.String("error", err.Error()), slog.String("device", device.Serial))
	} else {
		if len(history.Obj.Datas) == 0 {
			slog.Error("could not get device history, data empty", slog.String("device", device.Serial))
		} else {
			historyData := history.Obj.Datas[0]

			var batteries []models.BatteryPayload
			for i := 0; i < len(device.Batteries); i++ {
				switch i {
				case 0:
					batteries = append(batteries, models.BatteryPayload{
						SerialNumber: historyData.Battery1SerialNum,
						Soc:          float64(historyData.Battery1Soc),
						Temperature:  historyData.Battery1Temp,
					})
				case 1:
					batteries = append(batteries, models.BatteryPayload{
						SerialNumber: historyData.Battery2SerialNum,
						Soc:          float64(historyData.Battery2Soc),
						Temperature:  historyData.Battery2Temp,
					})
				case 2:
					batteries = append(batteries, models.BatteryPayload{
						SerialNumber: historyData.Battery3SerialNum,
						Soc:          float64(historyData.Battery3Soc),
						Temperature:  historyData.Battery3Temp,
					})
				case 3:
					batteries = append(batteries, models.BatteryPayload{
						SerialNumber: historyData.Battery4SerialNum,
						Soc:          float64(historyData.Battery4Soc),
						Temperature:  historyData.Battery4Temp,
					})
				}
			}

			for _, e := range g.endpoints {
				e.PublishBatteryDetails(device, batteries)
			}
		}
	}
}
