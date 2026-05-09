package specrt

import (
	"fmt"
	"sync"
)

type SpecPointEnum int

const (
	CompileTime SpecPointEnum = iota
	CodeGenPointArr
	Knob
)

type SpecializationPoint interface {
	Specialize(index int)
	String() string
	Type() SpecPointEnum
	NumVals() int
	GetName() string
	Incr(key int)
}

type CompileTimeSpecPoint[T comparable] struct {
	Name          string
	ParentFn      string
	Values        []T
	Current       T
	IridType      string
	GoType        string
	IsSpecialized bool
	Counter       []uint64
	sync.Mutex
}

func NewCompileTimeSpecPoint[T comparable](name string, values []T) *CompileTimeSpecPoint[T] {
	sp := &CompileTimeSpecPoint[T]{Name: name, Values: values, Counter: make([]uint64, len(values))}
	for idx, _ := range sp.Counter {
		sp.Counter[idx] = 0
	}
	return sp
}

func (sp *CompileTimeSpecPoint[T]) Incr(key int) {
	sp.Counter[key] = sp.Counter[key] + 1
}

func (sp *CompileTimeSpecPoint[T]) ResetStats() {
	sp.Lock()
	for k, _ := range sp.Counter {
		sp.Counter[k] = 0
	}
	sp.Unlock()
}

func (sp *CompileTimeSpecPoint[T]) String() string {
	s := sp.ParentFn + "." + sp.Name + ", values: " + fmt.Sprintf("%v", sp.Values)
	return s
}

func (sp *CompileTimeSpecPoint[T]) Specialize(index int) {
	sp.IsSpecialized = true
	sp.Current = sp.Values[index]
}

func (sp *CompileTimeSpecPoint[T]) NumVals() int {
	return len(sp.Values)
}

func (sp *CompileTimeSpecPoint[T]) Type() SpecPointEnum {
	return CompileTime
}

func (sp *CompileTimeSpecPoint[T]) GetName() string {
	return sp.Name
}

type CodeGenPointArrSpecPoint[T comparable] struct {
	Name     string
	Values   [][]T
	Default  []T
	Current  []T
	SetValFn func(val []T)
}

func NewCodeGenPointArrSpecPoint[T comparable](name string, values [][]T, default_val []T, SetValFn func([]T)) *CodeGenPointArrSpecPoint[T] {
	s := &CodeGenPointArrSpecPoint[T]{Name: name, Values: values, Default: default_val, Current: default_val, SetValFn: SetValFn}
	s.SetValFn(s.Current)
	return s
}

func (sp *CodeGenPointArrSpecPoint[T]) Specialize(index int) {
	sp.Current = sp.Values[index]
	sp.SetValFn(sp.Current)
}

func (sp *CodeGenPointArrSpecPoint[T]) Type() SpecPointEnum {
	return CodeGenPointArr
}

func (sp *CodeGenPointArrSpecPoint[T]) NumVals() int {
	return len(sp.Values)
}

func (sp *CodeGenPointArrSpecPoint[T]) String() string {
	s := sp.Name + ", values: " + fmt.Sprintf("%v", sp.Values) + ", default: " + fmt.Sprintf("%v", sp.Default) + ", current: " + fmt.Sprintf("%v", sp.Current)
	return s
}

func (sp *CodeGenPointArrSpecPoint[T]) Incr(key int) { /* Not implemented*/ }

func (sp *CodeGenPointArrSpecPoint[T]) GetName() string {
	return sp.Name
}

type KnobSpecPoint[T comparable] struct {
	Name     string
	Values   []T
	Default  T
	Current  T
	SetValFn func(val T)
}

func NewKnobSpecPoint[T comparable](name string, values []T, default_val T, SetValFn func(T)) *KnobSpecPoint[T] {
	s := &KnobSpecPoint[T]{Name: name, Values: values, Default: default_val, Current: default_val, SetValFn: SetValFn}
	s.SetValFn(s.Current)
	return s
}

func (sp *KnobSpecPoint[T]) Specialize(index int) {
	sp.Current = sp.Values[index]
	sp.SetValFn(sp.Current)
}

func (sp *KnobSpecPoint[T]) String() string {
	s := sp.Name + ", values: " + fmt.Sprintf("%v", sp.Values) + ", default: " + fmt.Sprintf("%v", sp.Default) + ", current: " + fmt.Sprintf("%v", sp.Current)
	return s
}

func (sp *KnobSpecPoint[T]) Type() SpecPointEnum {
	return Knob
}

func (sp *KnobSpecPoint[T]) NumVals() int {
	return len(sp.Values)
}

func (sp *KnobSpecPoint[T]) GetName() string {
	return sp.Name
}

func (sp *KnobSpecPoint[T]) Incr(key int) { /* Not implemented */ }
