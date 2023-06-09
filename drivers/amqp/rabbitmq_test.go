package amqp

import (
	"github.com/hq2005001/modules/logger"
	"testing"
	"time"
)

func TestRabbitMQ_SendMessage(t *testing.T) {
	New(&RabbitMQConf{
		Username: "admin",
		Password: "admin",
		Host:     "192.168.123.87",
		Port:     5672,
	}, logger.New(&logger.Conf{Level: "debug"})).SendMessage("test", "test")
}

func TestRabbitMQ_ConsumeMessage(t *testing.T) {
	client := New(&RabbitMQConf{
		Username: "admin",
		Password: "admin",
		Host:     "192.168.123.87",
		Port:     5672,
	}, logger.New(&logger.Conf{Level: "debug"}))
	go func() {
		time.Sleep(time.Second * 10)
		client.done <- true
	}()
	err := client.ConsumeMessage("test", func(data []byte) error {
		t.Log(string(data))
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}
