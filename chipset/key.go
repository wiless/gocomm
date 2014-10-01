package chipset

import (
	"math/rand"
)

var CharSet []byte

func init() {
	alphas := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	CharSet = []byte(alphas)
}
func RandIndentifier(N int) string {
	result := make([]byte, N)
	L := len(CharSet)
	for i := 0; i < N; i++ {
		n := rand.Intn(L)
		result[i] = CharSet[n]
	}
	return string(result)
}
