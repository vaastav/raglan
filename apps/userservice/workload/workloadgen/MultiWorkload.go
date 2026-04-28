package workloadgen

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/blueprint-uservices/blueprint/runtime/core/workload"
	"github.com/vaastav/iridescent/apps/userservice/workflow"
)

// Workload specific flags
var outfilew1 = flag.String("outfilew1", "statsw1.csv", "Outfile where individual request information will be stored for W1")
var outfilew2 = flag.String("outfilew2", "statsw2.csv", "Outfile where individual request information will be stored for W2")
var duration = flag.String("duration", "1m", "Duration for which the workload should be run")
var tput = flag.Int64("tput", 1000, "Desired throughput")
var mode = flag.String("mode", "irid", "One of irid|hard1|hard2")
var add_users = flag.Bool("addusers", false, "Add new users to the databases")

type MultiWorkload interface {
	ImplementsMultiWorkload(ctx context.Context) error
}

type multiWldGen struct {
	MultiWorkload

	userservice workflow.UserService
}

func NewMultiWorkload(ctx context.Context, userservice workflow.UserService) (MultiWorkload, error) {
	w := &multiWldGen{userservice: userservice}
	return w, nil
}

type FnType func() error

func statWrapper(fn FnType) workload.Stat {
	start := time.Now()
	err := fn()
	duration := time.Since(start)
	s := workload.Stat{}
	s.Start = start.UnixNano()
	s.Duration = duration.Nanoseconds()
	s.IsError = (err != nil)
	return s
}

func (w *multiWldGen) RunNAHandler(ctx context.Context) workload.Stat {
	user := &workflow.User{Username: fmt.Sprintf("NA_user_%d", rand.IntN(1000)+1)}
	return w.RunUserHandler(ctx, user)
}

func (w *multiWldGen) RunEUHandler(ctx context.Context) workload.Stat {
	user := &workflow.User{Username: fmt.Sprintf("EU_user_%d", rand.IntN(1000)+1)}
	return w.RunUserHandler(ctx, user)
}

func (w *multiWldGen) RunASHandler(ctx context.Context) workload.Stat {
	user := &workflow.User{Username: fmt.Sprintf("AS_user_%d", rand.IntN(1000)+1)}
	return w.RunUserHandler(ctx, user)
}

func (w *multiWldGen) RunAFHandler(ctx context.Context) workload.Stat {
	user := &workflow.User{Username: fmt.Sprintf("AF_user_%d", rand.IntN(1000)+1)}
	return w.RunUserHandler(ctx, user)
}

func (w *multiWldGen) RunSAHandler(ctx context.Context) workload.Stat {
	user := &workflow.User{Username: fmt.Sprintf("SA_user_%d", rand.IntN(1000)+1)}
	return w.RunUserHandler(ctx, user)
}

func (w *multiWldGen) RunOCHandler(ctx context.Context) workload.Stat {
	user := &workflow.User{Username: fmt.Sprintf("OC_user_%d", rand.IntN(1000)+1)}
	return w.RunUserHandler(ctx, user)
}

func (w *multiWldGen) RunUserHandler(ctx context.Context, user *workflow.User) workload.Stat {
	return statWrapper(func() error {
		var err error
		if *mode == "irid" {
			_, err = w.userservice.GetUserInfo(ctx, *user)
		} else if *mode == "hard1" {
			_, err = w.userservice.GetUserInfoHard1(ctx, *user)
		} else if *mode == "hard2" {
			_, err = w.userservice.GetUserInfoHard2(ctx, *user)
		}
		return err
	})
}

func (w *multiWldGen) RegisterUsers(ctx context.Context) error {
	num_users_per_region := 1000
	var j int64
	j = 0
	regions := []string{"NA", "SA", "EU", "AS", "AF", "OC"}
	for _, region := range regions {
		for i := 1; i <= num_users_per_region; i++ {
			fname, lname, password, email, address := gen_user_data()
			country := gen_random_country(region)
			username := fmt.Sprintf("%s_user_%d", region, i)
			err := w.userservice.RegsiterUser(ctx, username, fname, lname, password, int64(i)+j, email, address, country)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (w *multiWldGen) Run(ctx context.Context) error {
	if *mode != "irid" && *mode != "hard1" && *mode != "hard2" {
		return errors.New("Incorrect mode was provided")
	}
	initialize_arrays()
	if *add_users {
		err := w.RegisterUsers(ctx)
		if err != nil {
			return err
		}
	}
	w1 := workload.NewWorkload()
	w1.AddAPI("NA", w.RunNAHandler, 50)
	w1.AddAPI("EU", w.RunEUHandler, 10)
	w1.AddAPI("SA", w.RunSAHandler, 10)
	w1.AddAPI("AS", w.RunASHandler, 10)
	w1.AddAPI("AF", w.RunAFHandler, 10)
	w1.AddAPI("OC", w.RunOCHandler, 10)
	w2 := workload.NewWorkload()
	w2.AddAPI("NA", w.RunNAHandler, 10)
	w2.AddAPI("EU", w.RunEUHandler, 10)
	w2.AddAPI("SA", w.RunSAHandler, 10)
	w2.AddAPI("AS", w.RunASHandler, 10)
	w2.AddAPI("AF", w.RunAFHandler, 10)
	w2.AddAPI("OC", w.RunOCHandler, 50)

	e1, err := workload.NewEngine(*outfilew1, *tput, *duration, w1)
	if err != nil {
		return err
	}

	e2, err := workload.NewEngine(*outfilew2, *tput, *duration, w2)
	if err != nil {
		return err
	}

	// Run the 1st workload
	e1.RunOpenLoop(ctx)
	// Run the 2nd workload
	e2.RunOpenLoop(ctx)

	err = e1.PrintStats()
	if err != nil {
		return err
	}
	err = e2.PrintStats()
	if err != nil {
		return err
	}
	return nil
}

func (w *multiWldGen) ImplementsMultiWorkload(context.Context) error {
	return nil
}
