package base

import (
	"fmt"
	"strings"
)

const (
	EnvTest = "test" // 测试
	EnvPre  = "pre"  // 预发
	EnvProd = "prod" // 生产

	TraceContextKey = "sgo.trace_id"
)

var (
	// 环境（小写），-ldflags '-X github.com/100BitTech/sgo-api/base.ENV=dev'
	ENV string = ""

	// 配置路径，-ldflags '-X github.com/100BitTech/sgo-api/base.CONFIGPATH=./config.json'
	CONFIGPATH string = ""

	// 日志路径，-ldflags '-X github.com/100BitTech/sgo-api/base.LOGPATH=./logs'
	LOGPATH string = ""
)

func init() {
	ENV = strings.ToLower(strings.TrimSpace(ENV))

	fmt.Printf("环境：%v\n", ENV)
}

func IsTestENV() bool {
	return ENV == EnvTest
}

func IsPreENV() bool {
	return ENV == EnvPre
}

func IsProdENV() bool {
	return ENV == EnvProd
}
