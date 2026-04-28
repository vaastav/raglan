package autotune

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vaastav/raglan/iridescent_rt/specrt"
)

type IridescentRTIface interface {
	StartExploration(ctx context.Context) error
	ResetExploration(ctx context.Context) error
}

type IridescentRT struct {
	AllPoints map[string]specrt.SpecializationPoint
	ExpEngine *ExplorationEngine
	SpecRT    *specrt.SpecializationRuntime
}

var rt *IridescentRT

func NewIridescentRT(ctx context.Context, duration string, period string, strategy_name string, specialization_file string) (*IridescentRT, error) {
	// Ensure only once initialization
	if rt == nil {
		new_rt := &IridescentRT{}
		new_rt.AllPoints = make(map[string]specrt.SpecializationPoint)
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
			strat = NewLinearStrategy(new_rt.AllPoints)
		} else if strategy_name == "random" {
			strat = NewRandomStrategy(new_rt.AllPoints)
		} else {
			return nil, errors.New(fmt.Sprintf("Unknown strategy chosen: %s", strategy_name))
		}
		if specialization_file != "" {
			srt, err := specrt.NewSpecializationRuntime(ctx, specialization_file)
			if err != nil {
				return nil, err
			}
			new_rt.SpecRT = srt
		}
		// Specialization runtime will be nil if no target spec file is provided!
		engine := NewExplorationEngine(new_rt.AllPoints, dur, per, strat, new_rt.SpecRT)
		new_rt.ExpEngine = engine
		rt = new_rt
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

func (irid *IridescentRT) StartExploration(ctx context.Context) error {
	go rt.ExpEngine.StartExploration()
	return nil
}

func (irid *IridescentRT) ResetExploration(ctx context.Context) error {
	rt.ExpEngine.ResetExploration()
	return nil
}

func (irid *IridescentRT) RegisterKnob(name string, sp *specrt.KnobSpecPoint[any]) {
	rt.AllPoints[name] = sp
}

func (irid *IridescentRT) RegisterCompileTimeSpecPoint(name string, sp *specrt.CompileTimeSpecPoint[any]) {
	rt.AllPoints[name] = sp
}
