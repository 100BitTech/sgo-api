package base

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/samber/oops"
)

var (
	cache struct {
		l Logger
		c *bigcache.BigCache
	}
)

// 初始化缓存。参数 maxCacheSize 单位为 MB
func InitCache(log Logger, maxCacheSize int) {
	c, err := bigcache.New(context.Background(), bigcache.Config{
		Shards:           8,
		LifeWindow:       3 * time.Minute,
		CleanWindow:      1 * time.Second,
		HardMaxCacheSize: maxCacheSize,
	})
	if err != nil {
		panic(err)
	}

	cache = struct {
		l Logger
		c *bigcache.BigCache
	}{
		l: log,
		c: c,
	}
}

func Cache[T any](key string, newValue func() (T, error)) (T, error) {
	var zeroValue T

	cv, err := cache.c.Get(key)
	if err == nil {
		var value T
		if err := json.Unmarshal(cv, &value); err != nil {
			return zeroValue, oops.Wrap(err)
		} else {
			return value, nil
		}
	} else if !errors.Is(err, bigcache.ErrEntryNotFound) {
		return zeroValue, oops.Wrap(err)
	}

	value, err := newValue()
	if err != nil {
		return zeroValue, oops.Wrap(err)
	}

	jv, err := json.Marshal(value)
	if err != nil {
		return zeroValue, oops.Wrap(err)
	}
	if err := cache.c.Set(key, jv); err != nil {
		cache.l.Errorf("设置 %v 缓存错误：%+v", key, oops.Wrap(err))
	}

	return value, nil
}
