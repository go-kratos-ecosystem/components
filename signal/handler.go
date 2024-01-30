package signal

import (
	"os"
)

type Handler interface {
	Listen() []os.Signal
	Handle(os.Signal)
}

type asyncFeature interface {
	Async() bool
}

type AsyncFeature struct{}

func (*AsyncFeature) Async() bool {
	return true
}
