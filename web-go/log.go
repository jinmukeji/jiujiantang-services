package main

import (
	"os"

	mlog "github.com/jinmukeji/go-pkg/v2/log"
	"github.com/jinmukeji/jiujiantang-services/web-go/config"
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
