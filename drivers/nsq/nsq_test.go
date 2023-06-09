package nsq

import (
	"context"
	"fmt"
	"github.com/nsqio/go-nsq"
	"testing"
)

func TestProduce(t *testing.T) {
	err := Produce(&NsqConfig{
		Lookups: []string{"127.0.0.1:4161"},
		NSQD:    "127.0.0.1:4150",
	},
		"test1", []byte("延迟1分钟发送"))
	if err != nil {
		t.Error(err.Error())
	}
}

type testHandler struct {
}

func (t testHandler) HandleMessage(message *nsq.Message) error {
	fmt.Println(string(message.Body))
	return nil
}

// Topic 消息主题
func (t testHandler) Topic() string {
	return "test1"
}

// Channel 消息通道
func (t testHandler) Channel() string {
	return "test1"
}

func (t testHandler) Workers() int {
	return 1
}

func TestConsume(t *testing.T) {
	Consume(&NsqConfig{
		Lookups: []string{"127.0.0.1:4161"},
		NSQD:    "127.0.0.1:4150",
	}, context.TODO(), &testHandler{})
	select {}
}

func TestConsume1(t *testing.T) {

	Consume(&NsqConfig{
		Lookups: []string{"127.0.0.1:4161"},
		NSQD:    "127.0.0.1:4150",
	}, context.TODO(), &testHandler{})
	select {}
}
