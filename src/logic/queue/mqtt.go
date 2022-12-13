package queue

import (
	"context"
	"ehdw/smartiko-test/src/logic/service"
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

type queue struct {
	client mqtt.Client
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

var f2 mqtt.MessageHandler = func(c mqtt.Client, m mqtt.Message) {
	fmt.Printf("Finally: %s, %s", m.Topic(), m.Payload())
}

func NewQueue() *queue {
	o := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("device_proc").SetDefaultPublishHandler(f)
	c := mqtt.NewClient(o)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	logrus.Info("successfully initialised the MQTT client")
	return &queue{client: c}
}

func (q *queue) Stop() {
	if token := q.client.Unsubscribe("go-mqtt/sample"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	q.client.Disconnect(250)
	logrus.Info("Disconnected from the message queue")
}

func (q *queue) RegisterDevices(ctx context.Context, f service.MessageProcessingFunc, deviceNames ...string) error {
	var wrapper mqtt.MessageHandler = func(c mqtt.Client, m mqtt.Message) {
		f(m.Topic(), string(m.Payload()))
	}
	for _, device := range deviceNames {
		if token := q.client.Subscribe(device+"/out/data", 1, wrapper); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}
	return nil
}
