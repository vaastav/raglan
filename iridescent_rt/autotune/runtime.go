package autotune

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type IridescentRT struct {
	AllPoints map[string]*SpecPoint[any]
	ExpEngine *ExplorationEngine
}

var rt *IridescentRT

func NewIridescentRT(ctx context.Context, duration string, period string, strategy_name string) (*IridescentRT, error) {
	// Ensure only once initialization
	if rt == nil {
		rt = &IridescentRT{}
		rt.AllPoints = make(map[string]*SpecPoint[any])
		dur, err := time.ParseDuration(duration)
		if err != nil {
			return nil, err
		}
		per, err := time.ParseDuration(period)
		if err != nil {
			return nil, err
		}
		var strat Strategy
		if strategy_name == "linear" {
			strat = NewLinearStrategy(rt.AllPoints)
		} else if strategy_name == "random" {
			strat = NewRandomStrategy(rt.AllPoints)
		} else {
			return nil, errors.New(fmt.Sprintf("Unknown strategy chosen: %s", strategy_name))
		}
		engine := NewExplorationEngine(rt.AllPoints, dur, per, strat)
		rt.ExpEngine = engine
	}
	return rt, nil
}

func GetRuntime() *IridescentRT {
	return rt
}

func (irid *IridescentRT) RegisterMeasurementFn(mes MeasurementFn) {
	rt.ExpEngine.Measure = mes
}

func (irid *IridescentRT) RegisterObjFn(obj ObjectiveFn) {
	rt.ExpEngine.Objective = obj
}

func (irid *IridescentRT) StartExploration() {
	go rt.ExpEngine.StartExploration()
}

func (irid *IridescentRT) ResetExploration() {
	rt.ExpEngine.ResetExploration()
}

func (irid *IridescentRT) RegisterSpecPoint(name string, sp *SpecPoint[any]) {
	rt.AllPoints[name] = sp
}
