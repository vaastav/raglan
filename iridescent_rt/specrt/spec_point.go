package specrt

import "fmt"

type SpecPointEnum int

const (
	CompileTime SpecPointEnum = iota
	Knob
)

type SpecializationPoint interface {
	Specialize(index int)
	String() string
	Type() SpecPointEnum
	NumVals() int
}

type CompileTimeSpecPoint[T any] struct {
	Name          string
	ParentFn      string
	Values        []T
	Current       T
	IridType      string
	GoType        string
	IsSpecialized bool
}

func NewCompileTimeSpecPoint[T any](name string, values []T) *CompileTimeSpecPoint[T] {
	sp := &CompileTimeSpecPoint[T]{Name: name, Values: values}
	return sp
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
