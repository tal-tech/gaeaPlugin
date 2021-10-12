package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http/httputil"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	logger "github.com/tal-tech/loggerX"

	"github.com/tal-tech/gaeaPlugin/xesgin"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		ctx := xesgin.TransferToContext(c)
		logger.Dx(ctx, "[GIN]", "%v | %3d | %13v | %15s | %-7s %s\n%s",
			end.Format("2006/01/02 - 15:04:05"), statusCode, latency, clientIP, method, path, comment,
		)
	}
}

//创建logid
func HandleCtxInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("logid", strconv.FormatInt(logger.Id(), 10))
		ctx.Set("start", time.Now())
	}
}

//revover middleware
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := stack(3)
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				logger.E("[Recovery]", " %s panic recovered:\n%s\n%s\n%s", time.Now().Format("2006/01/02 - 15:04:05"), string(httprequest), err, stack)
				fmt.Println("[Recovery]", " %s panic recovered:\n%s\n%s\n%s", time.Now().Format("2006/01/02 - 15:04:05"), string(httprequest), err, string(stack))
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//  runtime/debug.*T·ptrmethod
	// and want
	//  *T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}
