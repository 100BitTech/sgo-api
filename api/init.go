package api

import (
	"sgo-api/base"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Log base.Logger

	Port int

	TraceContextKey string
	GetTraceID      func(c *gin.Context) string
	ErrorCodeMap    map[int][]string
}

func Init(conf Config, extend func(*gin.Engine)) {
	if conf.Port <= 0 {
		conf.Port = 80
	}

	log := conf.Log.WithTag("GIN")

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// 中间件
	r.Use(gin.Recovery()) // ErrorMiddleware 已经处理了恐慌问题，这里作为最后一道保险
	r.Use(LogMiddleware(log, conf.TraceContextKey, conf.GetTraceID))
	r.Use(ErrorMiddleware(log, conf.ErrorCodeMap))

	if extend != nil {
		extend(r)
	}

	for {
		log.Infof("运行端口：%d", conf.Port)
		if err := r.Run(":" + strconv.Itoa(conf.Port)); err != nil {
			log.Errorf("运行出错：%+v", err)
		}
	}
}
