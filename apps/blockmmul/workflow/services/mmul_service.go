package services

import (
	"context"
	"errors"
	"math/rand"

	"github.com/vaastav/iridescent/iridescent_rt/autotune"
)

type MatrixMulService interface {
	MatrixMultiply(ctx context.Context, n int) error
}

type MatrixMulServiceImpl struct {
	MatMulFn func(a [][]int, b [][]int, c [][]int, n int, s int) error
}

func NewMatrixMulServiceImpl(ctx context.Context) (MatrixMulService, error) {
	impl := &MatrixMulServiceImpl{}
	update_fn := func() error {
		srt := autotune.GetRuntime().SpecRT
		mat_mul, err := srt.Lookup("MatrixMultiply")
		if err != nil {
			return err
		}
		var ok bool
		impl.MatMulFn, ok = mat_mul.(func([][]int, [][]int, [][]int, int, int) error)
		if !ok {
			return errors.New("Failed to convert loaded symbol into desired type")
		}
		return nil
	}
	err := update_fn()
	if err != nil {
		return nil, err
	}
	autotune.GetRuntime().SpecRT.AddCallbackFn(update_fn)
	return impl, nil
}

func (m *MatrixMulServiceImpl) MatrixMultiply(ctx context.Context, n int) error {
	A := make([][]int, n)
	B := make([][]int, n)
	C := make([][]int, n)
	for i := 0; i < n; i++ {
		A[i] = make([]int, n)
		B[i] = make([]int, n)
		C[i] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			A[i][j] = int(rand.Int31())
			B[i][j] = int(rand.Int31())
		}
	}

	return m.MatMulFn(A, B, C, n, 4)
}
