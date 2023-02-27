package client

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/igzhang/ops_job/pkg"
)

const (
	ExecSuccess         = 0
	ExecFailure         = 1
	ServerCallbackTopic = "server"
)

func Run() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	mqttAddr := os.Getenv("mqtt")
	if len(mqttAddr) == 0 {
		log.Fatalln("env: mqtt is not specified!")
	}

	clientID := os.Getenv("id")
	if len(clientID) == 0 {
		log.Fatalln("env: id is not specified!")
	}

	mqttClient, err := pkg.NewMQTTClient(mqttAddr, clientID)
	if err != nil {
		log.Fatalln(err.Error())
	}

	if err = pkg.SubscribeTopic(mqttClient, clientID, clientSubscribeCallback); err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("running success!")
	<-sigChannel
}

func clientSubscribeCallback(client MQTT.Client, msg MQTT.Message) {
	recvMsg := string(msg.Payload())
	log.Printf("receive msg: %s", recvMsg)

	cmd := exec.Command(recvMsg)
	var RunCmdResult int
	if err := cmd.Run(); err != nil {
		RunCmdResult = ExecSuccess
	} else {
		RunCmdResult = ExecFailure
	}
	askMsg := fmt.Sprintf("%v:%d", msg.MessageID(), RunCmdResult)

	if token := client.Publish(ServerCallbackTopic, 2, false, askMsg); token.Wait() && token.Error() != nil {
		log.Printf("ask server cmd error: %s", token.Error().Error())
	}
}
