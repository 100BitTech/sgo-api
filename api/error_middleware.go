package api

import (
	"errors"
	"fmt"
	"net/http"
	"sgo-api/base"
	"strings"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware(log base.Logger, errorCodeMap map[int][]string) gin.HandlerFunc {
	log = log.WithTag("GIN")

	return func(c *gin.Context) {
		defer func() {
			if err := base.Recover(recover()); err != nil {
				log.Errorf("接口发生恐慌：%+v", err)
				rep(c, errorCodeMap, err)
			}
		}()

		c.Next()

		if len(c.Errors) <= 0 {
			return
		}

		err := c.Errors[len(c.Errors)-1]

		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("发生错误：%+v", err.Unwrap()))
		for i := len(c.Errors) - 2; i >= 0; i-- {
			sb.WriteString(fmt.Sprintf("\n%+v", c.Errors[i]))
		}
		log.Errorf(sb.String())

		rep(c, errorCodeMap, err)
	}
}

// 响应错误，返回 true 时说明有错误
func rep(c *gin.Context, errorCodeMap map[int][]string, err error) bool {
	if err == nil {
		return false
	}

	msg := err.Error()

	// 获取 code
	code := http.StatusInternalServerError

	var coder interface{ Code() int }
	if ok := errors.As(err, &coder); ok {
		code = coder.Code()
	} else {
		for c, es := range errorCodeMap {
			ok := false
			for _, e := range es {
				if strings.HasPrefix(msg, e) {
					code = c
					ok = true
					break
				}
			}
			if ok {
				break
			}
		}
	}

	c.JSON(code, gin.H{"code": code, "msg": msg})
	return true
}
