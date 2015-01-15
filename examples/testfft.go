package main

import (
	"fmt"
	// "math"
	// "github.com/gonum/stat"
	"github.com/wiless/gocomm/funcs"
	"github.com/wiless/gocomm/modem"
	"github.com/wiless/gocomm/sources"
	// "github.com/wiless/vlib"
)

func main() {
	NFFT := 64
	// x := sources.RandNCVec(NFFT, 1)
	// y := gocomm.FFT_C(vlib.VectorC(x), NFFT)
	// fmt.Printf("\n x=%f", x)
	// fmt.Printf("\n y=%f", y)
	// fmt.Printf("\n xcap=%f", gocomm.IFFT_C(y, NFFT))

	txmodem := modem.NewModem(2)

	txmodem.InitializeChip()
	fmt.Printf("\nconstellation=%f\n", txmodem.Constellation)
	bits := sources.RandB(NFFT * 2)
	txsymbols := txmodem.ModulateBits(bits)
	txOFDM := gocomm.IFFT_C(txsymbols, NFFT)

	fmt.Printf("\ntxbits=%d1", bits)
	fmt.Printf("\ntxsymbols=%f", txsymbols)
	fmt.Printf("\ntxofdm=%f", txOFDM)

	/// End of code
	fmt.Printf("\n")
}
