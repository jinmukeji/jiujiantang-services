package handler

import (
	mlog "github.com/jinmukeji/go-pkg/v2/log"
)

var (
	// log is the package global logger
	log = mlog.StandardLogger()
)

func init() {
	log.Debugln("jinmuid service is initialized.")
}
