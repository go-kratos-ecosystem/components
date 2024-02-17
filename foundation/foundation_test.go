package foundation

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-kratos-ecosystem/components/v2/coordinator"
	"github.com/go-kratos-ecosystem/components/v2/event"
)

var (
	wg sync.WaitGroup
	c  = make(chan string, 8)
)

type listener struct{}

var _ event.Listener = (*listener)(nil)

func (l *listener) Listen() []event.Event {
	return []event.Event{
		BootstrapName,
		ShutdownName,
	}
}

func (l *listener) Handle(event event.Event, data interface{}) {
	if event.String() == BootstrapName {
		if _, ok := data.(*BootstrapEvent); ok {
			c <- "bootstrap done from listener"
			wg.Done()
		}
	}

	if event.String() == ShutdownName {
		if _, ok := data.(*ShutdownEvent); ok {
			c <- "shutdown done from listener"
			wg.Done()
		}
	}
}

func TestBootAndShut(t *testing.T) {
	var (
		m    = coordinator.NewManager()
		d    = event.NewDispatcher()
		boot = NewBootstrap(WithManager(m), WithDispatcher(d))
		shut = NewShutdown(WithManager(m), WithDispatcher(d))
	)
	wg.Add(4)
	assert.Equal(t, 0, len(c))

	go func() {
		defer wg.Done()
		if <-m.Until(BootstrapName).Done(); true {
			c <- "bootstrap done"
		}
	}()

	go func() {
		defer wg.Done()
		if <-m.Until(ShutdownName).Done(); true {
			c <- "shutdown done"
		}
	}()

	d.AddListener(&listener{})

	assert.NoError(t, boot(context.Background()))
	assert.NoError(t, shut(context.Background()))

	wg.Wait()

	assert.Equal(t, 4, len(c))
}
