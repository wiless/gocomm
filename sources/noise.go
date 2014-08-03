package sources

import (
	"math"
	"math/rand"
)

func RandNC(variance float64) complex128 {
	var result complex128
	var Mean float64 = 0
	if Mean != 0 && variance != 1 {
		var StdDev float64 = math.Sqrt(variance)
		result = complex128(complex(rand.NormFloat64()*StdDev+Mean, rand.NormFloat64()*StdDev+Mean))

	} else {
		result = complex128(complex(rand.NormFloat64(), rand.NormFloat64()))
	}
	return result
}

func RandNCVec(samples int, variance float64) []complex128 {

	result := make([]complex128, samples)
	for i := 0; i < samples; i++ {
		result[i] = RandNC(variance)
	}
	return result

}
