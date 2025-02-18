package base

import (
	"encoding/json"

	"github.com/samber/oops"
)

// 基于 encoding/json, 适用于简单结构体的深度复制
func Copy[T any](src T) (T, error) {
	var dst T
	if bs, err := json.Marshal(src); err != nil {
		return dst, oops.Wrap(err)
	} else if err := json.Unmarshal(bs, &dst); err != nil {
		return dst, oops.Wrap(err)
	}
	return dst, nil
}
