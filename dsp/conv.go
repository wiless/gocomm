package dsp

import (
	"fmt"
	"github.com/wiless/vlib"
)

func Conv(in1, in2 vlib.VectorC) (result vlib.VectorC) {
	L1 := in1.Size()
	L2 := in2.Size()
	N := L1 + L2 - 1
	result = vlib.NewVectorC(N)
	fmt.Printf("\n in1=%v", in1)
	fmt.Printf("\n in2=%v", in2)
	for n := 0; n < N; n++ {

		for l := 0; l < L1; l++ {
			indx := n - l
			if indx < L2 && indx >= 0 {
				result[n] += in1[l] * in2[n-l]
			}

		}
		// result[n] = sum

	}
	fmt.Printf("\n result=%v", result)
	return result
}
