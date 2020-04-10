package rest

import (
	"fmt"
	"time"

	mlog "github.com/jinmukeji/go-pkg/log"
	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
)

var log = mlog.StandardLogger()

// ContextLogger 打印ctx中的cid
func ContextLogger(ctx iris.Context) *logrus.Entry {
	cid := GetCidFromContext(ctx)
	return log.WithField(ContextCidKey, cid)
}

// LogMiddleware Log中间件
func LogMiddleware(ctx iris.Context) {
	var startTime, endTime time.Time
	var latency time.Duration

	startTime = time.Now()
	ctx.Next()

	//no time.Since in order to format it well after
	endTime = time.Now()
	latency = endTime.Sub(startTime)

	logRequest(ctx, latency)
}

// logRequest 打印请求
func logRequest(ctx iris.Context, latency time.Duration) {
	//all except latency to string
	var status int
	var ip, userAgent, referrer, method, path, query string

	req := ctx.Request()

	cid := GetCidFromContext(ctx)
	ip = ctx.RemoteAddr()
	referrer = req.Referer()
	userAgent = req.UserAgent()
	method = req.Method
	path = req.URL.EscapedPath()
	query = req.URL.RawQuery
	status = ctx.GetStatusCode()

	l := log.WithFields(logrus.Fields{
		ContextCidKey:    cid,
		"req.ip":         ip,
		"req.user_agent": userAgent,
		"req.referrer":   referrer,
		"req.path":       fmt.Sprintf("%s %s", method, path),
		"req.query":      query,
		"req.status":     status,
		"req.latency":    latency.String(),
	})

	msg := GetContextMessage(ctx)
	l.Info(msg)
}
