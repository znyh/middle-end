package kafka

import (
	"fmt"
	"testing"
)

func TestProducer(t *testing.T) {
	fmt.Println("producer...")

	producer, err := NewProducer([]string{_addr})
	if err != nil {
		fmt.Println("failed to new producer,err:", err)
		return
	}
	defer producer.Close()

	err = producer.Producer(Message{
		Topic: "test",
		Value: []byte("wo shi yi tiao yu"),
	})

	if err != nil {
		fmt.Println("failed to send message,err:", err)
		return
	}
}
