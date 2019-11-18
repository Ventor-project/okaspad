package mqtt

import (
	"errors"
	"github.com/daglabs/btcd/apiserver/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// client is an instance of the MQTT client, in case we have an active connection
var client mqtt.Client

// GetClient returns an instance of the MQTT client, in case we have an active connection
func GetClient() (mqtt.Client, error) {
	if client == nil {
		return nil, errors.New("MQTT is not connected")
	}
	return client, nil
}

// Connect initiates a connection to the MQTT server, if defined
func Connect() error {
	cfg := config.ActiveConfig()
	if cfg.MQTTBrokerAddress == "" {
		// MQTT broker not defined -- nothing to do
		return nil
	}

	options := mqtt.NewClientOptions()
	options.AddBroker(cfg.MQTTBrokerAddress)
	options.SetUsername(cfg.MQTTUser)
	options.SetPassword(cfg.MQTTPassword)
	options.SetAutoReconnect(true)

	newClient := mqtt.NewClient(options)
	if token := newClient.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	client = newClient

	return nil
}

// Close closes the connection to the MQTT server, if previously connected
func Close() {
	if client == nil {
		return
	}
	client.Disconnect(250)
	client = nil
}