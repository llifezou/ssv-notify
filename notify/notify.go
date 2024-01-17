package notify

import (
	"errors"
	"fmt"
	"log"
)

type Sender interface {
	Send(msg string) error
	Platform() string
}

type Notify struct {
	Senders []Sender
}

func NewNotify(senders ...Sender) (*Notify, error) {
	n := &Notify{
		Senders: senders,
	}
	// errs := n.Send("Monitoring service startup test message")
	// if len(errs) != 0 {
	// 	return nil, errors.New("NewNotify failed")
	// }
	return n, nil
}

func (n *Notify) Send(msg string) []error {
	var errs []error
	for _, sender := range n.Senders {
		err := sender.Send(msg)
		if err != nil {
			name := sender.Platform()
			werr := fmt.Errorf("Notification failed! platform: %s, err: %s \n", name, err)
			errs = append(errs, werr)
			log.Println(err.Error())
		}
	}
	return errs
}
