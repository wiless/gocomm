package main

import (
	"fmt"
	"math/rand"

	"time"
	"github.com/wiless/gocomm"
	"github.com/wiless/gocomm/sources"
	"github.com/wiless/vlib"
)

func main() {

	rand.Seed(time.Now().Unix())
	// rand.Seed(1)

	start := time.Now()
	var N = 10
	var N0 float64 = 1.0 / 10
	txbits := vlib.VectorB(sources.RandB(N))

	qpskModem := new(gocomm.Modem)
	var BitsPerSymbol = 2
	qpskModem.Init(BitsPerSymbol)
	fmt.Printf("\n%v", qpskModem)

	/// Take every N-bit and push to modem-channel-demod

	length := len(txbits)
	// slength := length / BitsPerSymbol

	//var result = make([]complex128, slength)
	/// Actual Modulation happens here

	cnt := 0
	// rxcnt := 0
	// symCH := make([]chan complex128, length)
	COUNT := length / 2
	rxCH := make([]chan []uint8, COUNT)

	for i := 0; i < COUNT; i++ {
		rxCH[i] = make(chan []uint8)
	}

	for i := 0; i < length; i += BitsPerSymbol {

		symCH := make(chan complex128)

		go qpskModem.GenerateSymbolCH(cnt, txbits[i:i+BitsPerSymbol], symCH)
		go func(twowayChannel chan complex128, cnt int) {

			symbol := <-twowayChannel
			noise := sources.RandNC(N0)
			fmt.Printf("\n Added Noise %d %v", cnt, noise)
			rxsymbol := symbol + noise
			twowayChannel <- rxsymbol
		}(symCH, cnt)
		go qpskModem.DemodSymbolCH(cnt, symCH, rxCH[cnt])

		cnt++
	}

	ORDER := true
	//// UNORDERED receiver bits...
	for i := 0; i < COUNT; i++ {

		if ORDER {
			func(cnt int, tmpRxCh chan []uint8) {
				for rxbits := range tmpRxCh {
					fmt.Printf("\n Unordered Loop Detected bits %d %v", cnt, rxbits)
				}

			}(i, rxCH[i])
		} else {
			go func(cnt int, tmpRxCh chan []uint8) {
				for rxbits := range tmpRxCh {
					fmt.Printf("\n Unordered Loop Detected bits %d %v", cnt, rxbits)
				}

			}(i, rxCH[i])
			fmt.Printf("Waiting for routines to finish")
			time.Sleep(2 * time.Second)

		}
	}

	/// Ordered read from the Receiver channels
	// cnt = 0
	// for _, channel := range rxCH {

	// 	for rxbits := range channel {
	// 		fmt.Printf("\n Outer Loop Detected bits %d %v", cnt, rxbits)
	// 	}
	// 	cnt++
	// }

	// txsymbols := qpskModem.ModulateBits(txbits)
	// rxsymbols := vlib.ElemAddCmplx(txsymbols, sources.Noise(len(txsymbols), N0))
	// rxbits := qpskModem.DeModulateBits(rxsymbols)

	// // fmt.Printf("\nTx=%v", txbits)
	// // fmt.Printf("\nRx=%v", rxbits)
	// fmt.Printf("\nERR=%v/%d = %f", txbits.CountErrors(rxbits), len(txbits), float64(txbits.CountErrors(rxbits))/float64(len(txbits)))
	fmt.Printf("\n\n Elapsed %v \n", time.Since(start).String())
}
