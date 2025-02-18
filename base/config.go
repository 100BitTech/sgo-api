package base

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/samber/oops"
)

type BaseConfig struct {
	Log LogConfig `json:"log"`
}

// 初始化缓存
//   - 默认配置目录为 ./
//   - 可以通过入参指定目录（优先级最高）
//   - 可以编译时加上 -ldflags '-X github.com/100BitTech/sgo-api/base.CONFIGPATH=./xxx' 指定目录
//   - 如果设置了环境，会优先加载指定环境的配置。比如配置目录为 ./xxx，环境为 prod，那么优先加载 ./xxx/config.prod.json，不存在再加载 ./xxx/config.json
func InitConfig[T any](name string) T {
	if name == "" {
		name = CONFIGPATH
	}
	if name == "" {
		name = "./"
	}

	if ENV != "" {
		if conf, ok := readConfig[T](path.Join(name, fmt.Sprintf("config.%v.json", ENV))); ok {
			return conf
		}
	}

	if conf, ok := readConfig[T](path.Join(name, "config.json")); ok {
		return conf
	}

	panic("无法加载配置")
}

func readConfig[T any](name string) (T, bool) {
	if _, err := os.Stat(name); err != nil {
		var conf T
		return conf, false
	}
	fmt.Printf("配置路径：%v\n", name)

	data, err := os.ReadFile(name)
	if err != nil {
		panic(oops.Wrapf(err, "读取 %v 配置失败", name))
	}

	var conf T
	if err = json.Unmarshal(data, &conf); err != nil {
		panic(oops.Wrapf(err, "解析 %v 配置失败", name))
	}
	return conf, true
}
