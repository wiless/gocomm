package main

import (
	"fmt"
	"time"
	"wiless/gocomm"
	"wiless/gocomm/customchips"
	// "wiless/gocomm/dsp"
	"wiless/gocomm/sources"
)

func main() {

	N := 2048
	begin := time.Now()
	var ofdm customchips.OFDM
	ofdm.InitializeChip()
	ofdm.NPoint = N

	// fmt.Printf("\n OFDM = %v", ofdm)
	/// Input
	inCHA := gocomm.NewComplex128AChannel()
	// inCH := gocomm.NewComplex128Channel()
	var dataArray gocomm.SComplex128AObj

	x := sources.RandNCVec(N, 1)
	dataArray.Ch = x
	dataArray.MaxExpected = 1
	inCHA <- dataArray
	go ofdm.Ifft(inCHA)
	go ofdm.Fft(ofdm.PinByName("outputPin0").Channel.(gocomm.Complex128AChannel))
	outCHA := ofdm.PinByName("outputPin1").Channel.(gocomm.Complex128AChannel)
	data := <-outCHA
	/// Analyse
	fmt.Printf("\nX=%f", x[0:16])
	fmt.Printf("\nOutput=%f", data.Ch[0:16])
	// xcap := dsp.ExtFFT_C(data.Ch, N)
	// fmt.Printf("\nXcap=%f", xcap[0:16])

	fmt.Printf("\nTime Elapsed : %v\n", time.Since(begin))
}
