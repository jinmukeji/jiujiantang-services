package rest

import (
	"time"

	mlog "github.com/jinmukeji/go-pkg/v2/log"
	"github.com/sirupsen/logrus"
)

var (
	// log is the package global logger
	log = mlog.StandardLogger()
)

func logFunc(now time.Time, latency time.Duration, status, ip, method, path string, message interface{}, headerMessage interface{}) {
	l := log.WithFields(logrus.Fields{
		"latency": latency,
		"status":  status,
		"ip":      ip,
		"method":  method,
		"path":    path,
	})

	if message != nil {
		l.Info(message)
	} else {
		l.Info()
	}
}
