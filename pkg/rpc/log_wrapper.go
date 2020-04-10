package rpc

import (
	"context"
	"time"

	mlog "github.com/jinmukeji/go-pkg/log"
	"github.com/micro/go-micro/server"
	"github.com/sirupsen/logrus"
)

var (
	// log is the package global logger
	log = mlog.StandardLogger()
)

const (
	logCidKey     = "cid"
	logLatencyKey = "latency"
	logRpcCallKey = "rpc.call"

	// rpcMetadata   = "[RPC METADATA]"
	rpcFailed = "[RPC ERR]"
	rpcOk     = "[RPC OK]"
)

// ContextLogger 打印ctx中的cid
func ContextLogger(ctx context.Context) *logrus.Entry {
	cid := ContextGetCid(ctx)
	return log.WithField(logCidKey, cid)
}

// LogWrapper is a handler wrapper that logs server request.
func LogWrapper(fn server.HandlerFunc) server.HandlerFunc {

	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		start := time.Now()
		err := fn(ctx, req, rsp)
		// RPC 计算经历的时间长度
		//no time.Since in order to format it well after
		end := time.Now()
		latency := end.Sub(start)
		cid := ContextGetCid(ctx)

		// l.Infof("%s %s", rpcMetadata, flatMetadata(md))

		l := log.
			WithField(logRpcCallKey, req.Method()).
			WithField(logCidKey, cid).
			WithField(logLatencyKey, latency.String())

		// Log rpc call execution result
		if err != nil {
			l.WithError(err).Warn(rpcFailed)
		} else {
			l.Info(rpcOk)
		}

		return err
	}
}

// flatMetadata 将 Metadata 打平为 "k=v" 形式的字符串序列
// func flatMetadata(md metadata.Metadata) string {
// 	var buffer bytes.Buffer
// 	for k, v := range md {
// 		buffer.WriteString(strconv.Quote(k))
// 		buffer.WriteString("=")
// 		buffer.WriteString(strconv.Quote(v))
// 		buffer.WriteString(" ")
// 	}

// 	return buffer.String()
// }
