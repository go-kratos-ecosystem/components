package pubsub

import (
	"fmt"
	"testing"
)

type loginEvent struct {
	UserID string
}

func TestTopic(*testing.T) {
	topic := NewTopic[*loginEvent]()

	topic.Subscribe(func(msg *loginEvent) error {
		fmt.Println(msg.UserID)
		return nil
	})

	sub2 := topic.Subscribe(func(_ *loginEvent) error {
		return fmt.Errorf("error")
	})

	err := topic.Publish(&loginEvent{UserID: "123"})
	if err != nil {
		return
	}

	sub2.Unsubscribe()

	if err = topic.Publish(&loginEvent{UserID: "456"}, PublishAsync(), PublishSkipErrors()); err != nil {
		return
	}

	topic2 := NewTopic[*loginEvent]()
	topic2.Subscribe(func(msg *loginEvent) error {
		fmt.Println("Topic 2", msg.UserID)
		return nil
	})

	if err = topic2.Publish(&loginEvent{UserID: "789"}); err != nil {
		return
	}
}
