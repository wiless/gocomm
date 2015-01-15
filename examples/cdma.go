package main

import (
	"fmt"
	"time"
	"github.com/wiless/gocomm"

	"github.com/wiless/gocomm/customchips"
	"github.com/wiless/gocomm/sources"
	"github.com/wiless/vlib"
)

func main() {

	N := 20 /// 20 samples
	L := 4  /// 5tap channel
	begin := time.Now()

	var cdma customchips.CDMA
	cdma.InitializeChip()
	cdma.SetSpreadCode(vlib.NewOnesC(L))

	samples := vlib.VectorC(sources.RandNCVec(N, 1))
	var data gocomm.SComplex128Obj
	/// METHOD A
	data.Ts = 1
	data.TimeStamp = 0
	data.MaxExpected = N
	data.Message = ""
	for i := 0; i < N; i++ {
		data.Next(samples[i])
		chips := cdma.SpreadFn(data)

		output := cdma.DeSpreadFn(chips)

		fmt.Printf("\nTxSymbol %v ", data)
		// fmt.Printf("\nTx %v ", chips)
		fmt.Printf("\nRxSymbol %v ", output)
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

	//fmt.Printf("\nFilter Residues %v", filter.FilterMemory)

	//  Of code
	fmt.Printf("\nTime Elapsed : %v\n", time.Since(begin))
}
