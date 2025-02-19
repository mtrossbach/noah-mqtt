package endpoint_mqtt

import "fmt"

func deviceStateTopic(topicPrefix string, serialNumber string) string {
	return fmt.Sprintf("%s/%s", topicPrefix, serialNumber)
}

func stateTopicBattery(topicPrefix string, serialNumber string, index int) string {
	return fmt.Sprintf("%s/%s/BAT%d", topicPrefix, serialNumber, index)
}

func parameterStateTopic(topicPrefix string, serialNumber string) string {
	return fmt.Sprintf("%s/%s/parameters", topicPrefix, serialNumber)
}

func parameterCommandTopic(topicPrefix string, serialNumber string) string {
	return fmt.Sprintf("%s/%s/parameters/set", topicPrefix, serialNumber)
}
