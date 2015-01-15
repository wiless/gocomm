package main

import (
	"fmt"
	"time"
	"github.com/wiless/gocomm"
	// "github.com/wiless/gocomm/chipset"
	"github.com/wiless/gocomm/customchips"
	"github.com/wiless/gocomm/sources"
	"github.com/wiless/vlib"
)

func main() {

	N := 20 /// 20 samples
	L := 5  /// 5tap channel
	begin := time.Now()

	var filter customchips.Filter
	filter.InitializeChip()
	filter.SetCoeff(vlib.NewOnesC(L))

	// fmt.Printf("\n OFDM = %v", ofdm)
	/// Input
	// inCH := gocomm.NewComplex128Channel()
	samples := vlib.VectorC(sources.RandNCVec(N, 1))
	var dataArray gocomm.SComplex128Obj
	/// METHOD A
	for i := 0; i < N; i++ {
		dataArray.Ch = samples[i]
		fmt.Printf("\n%d I/O : %f ==> %f", i, dataArray.Ch, filter.FilterFn(dataArray).Ch)
	}

	/// METHOD B
	// dataArray.MaxExpected = samples.Size()
	// inCHA := gocomm.NewComplex128Channel()
	// outputPin := filter.PinByID(1)

	// go filter.Filter(inCHA)
	// go chipset.Sink(outputPin)
	// /// Actual data pushing
	// for i := 0; i < N; i++ {
	// 	dataArray.MaxExpected = N
	// 	dataArray.Ch = samples[i]
	// 	inCHA <- dataArray
	// }

	fmt.Printf("\nFilter Residues %v", filter.FilterMemory)

	//  Of code
	fmt.Printf("\nTime Elapsed : %v\n", time.Since(begin))
}
