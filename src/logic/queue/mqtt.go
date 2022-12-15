package queue

import (
	"context"
	"ehdw/smartiko-test/src/config"
	"ehdw/smartiko-test/src/logic/service"
	"ehdw/smartiko-test/src/util"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

type queue struct {
	client        mqtt.Client
	activeListens *util.Set[string]
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func NewQueue() *queue {
	broker := fmt.Sprintf("tcp://%s:1883", config.Config().MosquittoHost)
	o := mqtt.NewClientOptions().AddBroker(broker).SetClientID("device_proc").SetDefaultPublishHandler(f)
	c := mqtt.NewClient(o)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		//panic(token.Error())
		logrus.Warn("failed to connect to mqtt: " + token.Error().Error())
		return nil
	}
	logrus.Info("successfully initialised the MQTT client")
	return &queue{client: c, activeListens: util.NewSet[string]()}
}

func (q *queue) Stop() {
	for _, s := range q.activeListens.ToList() {
		if token := q.client.Unsubscribe(s); token.Wait() && token.Error() != nil {
			logrus.Warn("Failed to unsubscribe for " + s + " , err: " + token.Error().Error())
		}
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
		q.activeListens.Insert(device + "/out/data")
	}
	return nil
}

func (q *queue) UnsubscribeDevices(ctx context.Context, deviceNames ...string) error {
	for _, d := range deviceNames {
		if q.activeListens.Has(d) {
			if token := q.client.Unsubscribe(d); token.Wait() && token.Error() != nil {
				logrus.Warn("Failed to unsubscribe for " + d + " , err: " + token.Error().Error())
			}
		}
	}
	return nil
}
