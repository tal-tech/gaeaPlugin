package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tal-tech/loggerX"
	"github.com/tal-tech/loggerX/logtrace"
	"github.com/tal-tech/xtools/perfutil"
	"github.com/tal-tech/xtools/traceutil"
)

var (
	hostname, _ = os.Hostname()
)

// LoggerMiddleware 创建logger相关
func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("logid", strconv.FormatInt(logger.Id(), 10))
		ctx.Set("hostname", hostname)
		ctx.Set("start", time.Now())
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery
		var body []byte
		if ctx.Request.Body != nil {
			body, _ = ioutil.ReadAll(ctx.Request.Body)
		}
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		if raw != "" {
			path = path + "?" + raw
		}

		logtraceMap := logtrace.GenLogTraceMetadata()
		logtraceMap.Set("request_uri", fmt.Sprintf("\"%s\"", path))
		if len(body) > 0 && bytes.HasPrefix(body, []byte("{")) {
			logtraceMap.Set("x_request_param", fmt.Sprintf("%s", body))
		} else {
			logtraceMap.Set("request_param", fmt.Sprintf("\"%s\"", body))
		}
		logtraceMap.Set("request_method", fmt.Sprintf("\"%s\"", ctx.Request.Method))
		logtraceMap.Set("request_client_ip", fmt.Sprintf("\"%s\"", ctx.ClientIP()))
		if traceId := ctx.GetHeader("traceid"); traceId != "" {
			logtraceMap.Set("x_trace_id", "\""+traceId+"\"")
			if strings.HasPrefix(traceId, "pts_") {
				ctx.Set("IS_BENCHMARK", "1")
			}
		}
		if traceId := ctx.GetHeader("traceId"); traceId != "" {
			logtraceMap.Set("x_trace_id", "\""+traceId+"\"")
		}
		if rpcId := ctx.GetHeader("rpcid"); rpcId != "" {
			rpcId = rpcId + ".0"
			logtraceMap.Set("x_rpcid", "\""+rpcId+"\"")
		}
		if rpcId := ctx.GetHeader("rpcId"); rpcId != "" {
			rpcId = rpcId + ".0"
			logtraceMap.Set("x_rpcid", "\""+rpcId+"\"")
		}
		ctx.Set(logtrace.GetMetadataKey(), logtraceMap)
	}
}

// PerfMiddleware Perf监控
func PerfMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		perfutil.CountI(ctx.Request.URL.Path)
		defer perfutil.AutoElapsed(ctx.Request.URL.Path, time.Now())
		ctx.Next()
	}
}

// TraceMiddleware Zipkin链路监控
func TraceMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		span, _ := traceutil.Trace(context.Background(), ctx.Request.URL.Path)
		if span != nil {
			defer span.Finish()
			ctx.Set(traceutil.SpanKey, span)
		}
		ctx.Next()
	}
}
