package mutex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testJob struct{}

func (t *testJob) Run() {
}

func (t *testJob) Slug() string {
	return "testJob"
}

func (t *testJob) IsMutexJob() {}

func TestJobWrapper_Job(t *testing.T) {
	_, ok1 := interface{}(struct{}{}).(Job)
	assert.False(t, ok1)

	_, ok2 := interface{}(&testJob{}).(Job)
	assert.True(t, ok2)
}
