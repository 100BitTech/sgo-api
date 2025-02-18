package base

import (
	"github.com/dranikpg/dto-mapper"
	"github.com/samber/oops"
)

var (
	defultMapper = NewMapper()
)

// 快速映射对象
func Map(src any, dst any) error {
	return defultMapper.Map(src, dst)
}

type Mapper interface {
	Map(src any, dst any) error

	// 类型转换器
	//   A -> B：AddConvFunc(func(s A) B { ... })
	//   A -> B：AddConvFunc(func(s A) (B, error) { ... })
	AddTypeConverter(f any) Mapper

	// 映射成功后的触发器，目标必须是指针
	//   A -> B：AddAfterTrigger(func(t *B) { ... })
	//   A -> B：AddAfterTrigger(func(t *B, s A) { ... })
	//   A -> B：AddAfterTrigger(func(t *B) error { ... })
	//   A -> B：AddAfterTrigger(func(t *B, s A) error { ... })
	AddAfterTrigger(f any) Mapper
}

func NewMapper() Mapper {
	m := NewDtoMapper()
	Expand(m)
	return m
}

func Expand(m Mapper) {
	// 需要增加一些全局的扩展，可以在这里加上
}

type DtoMapper struct {
	m *dto.Mapper
}

func NewDtoMapper() Mapper {
	return &DtoMapper{
		m: &dto.Mapper{},
	}
}

func (m *DtoMapper) Map(src any, dst any) error {
	if err := m.m.Map(dst, src); err != nil {
		return oops.Wrap(err)
	} else {
		return nil
	}
}

func (m *DtoMapper) AddTypeConverter(f any) Mapper {
	m.m.AddConvFunc(f)
	return m
}

func (m *DtoMapper) AddAfterTrigger(f any) Mapper {
	m.m.AddInspectFunc(f)
	return m
}
