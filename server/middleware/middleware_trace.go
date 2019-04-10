package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/j75689/easybot/pkg/logger"
)

// TraceMiddleware trace entry point
func TraceMiddleware(skipper RouteSkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if skipper != nil {
			if skipper(c) {
				c.Next()
				return
			}
		}
		var formdata bytes.Buffer
		io.Copy(&formdata, c.Request.Body)
		c.Request.Body = ioutil.NopCloser(&formdata)
		logger.Infow("[trace]",
			"ip", c.ClientIP(),
			"proto", c.Request.Proto,
			"method", c.Request.Method,
			"url", c.Request.URL.String(),
			"header", c.Request.Header,
			"user_agent", c.GetHeader("User-Agent"),
			"time", time.Now().Format(time.RFC3339),
			"formdata", formdata.String())
	}
}
