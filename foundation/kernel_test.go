package foundation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type application struct{}

func (a *application) Run() error {
	return nil
}

type provider struct{}

func (p *provider) Bootstrap(Application) error {
	return nil
}

func (p *provider) Terminate(Application) error {
	return nil
}

func TestKernel(t *testing.T) {
	app := &application{}
	k := NewKernel(app)
	k.Register(&provider{}, &provider{})
	assert.NoError(t, k.Run())
}
