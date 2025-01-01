package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/vaastav/iridescent/iridescent_rt/specrt"
)

func main() {
	fmt.Println("Mat mul main")
	srt, err := specrt.NewSpecializationRuntime(context.Background(), "guest/mmul.go")
	if err != nil {
		log.Fatal(err)
	}
	var mat_mul_fn func(int, int, int, int) (float64, error)
	update_fn := func() error {
		log.Println("Updating symbols")
		mat_mul, err := srt.Lookup("MatrixMultiply")
		if err != nil {
			return err
		}
		var ok bool
		mat_mul_fn, ok = mat_mul.(func(int, int, int, int) (float64, error))
		if !ok {
			return errors.New("Failed to convert symbol into type")
		}
		return nil
	}
	update_fn()
	srt.AddCallbackFn(update_fn)
	val, err := mat_mul_fn(32, 32, 32, 4)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(val)
	srt.Specialize("s", 0)
	err = srt.UpdatePlugin()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Post update")
	val, err = mat_mul_fn(32, 32, 32, 2)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(val)
	srt.Specialize("s", 1)
	err = srt.UpdatePlugin()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Post 2nd update")
	val, err = mat_mul_fn(32, 32, 32, 4)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(val)
}
