package pkg

import (
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func NewMQTTClient(broker string, clientID string) (*MQTT.Client, error) {
	opts := MQTT.NewClientOptions().AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetOnConnectHandler(func(c MQTT.Client) {
		log.Println("mqtt connection established")
	})

	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return &c, nil
}

func SubscribeTopic(mqttClient *MQTT.Client, topic string, callback MQTT.MessageHandler) error {
	if token := (*mqttClient).Subscribe(topic, 2, callback); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	log.Printf("subscribed to %s success!", topic)
	return nil
}
