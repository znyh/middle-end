package kafka

import (
	"fmt"
	"testing"
)

const (
	//_addr = "172.0.0.1:9092"
	_addr = "192.168.131.131:9092"
)

var (
	handler = func(msg []byte, args ...interface{}) {
		if len(args) != 1 {
			fmt.Println("arg is empty,")
			return
		}
		value, ok := args[0].(int)
		if !ok {
			fmt.Println("args 1 not type int")
			return
		}
		fmt.Println("value:", value)
		fmt.Println("msg:", string(msg))
	}
)

func TestConsumer(t *testing.T) {
	fmt.Println("consumer...")
	consumer, err := NewConsumer([]string{_addr}, "zhuma")
	if err != nil {
		fmt.Println("failed to new consumer,err:", err)
		return
	}
	defer consumer.Close()
	err = consumer.Consume(map[string]Handler{
		"test": Handler{
			Run:  handler,
			Args: []interface{}{1},
		},
	})
}
