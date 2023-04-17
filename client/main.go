package client

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/igzhang/ops_job/pkg"
	jsoniter "github.com/json-iterator/go"
)

const (
	ServerCallbackTopic = "server"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type RunCmdResult struct {
	MsgID     uint16
	ClientID  string
	IsSuccess bool
	Stdout    string
	Stderr    string
}

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

	_, err := pkg.NewMQTTClient(mqttAddr, clientID, clientSubscribeCallback)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("running success!")
	<-sigChannel
}

func clientSubscribeCallback(client MQTT.Client, msg MQTT.Message) {
	recvMsg := string(msg.Payload())
	log.Printf("receive msg: %s", recvMsg)

	cmd := exec.Command("bash", "-c", recvMsg)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runCMDResult := RunCmdResult{MsgID: msg.MessageID(), ClientID: msg.Topic(), IsSuccess: false}
	if err := cmd.Run(); err == nil {
		runCMDResult.IsSuccess = true
	}
	runCMDResult.Stdout = stdout.String()
	runCMDResult.Stderr = stderr.String()

	askMsgBytes, err := json.Marshal(&runCMDResult)
	if err != nil {
		log.Printf("json err: %s", err.Error())
	}

	if token := client.Publish(ServerCallbackTopic, 2, false, askMsgBytes); token.Wait() && token.Error() != nil {
		log.Printf("ask server cmd error: %s", token.Error().Error())
	}
}
