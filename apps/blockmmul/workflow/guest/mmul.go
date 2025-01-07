package main

func MatrixMultiply(A [][]int, B [][]int, C [][]int, n int, s int) error {
	s = iridescent_instr_general_int("s", s, 2, 4, 8, 16, 32)

	en := s * (n / s)

	for kk := 0; kk < en; kk += s {
		for jj := 0; jj < en; jj += s {
			for i := 0; i < n; i += 1 {
				for jjj := 0; jjj < s; jjj += 1 {
					j := jj + jjj
					sum := C[i][j]
					for kkk := 0; kkk < s; kkk += 1 {
						k := kk + kkk
						sum += A[i][k] * B[k][j]
					}
					C[i][j] = sum
				}
			}
		}
	}

	return nil
}
