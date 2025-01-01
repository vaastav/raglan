package autotune

type SpecPoint[T comparable] struct {
	Name     string
	Values   []T
	Default  T
	Current  T
	SetValFn func(val T)
}

func NewSpecPoint[T comparable](name string, values []T, default_val T, SetValFn func(T)) *SpecPoint[T] {
	sp := &SpecPoint[T]{Name: name, Values: values, Default: default_val, Current: default_val, SetValFn: SetValFn}
	sp.SetValFn(sp.Current)
	return sp
}

func (sp *SpecPoint[T]) Specialize(index int) {
	sp.Current = sp.Values[index]
	sp.SetValFn(sp.Current)
}
