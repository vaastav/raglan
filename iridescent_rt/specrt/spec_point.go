package specrt

import "fmt"

type SpecPoint[T comparable] struct {
	Name          string
	ParentFn      string
	Values        []T
	Current       T
	IridType      string
	GoType        string
	IsSpecialized bool
}

func NewSpecPoint[T comparable](name string, values []T) *SpecPoint[T] {
	sp := &SpecPoint[T]{Name: name, Values: values}
	return sp
}

func (sp *SpecPoint[T]) String() string {
	s := sp.ParentFn + "." + sp.Name + ", values: " + fmt.Sprintf("%v", sp.Values)
	return s
}

func (sp *SpecPoint[T]) Specialize(index int) {
	sp.IsSpecialized = true
	sp.Current = sp.Values[index]
}
