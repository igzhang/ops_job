package pkg

import (
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const qos = 2

func NewMQTTClient(broker string, clientID string, callback MQTT.MessageHandler) (*MQTT.Client, error) {
	opts := MQTT.NewClientOptions().AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetOnConnectHandler(func(c MQTT.Client) {
		log.Println("mqtt connection established")
		if token := c.Subscribe(clientID, qos, callback); token.Wait() && token.Error() != nil {
			log.Printf("subscribe err: %s", token.Error())
		}
		log.Printf("subscribed to %s success!", clientID)
	})

	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return &c, nil
}
