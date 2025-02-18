package base

import (
	"sync"
)

// 支持泛型的并发映射
type SyncMap[K comparable, V any] struct {
	m sync.Map
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		m: sync.Map{},
	}
}

func (m *SyncMap[K, V]) Load(key K) (V, bool) {
	value, ok := m.m.Load(key)
	return To[V](value), ok
}

func (m *SyncMap[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}

func (m *SyncMap[K, V]) Delete(key K) {
	m.m.Delete(key)
}

func (m *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(key, value any) bool {
		return f(key.(K), To[V](value))
	})
}

func (m *SyncMap[K, V]) LoadOrStore(key K, value V) (V, bool) {
	actual, loaded := m.m.LoadOrStore(key, value)
	return To[V](actual), loaded
}

func (m *SyncMap[K, V]) LoadAndDelete(key K) (V, bool) {
	actual, loaded := m.m.LoadAndDelete(key)
	return To[V](actual), loaded
}

type UpdateMapItemConfig interface {
	Key() string
}

type UpdateMapConfig[C UpdateMapItemConfig, T any] struct {
	Name           string
	UpdateInterval int64

	Log Logger

	GetItemConfigs func() ([]C, error)
	NewItem        func(C) (T, error)
}

func GetMapValue(data map[string]any, keys ...string) any {
	if len(data) == 0 {
		return nil
	}

	var v any
	v = data
	for _, k := range keys {
		vv, ok := v.(map[string]any)
		if !ok {
			return nil
		}

		v, ok = vv[k]
		if !ok || v == nil {
			return nil
		}
	}
	return v
}
