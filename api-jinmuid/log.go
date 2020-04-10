package main

import (
	"os"

	"github.com/jinmukeji/gf-api2/api-v2/config"
	mlog "github.com/jinmukeji/go-pkg/log"
	"github.com/sirupsen/logrus"
)

var (
	// log is the package global logger
	log *mlog.Logger
)

func init() {
	log = mlog.StandardLogger()
	setupLogger(log, config.ServiceName)
}

func setupLogger(logger *mlog.Logger, svc string) {
	// Setup formatter
	if os.Getenv("LOG_FORMAT") == "logstash" {
		formatter := mlog.NewLogstashFormatter(logrus.Fields{
			"svc": svc,
		})
		log.SetFormatter(formatter)
	}

	// Setup Log level
	if lvl := os.Getenv("LOG_LEVEL"); lvl != "" {
		if level, err := logrus.ParseLevel(lvl); err != nil {
			log.Fatal(err.Error())
		} else {
			log.SetLevel(level)
		}
	}
}
