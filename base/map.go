package base

import (
	"fmt"
	"sync"
	"time"

	"github.com/samber/oops"
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

// 并发安全的，可以定时更新元素的映射
type UpdateMap[C UpdateMapItemConfig, T any] struct {
	m *SyncMap[string, T]

	config UpdateMapConfig[C, T]

	lock sync.Mutex
}

func NewUpdateMap[C UpdateMapItemConfig, T any](conf UpdateMapConfig[C, T]) (*UpdateMap[C, T], error) {
	if conf.Name == "" {
		conf.Name = NewNanoID()
	}

	m := UpdateMap[C, T]{m: NewSyncMap[string, T](), config: conf}

	if err := m.update(true); err != nil {
		return nil, oops.Wrap(err)
	}

	return &m, nil
}

func (m *UpdateMap[C, T]) Get(key string) (T, error) {
	value, ok := m.m.Load(key)
	if !ok {
		return value, fmt.Errorf("[my.UpdateMap] 映射 %v 无法获取 %v 元素", m.config.Name, key)
	}

	return value, nil
}

func (m *UpdateMap[C, T]) GetKeys() []string {
	keys := []string{}

	m.m.Range(func(key string, value T) bool {
		keys = append(keys, key)
		return true
	})

	return keys
}

func (m *UpdateMap[C, T]) Contains(key string) bool {
	_, err := m.Get(key)
	return err == nil
}

func (m *UpdateMap[C, T]) Range(f func(key string, value T) bool) {
	m.m.Range(func(key string, value T) bool {
		return f(key, value)
	})
}

func (m *UpdateMap[C, T]) Update() error {
	return m.update(false)
}

func (m *UpdateMap[C, T]) update(isAfter bool) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.config.GetItemConfigs == nil || m.config.NewItem == nil {
		return oops.Errorf("[my.UpdateMap] 映射 %v 无法初始化，请检查配置", m.config.Name)
	}

	confs, err := m.config.GetItemConfigs()
	if err != nil {
		return oops.Wrap(err)
	}

	keySet := NewSet[string]()
	for _, conf := range confs {
		keySet.Add(conf.Key())

		if _, ok := m.m.Load(conf.Key()); ok {
			m.config.Log.Warnf("[my.UpdateMap] 映射 %v 元素 %v 加载成功", m.config.Name, conf.Key())
			continue
		}

		item, err := m.config.NewItem(conf)
		if err != nil {
			m.config.Log.Errorf("[my.UpdateMap] 映射 %v 元素 %v 加载失败：%+v", m.config.Name, conf.Key(), err)
			continue
		}

		m.m.Store(conf.Key(), item)
		m.config.Log.Warnf("[my.UpdateMap] 映射 %v 元素 %v 加载成功", m.config.Name, conf.Key())
	}

	// 移除多余的元素
	m.m.Range(func(key string, item T) bool {
		if !keySet.Contains(key) {
			m.m.Delete(key)
			m.config.Log.Warnf("[my.UpdateMap] 映射 %v 元素 %v 移除成功", m.config.Name, key)
		}

		return true
	})

	if isAfter {
		// 定时更新
		if m.config.UpdateInterval > 0 {
			time.AfterFunc(time.Minute*time.Duration(m.config.UpdateInterval), func() {
				Try(func() {
					m.config.Log.Warnf("[my.UpdateMap] 映射 %v 每 %d 分钟刷新中...", m.config.Name, m.config.UpdateInterval)

					if err := m.update(true); err != nil {
						m.config.Log.Errorf("[my.UpdateMap] 映射 %v 刷新失败：%+v", m.config.Name, err)
					}
				}).Catch(func(err error) {
					m.config.Log.Errorf("[my.UpdateMap] 映射 %v 刷新恐慌：%+v", m.config.Name, err)
				}).Do()
			})
		}
	}

	return nil
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
