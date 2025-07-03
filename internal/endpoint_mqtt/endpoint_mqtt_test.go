package endpoint_mqtt

import (
	"encoding/json"
	"fmt"
	"math"
	"nexa-mqtt/internal/homeassistant"
	"nexa-mqtt/pkg/models"
	"sync"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ----- Mocks --------------------------------------------------------------

// MockToken implements mqtt.Token
type MockToken struct {
	mock.Mock
	done chan struct{}
}

func NewMockToken() *MockToken {
	done := make(chan struct{})
	close(done) // sofort abgeschlossen
	return &MockToken{done: done}
}

func (m *MockToken) Wait() bool                     { return true }
func (m *MockToken) WaitTimeout(time.Duration) bool { return true }
func (t *MockToken) Done() <-chan struct{}          { return t.done }
func (t *MockToken) Error() error {
	args := t.Called("Error")
	return args.Error(0)
}

// MockMqttClient implements mqtt.Client
type MockMqttClient struct {
	mock.Mock
	mqtt.Client
}

func (m *MockMqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	args := m.Called(topic, qos, retained, payload)
	return args.Get(0).(mqtt.Token)
}

func (m *MockMqttClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	args := m.Called(topic, qos, callback)
	return args.Get(0).(mqtt.Token)
}

func (m *MockMqttClient) Unsubscribe(topics ...string) mqtt.Token {
	ifaceArgs := make([]interface{}, len(topics))
	for i, v := range topics {
		ifaceArgs[i] = v
	}
	args := m.Called(ifaceArgs...)
	return args.Get(0).(mqtt.Token)
}

// MockMqttMessage implements mqtt.Message
type MockMqttMessage struct {
	mock.Mock
	mqtt.Message
}

func (m *MockMqttMessage) Payload() []byte {
	args := m.Called()
	return args.Get(0).([]byte)
}

// MockParameterApplier implements endpoint.ParameterApplier
type MockParameterApplier struct {
	mock.Mock
}

func (p *MockParameterApplier) SetOutputPowerW(device models.NoahDevicePayload, mode models.WorkMode, power float64) bool {
	args := p.Called(device, mode, power)
	return args.Get(0).(bool)
}

func (p *MockParameterApplier) SetChargingLimits(device models.NoahDevicePayload, chargingLimit float64, dischargeLimit float64) bool {
	args := p.Called(device, chargingLimit, dischargeLimit)
	return args.Get(0).(bool)
}

// MockHaClient implements homeassistant.HaClient
type MockHaClient struct {
	mock.Mock
}

func (s *MockHaClient) SetDevices(devices []homeassistant.DeviceInfo) {
	s.Called(devices)
}

// ----- Test functions -----------------------------------------------------

func TestNewEndpoint(t *testing.T) {
	endpoint := NewEndpoint(Options{})

	assert.Equal(t, models.EmptyParameterPayload(), endpoint.lastParameter)
	assert.Equal(t, models.ParameterPayload{}, endpoint.newParameter)
}

func TestSetDevices(t *testing.T) {
	mockClient := new(MockMqttClient)
	mockToken := NewMockToken()
	haClient := &MockHaClient{}
	endpoint := &Endpoint{
		opts: Options{
			MqttClient:  mockClient,
			TopicPrefix: "test",
			HaClient:    haClient,
		},
	}

	devices1 := []models.NoahDevicePayload{
		{Serial: "device123", Batteries: []models.NoahDeviceBatteryPayload{{Alias: "A"}, {Alias: "B"}}},
		{Serial: "device234", Batteries: []models.NoahDeviceBatteryPayload{{Alias: "C"}}},
	}

	mockClient.On(
		"Subscribe",
		"test/device123/parameters/set",
		byte(0),
		mock.AnythingOfType("mqtt.MessageHandler"),
	).Return(mockToken)
	mockClient.On(
		"Subscribe",
		"test/device234/parameters/set",
		byte(0),
		mock.AnythingOfType("mqtt.MessageHandler"),
	).Return(mockToken)

	haClient.On(
		"SetDevices",
		[]homeassistant.DeviceInfo{
			{
				SerialNumber:          "device123",
				StateTopic:            "test/device123",
				ParameterStateTopic:   "test/device123/parameters",
				ParameterCommandTopic: "test/device123/parameters/set",
				Batteries: []homeassistant.BatteryInfo{
					{
						Alias:      "A",
						StateTopic: "test/device123/BAT0",
					},
					{
						Alias:      "B",
						StateTopic: "test/device123/BAT1",
					},
				},
			},
			{
				SerialNumber:          "device234",
				StateTopic:            "test/device234",
				ParameterStateTopic:   "test/device234/parameters",
				ParameterCommandTopic: "test/device234/parameters/set",
				Batteries: []homeassistant.BatteryInfo{
					{
						Alias:      "C",
						StateTopic: "test/device234/BAT0",
					},
				},
			},
		},
	)

	endpoint.SetDevices(devices1)

	mockClient.AssertExpectations(t)
	haClient.AssertExpectations(t)

	devices2 := []models.NoahDevicePayload{
		{Serial: "device345", Batteries: []models.NoahDeviceBatteryPayload{}},
	}

	mockClient.On(
		"Unsubscribe",
		"test/device123/parameters/set",
	).Return(mockToken)
	mockClient.On(
		"Unsubscribe",
		"test/device234/parameters/set",
	).Return(mockToken)
	mockClient.On(
		"Subscribe",
		"test/device345/parameters/set",
		byte(0),
		mock.AnythingOfType("mqtt.MessageHandler"),
	).Return(mockToken)

	haClient.On(
		"SetDevices",
		[]homeassistant.DeviceInfo{
			{
				SerialNumber:          "device345",
				StateTopic:            "test/device345",
				ParameterStateTopic:   "test/device345/parameters",
				ParameterCommandTopic: "test/device345/parameters/set",
				Batteries:             nil,
			},
		},
	)

	endpoint.SetDevices(devices2)

	mockClient.AssertExpectations(t)
	haClient.AssertExpectations(t)
}

func TestPublishDeviceStatus_Success(t *testing.T) {
	mockClient := new(MockMqttClient)
	mockToken := NewMockToken()

	mockClient.On(
		"Publish",
		"test/device123",
		byte(0),
		false,
		`{"output_w":0,"solar_w":0,"soc":0,"charge_w":0,"discharge_w":0,"battery_num":0,"generation_total_kwh":0,"generation_today_kwh":0}`,
	).Return(mockToken)

	endpoint := &Endpoint{
		opts: Options{
			MqttClient:  mockClient,
			TopicPrefix: "test",
		},
	}

	device := models.NoahDevicePayload{Serial: "device123"}
	status := models.DevicePayload{}

	endpoint.PublishDeviceStatus(device, status)

	mockClient.AssertExpectations(t)
}

func TestPublishDeviceStatus_Fail(t *testing.T) {
	mockClient := new(MockMqttClient)

	endpoint := &Endpoint{
		opts: Options{
			MqttClient:  mockClient,
			TopicPrefix: "test",
		},
	}

	device := models.NoahDevicePayload{Serial: "device123"}
	status := models.DevicePayload{OutputPower: math.NaN()}

	endpoint.PublishDeviceStatus(device, status)

	mockClient.AssertExpectations(t)
}

func TestPublishBatteryDetails_Success(t *testing.T) {
	mockClient := new(MockMqttClient)
	mockToken := NewMockToken()

	mockClient.On(
		"Publish",
		"test/device123/BAT0",
		byte(0),
		false,
		`{"serial":"","soc":0,"temp":0}`,
	).Return(mockToken)

	endpoint := &Endpoint{
		opts: Options{
			MqttClient:  mockClient,
			TopicPrefix: "test",
		},
	}

	device := models.NoahDevicePayload{Serial: "device123"}
	details := []models.BatteryPayload{{}}

	endpoint.PublishBatteryDetails(device, details)

	mockClient.AssertExpectations(t)
}

func TestPublishBatteryDetails_SuccessMult(t *testing.T) {
	mockClient := new(MockMqttClient)
	mockToken := NewMockToken()

	mockClient.On(
		"Publish",
		"test/device123/BAT0",
		byte(0),
		false,
		`{"serial":"E","soc":20,"temp":0}`,
	).Return(mockToken)
	mockClient.On(
		"Publish",
		"test/device123/BAT1",
		byte(0),
		false,
		`{"serial":"F","soc":30,"temp":0}`,
	).Return(mockToken)
	mockClient.On(
		"Publish",
		"test/device123/BAT2",
		byte(0),
		false,
		`{"serial":"G","soc":40,"temp":0}`,
	).Return(mockToken)

	endpoint := &Endpoint{
		opts: Options{
			MqttClient:  mockClient,
			TopicPrefix: "test",
		},
	}

	device := models.NoahDevicePayload{Serial: "device123"}
	details := []models.BatteryPayload{{SerialNumber: "E", Soc: 20}, {SerialNumber: "F", Soc: 30}, {SerialNumber: "G", Soc: 40}}

	endpoint.PublishBatteryDetails(device, details)

	mockClient.AssertExpectations(t)
}

func TestPublishBatteryDetails_Fail(t *testing.T) {
	mockClient := new(MockMqttClient)

	endpoint := &Endpoint{
		opts: Options{
			MqttClient:  mockClient,
			TopicPrefix: "test",
		},
	}

	device := models.NoahDevicePayload{Serial: "device123"}
	details := []models.BatteryPayload{{Soc: math.NaN()}}

	endpoint.PublishBatteryDetails(device, details)

	mockClient.AssertExpectations(t)
}

func TestPublishParameterData_Success(t *testing.T) {
	mockClient := new(MockMqttClient)
	mockToken := NewMockToken()

	param := models.EmptyParameterPayload()
	json, err := json.Marshal(param)
	assert.Nil(t, err)

	mockClient.On(
		"Publish",
		"test/device123/parameters",
		byte(0),
		false,
		string(json),
	).Return(mockToken)

	endpoint := &Endpoint{
		opts: Options{
			MqttClient:  mockClient,
			TopicPrefix: "test",
		},
	}

	device := models.NoahDevicePayload{Serial: "device123"}

	endpoint.PublishParameterData(device, param)

	mockClient.AssertExpectations(t)
}

func TestPublishParameterData_Fail(t *testing.T) {
	mockClient := new(MockMqttClient)

	endpoint := &Endpoint{
		opts: Options{
			MqttClient:  mockClient,
			TopicPrefix: "test",
		},
	}

	device := models.NoahDevicePayload{Serial: "device123"}
	param := models.EmptyParameterPayload()
	*param.ChargingLimit = math.NaN()

	endpoint.PublishParameterData(device, param)

	mockClient.AssertExpectations(t)
}

func Test_parametersSubscription_NoApplier(t *testing.T) {
	mockClient := new(MockMqttClient)
	endpoint := NewEndpoint(Options{MqttClient: mockClient, TopicPrefix: "test"})

	device := models.NoahDevicePayload{Serial: "device123"}
	f1 := endpoint.parametersSubscription(device)

	mockMqttMessage := MockMqttMessage{}

	f1(mockClient, &mockMqttMessage)

	mockMqttMessage.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

func setup_parametersSubscription() (*MockToken, *MockMqttClient, *MockParameterApplier, *Endpoint, models.NoahDevicePayload, func(client mqtt.Client, message mqtt.Message)) {
	mockToken := NewMockToken()
	mockClient := new(MockMqttClient)
	mockApplier := MockParameterApplier{}
	endpoint := NewEndpoint(Options{MqttClient: mockClient, TopicPrefix: "test"})
	endpoint.SetParameterApplier(&mockApplier)
	device := models.NoahDevicePayload{Serial: "device123"}
	f1 := endpoint.parametersSubscription(device)
	return mockToken, mockClient, &mockApplier, endpoint, device, f1
}

func Test_parametersSubscription_InvalidPayload(t *testing.T) {
	_, mockClient, mockApplier, _, _, f1 := setup_parametersSubscription()

	mockMqttMessage := MockMqttMessage{}
	mockMqttMessage.On("Payload").
		Return([]byte(`{"charging_limit":"invalid string"}`))

	f1(mockClient, &mockMqttMessage)

	mockMqttMessage.AssertExpectations(t)
	mockApplier.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

func Test_parametersSubscription_ChargingLimit(t *testing.T) {
	mockToken, mockClient, mockApplier, endpoint, device, f1 := setup_parametersSubscription()
	empty := models.EmptyParameterPayload()

	mockMqttMessage := MockMqttMessage{}
	mockMqttMessage.On("Payload").
		Return([]byte(`{"charging_limit":90}`))

	var wg sync.WaitGroup

	mockApplier.On("SetChargingLimits", device, 90.0, *empty.DischargeLimit).
		Run(func(args mock.Arguments) {
			wg.Done()
		}).
		Return(true)

	mockClient.On(
		"Publish",
		"test/device123/parameters",
		byte(0),
		false,
		fmt.Sprintf(`{"charging_limit":90,"discharge_limit":%v,"default_output_w":%v,"default_mode":"%v"}`,
			*empty.DischargeLimit, *empty.DefaultACCouplePower, *empty.DefaultMode),
	).Run(func(args mock.Arguments) {
		wg.Done()
	}).Return(mockToken)

	wg.Add(2)
	f1(mockClient, &mockMqttMessage)
	wg.Wait()

	mockMqttMessage.AssertExpectations(t)
	mockApplier.AssertExpectations(t)
	mockClient.AssertExpectations(t)

	assert.Equal(t, models.ParameterPayload{}, endpoint.newParameter)
}

func Test_parametersSubscription_ChargingAndDischargeLimit(t *testing.T) {
	mockToken, mockClient, mockApplier, endpoint, device, call_parametersSubscription := setup_parametersSubscription()
	empty := models.EmptyParameterPayload()

	mockMqttMessage1 := MockMqttMessage{}
	mockMqttMessage1.On("Payload").
		Return([]byte(`{"charging_limit":95}`))

	mockMqttMessage2 := MockMqttMessage{}
	mockMqttMessage2.On("Payload").
		Return([]byte(`{"discharge_limit":5}`))

	var wg sync.WaitGroup

	mockApplier.On("SetChargingLimits", device, 95.0, 5.0).
		Run(func(args mock.Arguments) {
			wg.Done()
		}).
		Return(true)

	mockClient.On(
		"Publish",
		"test/device123/parameters",
		byte(0),
		false,
		fmt.Sprintf(`{"charging_limit":95,"discharge_limit":5,"default_output_w":%v,"default_mode":"%v"}`,
			*empty.DefaultACCouplePower, *empty.DefaultMode),
	).Run(func(args mock.Arguments) {
		wg.Done()
	}).Return(mockToken)

	wg.Add(2)
	call_parametersSubscription(mockClient, &mockMqttMessage1)
	call_parametersSubscription(mockClient, &mockMqttMessage2)
	wg.Wait()

	mockMqttMessage1.AssertExpectations(t)
	mockApplier.AssertExpectations(t)
	mockClient.AssertExpectations(t)

	assert.Equal(t, models.ParameterPayload{}, endpoint.newParameter)
}

func Test_parametersSubscription_ChargingLimitAndMode(t *testing.T) {
	mockToken, mockClient, mockApplier, endpoint, device, call_parametersSubscription := setup_parametersSubscription()
	empty := models.EmptyParameterPayload()

	mockMqttMessage1 := MockMqttMessage{}
	mockMqttMessage1.On("Payload").
		Return([]byte(`{"charging_limit":75}`))

	mockMqttMessage2 := MockMqttMessage{}
	mockMqttMessage2.On("Payload").
		Return([]byte(`{"default_mode":"battery_first"}`))

	var wg sync.WaitGroup

	mockApplier.On("SetChargingLimits", device, 75.0, *empty.DischargeLimit).
		Run(func(args mock.Arguments) {
			wg.Done()
		}).
		Return(true)

	mockApplier.On("SetOutputPowerW", device, models.WorkMode("battery_first"), *empty.DefaultACCouplePower).
		Run(func(args mock.Arguments) {
			wg.Done()
		}).
		Return(true)

	mockClient.On(
		"Publish",
		"test/device123/parameters",
		byte(0),
		false,
		fmt.Sprintf(`{"charging_limit":75,"discharge_limit":%v,"default_output_w":%v,"default_mode":"battery_first"}`,
			*empty.DischargeLimit, *empty.DefaultACCouplePower),
	).Run(func(args mock.Arguments) {
		wg.Done()
	}).Return(mockToken)

	wg.Add(3)
	call_parametersSubscription(mockClient, &mockMqttMessage1)
	call_parametersSubscription(mockClient, &mockMqttMessage2)
	wg.Wait()

	mockMqttMessage1.AssertExpectations(t)
	mockApplier.AssertExpectations(t)
	mockClient.AssertExpectations(t)

	assert.Equal(t, models.ParameterPayload{}, endpoint.newParameter)
}
