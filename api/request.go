package api

import (
	"bytes"
	"io"
	"sgo-api/base"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/oops"
)

const (
	requestContextKey = "sgo-api.request"
)

type Request struct {
	IP        string
	Method    string
	Path      string
	StartTime time.Time

	Error error

	body []byte
}

func newRequest(c *gin.Context) (req *Request) {
	req = &Request{}

	defer func() {
		if err := base.Recover(recover()); err != nil {
			req.Error = oops.Wrapf(err, "读取请求发生恐慌")
		}
	}()

	req.IP = c.ClientIP()
	req.StartTime = time.Now()
	req.Method = c.Request.Method

	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	if raw != "" {
		path = path + "?" + raw
	}
	req.Path = path

	if body, err := io.ReadAll(c.Request.Body); err != nil {
		req.Error = oops.Wrapf(err, "读取请求 Body 出错")
	} else {
		if body == nil {
			body = []byte{}
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		req.body = body
	}

	return
}

func (req *Request) Body() ([]byte, error) {
	if req.Error != nil {
		return nil, oops.Wrap(req.Error)
	}

	return req.body, nil
}

func GetRequest(c *gin.Context) *Request {
	if v, ok := c.Get(requestContextKey); ok {
		if req, ok := v.(*Request); ok {
			return req
		}
	}

	req := newRequest(c)
	c.Set(requestContextKey, req)
	return req
}
