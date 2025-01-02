package autotune

import (
	"time"

	"github.com/vaastav/iridescent/iridescent_rt/specrt"
)

// Configuration contains the mappings for the specialized value for each point.
// The value isn't directly stored but rather an index into the point's values array is stored.
type Configuration struct {
	Id       uint64
	Mappings map[string]int
}

// Stores collected statistics in a map.
// Key is the name of the statistic (eg: tput) and Value is the measured value (eg: 1000)
type Stats struct {
	Values map[string]uint64
}

// Function types for convenience
type ObjectiveFn func(s Stats) uint64
type MeasurementFn func() Stats

type ExplorationEngine struct {
	CurConfig    uint64
	Configs      map[uint64]*Configuration
	SpecPoints   map[string]specrt.SpecializationPoint
	ConfigScores map[uint64]uint64
	Measure      MeasurementFn
	Objective    ObjectiveFn
	Dur          time.Duration
	Period       time.Duration
	Strat        Strategy
	SpecRT       *specrt.SpecializationRuntime
}

func NewExplorationEngine(specpoints map[string]specrt.SpecializationPoint, duration time.Duration, period time.Duration, strategy Strategy, srt *specrt.SpecializationRuntime) *ExplorationEngine {
	e := &ExplorationEngine{}
	e.SpecPoints = specpoints
	e.Configs = make(map[uint64]*Configuration)
	e.ConfigScores = make(map[uint64]uint64)
	e.CurConfig = 0
	e.Dur = duration
	e.Period = period
	e.Strat = strategy
	e.SpecRT = srt
	return e
}

func (e *ExplorationEngine) StartExploration() {
	e.Strat.Init()
	stop_chan := make(chan bool)
	finished := make(chan bool)
	go func() {
		for {
			select {
			case <-stop_chan:
				// Exploration is done!
				finished <- true
				return
			default:
				e.CurConfig += 1
				c := e.Strat.NextConfig()
				c.Id = e.CurConfig
				// Set the configuration
				e.SelectConfig(c)
				time.Sleep(e.Period)
				// Measure the impact of the selected configuration
				stats := e.Measure()
				obj := e.Objective(stats)
				e.ConfigScores[c.Id] = obj
			}
		}

	}()
	time.Sleep(e.Dur)
	// Send stop signal for exploration
	stop_chan <- true
	// Wait till the exploration finishes so that we don't interfere with an ongoing config
	<-finished
	e.Finalize()
}

func (e *ExplorationEngine) SelectConfig(c *Configuration) error {
	// Specialize the points to the best config
	for name, idx := range c.Mappings {
		sp := e.SpecPoints[name]
		sp.Specialize(idx)
	}
	if e.SpecRT != nil {
		err := e.SpecRT.UpdatePlugin()
		if err != nil {
			return err
		}
	}
	return nil
}

// Finalize specializes the specialization points once the exploration ends.
func (e *ExplorationEngine) Finalize() error {
	// Select the best configuration
	highest_idx := uint64(0)
	highest_val := uint64(0)
	for k, v := range e.ConfigScores {
		if v > highest_val {
			highest_val = v
			highest_idx = k
		}
	}
	chosen_config := e.Configs[highest_idx]
	return e.SelectConfig(chosen_config)
}

func (e *ExplorationEngine) ResetExploration() {
	clear(e.ConfigScores)
	clear(e.Configs)
	e.CurConfig = 0
}
