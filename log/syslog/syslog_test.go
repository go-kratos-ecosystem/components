//go:build !windows && !plan9

package syslog

import (
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

func TestSyslogLogger(t *testing.T) {
	t.Skip("skip syslog test")

	logger, err := New(&Config{
		Network: "udp",
		Addr:    "192.168.8.92:30732",
		Tag:     "test",
	})
	defer logger.Close()

	if err != nil {
		t.Fatal(err)
	}

	err = logger.Log(log.LevelDebug, "test", "test")
	if err != nil {
		t.Fatal(err)
	}
}
