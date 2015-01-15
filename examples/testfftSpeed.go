package main

import (
	"fmt"
	// "log"
	"math"
	"runtime"
	// "github.com/gonum/stat"

	"time"
	"github.com/wiless/gocomm/dsp"

	"github.com/mjibson/go-dsp/fft"
	"github.com/wiless/gocomm/sources"
	"github.com/wiless/vlib"
)

func main() {
	NFFT := 2048 * 4
	x := sources.RandNCVec(NFFT, 1)
	// fmt.Printf("\n x=%f", x)
	runtime.GOMAXPROCS(8)
	fmt.Printf("\n FFT Point  %v", NFFT)
	fmt.Printf("\nStatus %v", runtime.GOMAXPROCS(-1))
	t1 := time.Now()
	fmt.Printf("\n Start : %v", t1)
	y := gocomm.FFT_C(vlib.VectorC(x), NFFT)
	fmt.Printf("\n Vanilla Elapsed : %v", time.Since(t1).String())
	fmt.Printf("\n y1=%f", y[0])

	fmt.Print("\n\n===== github.com/wiless/gocomm/fft ========== \n\n")
	t2 := time.Now()
	fmt.Printf("\n Start : %v", t2)
	y = gocomm.GoFFT_C(vlib.VectorC(x), NFFT)
	fmt.Printf("\n Concurrent Elapsed : %s", time.Since(t2).String())

	fmt.Printf("\n y2=%f", y[0])

	fmt.Print("\n\n===== mjibson/go-dsp/fft ========== \n\n")
	t3 := time.Now()
	fmt.Printf("\n Start : %v", t3)
	// y = gocomm.GoFFT_C(vlib.VectorC(x), NFFT)
	y = fft.FFT(x)
	y = y.Scale(math.Sqrt(1.0 / float64(NFFT)))
	fmt.Printf("\n Concurrent Elapsed : %s", time.Since(t3).String())

	fmt.Printf("\n y3=%f", y[0])
	/// End of code
	fmt.Printf("\n")
}
