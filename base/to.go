package base

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// 等同于 v.(T)，但是处理了 nil 的情况
func To[T any](v any) T {
	if v == nil {
		var vv T
		return vv
	} else {
		return v.(T)
	}
}

func ToTime(t any, loc *time.Location, layouts ...string) (time.Time, bool) {
	if t == nil {
		return time.Time{}, true
	}
	if loc == nil {
		loc = time.Local
	}

	if tt, ok := t.(time.Time); ok {
		return tt, true
	}
	if tt, ok := t.(*time.Time); ok {
		return *tt, true
	}

	s := fmt.Sprintf("%v", t)

	for _, layout := range NewSet[string](append(layouts,
		time.RFC3339,
		DateTimeLayout,
		DateLayout,
	)...).ToSlice() {
		if t, err := time.ParseInLocation(layout, s, loc); err == nil {
			return t, true
		}
	}

	// 兼容 time.Time 的 String 方法
	if strings.Contains(s, " m=") {
		ss := strings.Split(s, " m=")

		if t, err := time.ParseInLocation("2006-01-02 15:04:05.999999999 -0700 MST", ss[0], loc); err == nil {
			return t, true
		}
	}

	// 兼容“2023-10-09 12”之类的情况
	if f := "0001-01-01 00:00:00"; len(s) < len(f) {
		_s := s + f[len(s):]

		if t, err := time.ParseInLocation(DateTimeLayout, _s, loc); err == nil {
			return t, true
		}
	}

	return time.Time{}, false
}

func ToNumber(n any) (float64, bool) {
	if n == nil {
		return 0, true
	}

	switch v := n.(type) {
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case *int:
		if v == nil {
			return 0, true
		} else {
			return float64(*v), true
		}
	case *int8:
		if v == nil {
			return 0, true
		} else {
			return float64(*v), true
		}
	case *int16:
		if v == nil {
			return 0, true
		} else {
			return float64(*v), true
		}
	case *int32:
		if v == nil {
			return 0, true
		} else {
			return float64(*v), true
		}
	case *int64:
		if v == nil {
			return 0, true
		} else {
			return float64(*v), true
		}
	case *uint:
		if v == nil {
			return 0, true
		} else {
			return float64(*v), true
		}
	case *uint8:
		if v == nil {
			return 0, true
		} else {
			return float64(*v), true
		}
	case *uint16:
		if v == nil {
			return 0, true
		} else {
			return float64(*v), true
		}
	case *uint32:
		if v == nil {
			return 0, true
		} else {
			return float64(*v), true
		}
	case *uint64:
		if v == nil {
			return 0, true
		} else {
			return float64(*v), true
		}
	case *float32:
		if v == nil {
			return 0, true
		} else {
			return float64(*v), true
		}
	case *float64:
		if v == nil {
			return 0, true
		} else {
			return *v, true
		}
	default:
		s := fmt.Sprintf("%v", n)
		if nn, err := strconv.ParseFloat(s, 64); err != nil {
			return 0, false
		} else {
			return nn, true
		}
	}
}
