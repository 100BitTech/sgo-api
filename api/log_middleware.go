package api

import (
	"context"
	"sgo-api/base"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	dateTimeLayout = "2006-01-02 15:04:05.000"
)

func LogMiddleware(log base.Logger, traceContextKey string, getTraceID func(c *gin.Context) string) gin.HandlerFunc {
	log = log.WithTag("GIN")

	if traceContextKey == "" {
		traceContextKey = base.TraceContextKey
	}

	return func(c *gin.Context) {
		traceID := func() string {
			defer func() {
				if err := base.Recover(recover()); err != nil {
					log.Warnf("日志中间获取 Trace ID 恐慌：%+v", err)
				}
			}()

			value := ""
			if getTraceID != nil {
				value = getTraceID(c)
			}
			if value == "" {
				value = base.NewNanoID()
			}
			return value
		}()

		ctx := context.WithValue(c.Request.Context(), traceContextKey, traceID)
		c.Request = c.Request.WithContext(ctx)

		log := log.WithTrace(ctx, traceContextKey)

		req := GetRequest(c)
		if req.Error != nil {
			log.Errorf("读取请求失败：%+v", req.Error)
		}

		body, _ := req.Body()

		log.Infof("REQ:%v %v | %v %v\n%v",
			req.IP, req.StartTime.Format(dateTimeLayout),
			req.Method, req.Path,
			string(body),
		)

		// 重载 Writer，以便后续获取 Response 的 Body
		w := NewWriter(c.Writer)
		c.Writer = w

		c.Next()

		endTime := time.Now()
		latency := endTime.Sub(req.StartTime)
		statusCode := c.Writer.Status()

		log.Infof("RESP:%v %v | %v %v\n%d %v\n%v",
			req.IP, endTime.Format(dateTimeLayout),
			req.Method, req.Path,
			statusCode, latency,
			w.String(),
		)
	}
}
