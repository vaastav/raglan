package autotune

import "math/rand"

type Strategy interface {
	Name() string
	Init()
	NextConfig() *Configuration
}

type LinearStrategy struct {
	all_points map[string]*SpecPoint[any]
	max_vals   map[string]int
	cur_config *Configuration
}

func NewLinearStrategy(specpoints map[string]*SpecPoint[any]) *LinearStrategy {
	ls := &LinearStrategy{all_points: specpoints, max_vals: make(map[string]int)}
	return ls
}

func (ls *LinearStrategy) Name() string {
	return "LinearStrategy"
}

func (ls *LinearStrategy) Init() {
	for k, sp := range ls.all_points {
		ls.max_vals[k] = len(sp.Values)
	}
}

func (ls *LinearStrategy) NextConfig() *Configuration {
	conf := &Configuration{Mappings: make(map[string]int)}
	if ls.cur_config == nil {
		for k := range ls.all_points {
			// Pick the first index at initialization
			conf.Mappings[k] = 0
		}
		ls.cur_config = conf
	} else {
		for k := range ls.all_points {
			conf.Mappings[k] = (ls.cur_config.Mappings[k] + 1) % ls.max_vals[k]
		}
		ls.cur_config = conf
	}
	return conf
}

type RandomStrategy struct {
	all_points map[string]*SpecPoint[any]
	max_vals   map[string]int
	cur_config *Configuration
}

func NewRandomStrategy(specpoints map[string]*SpecPoint[any]) *RandomStrategy {
	rs := &RandomStrategy{all_points: specpoints, max_vals: make(map[string]int)}
	return rs
}

func (rs *RandomStrategy) Name() string {
	return "RandomStrategy"
}

func (rs *RandomStrategy) Init() {
	for k, sp := range rs.all_points {
		rs.max_vals[k] = len(sp.Values)
	}
}

func (rs *RandomStrategy) NextConfig() *Configuration {
	conf := &Configuration{Mappings: make(map[string]int)}
	for k := range rs.all_points {
		conf.Mappings[k] = rand.Intn(rs.max_vals[k])
	}
	rs.cur_config = conf
	return conf
}
