package main

import (
	"log"
	"math/rand/v2"

	"gonum.org/v1/gonum/mat"
)

func MatrixMultiply(m int, n int, k int, s int) (float64, error) {
	m = iridescent_instr_general_int("m", m, 32, 64, 128)
	n = iridescent_instr_general_int("n", n, 32, 64, 128)
	k = iridescent_instr_general_int("k", k, 32, 64, 128)
	s = iridescent_instr_general_int("s", s, 4, 2)
	log.Println(s)

	a_data := make([]float64, m*k)
	b_data := make([]float64, k*n)

	for i := range a_data {
		a_data[i] = rand.NormFloat64()
	}
	for i := range b_data {
		b_data[i] = rand.NormFloat64()
	}

	A := mat.NewDense(m, k, a_data)
	B := mat.NewDense(k, n, b_data)

	var C mat.Dense
	C.Mul(A, B)

	return C.Norm(1.0), nil
}
