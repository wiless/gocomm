package main

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
	"time"
	"wiless/gocomm"
	"wiless/gocomm/channel"
	"wiless/gocomm/chipset"
	"wiless/gocomm/core"
	"wiless/gocomm/dsp"
	"wiless/gocomm/sources"
	"wiless/vlib"
)

var wg sync.WaitGroup
var BlockSize int
var NBlocks = 200

func main() {

	runtime.GOMAXPROCS(5)
	BlockSize = 1000 // 20 samples
	SF := 1
	begin := time.Now()
	hn := vlib.NewOnesF(1)
	spcode := vlib.NewOnesC(SF)
	snr := vlib.VectorF{0, 10, 100}
	snr = vlib.ToVectorF("0:2:20")
	fmt.Printf("\nSNR : %v", snr)
	fmt.Printf("\nhn : %v", hn)
	linkresult := make(map[float64]float64, 1)

	outCH := make([]gocomm.FloatChannel, snr.Size())
	for i := 0; i < snr.Size(); i++ {
		// wg.Add(1)
		outCH[i] = gocomm.NewFloatChannel()
		go SimulateLinkFn(snr[i], hn, spcode, outCH[i])

	}
	// wg.Wait()

	for i := 0; i < snr.Size(); i++ {
		data := <-outCH[i]
		linkresult[snr[i]] = data.Ch
	}
	Print(linkresult, "SNR", "BER")
	//for i := 0; i < snr.Size(); i++ {
	//	fmt.Printf("\n%v = %v", snr[i], linkresult[snr[i]])
	//}

	//fmt.Printf("\nLinkResult %v", linkresult)
	/// Actual data pushing
	// outputPin := channel.PinByID(1)
	//fmt.Printf("\nFilter Residues %v", channel.FilterMemory)

	//  Of code
	fmt.Printf("\nTime Elapsed : %v\n", time.Since(begin))
}

func SimulateLinkFn(SNR float64, pdp vlib.VectorF, spcode vlib.VectorC, outCH gocomm.FloatChannel) {
	var floatobj gocomm.SFloatObj
	floatobj.Ch = 0
	for j := 0; j < NBlocks; j++ {
		result := SimulateLink(SNR, pdp, spcode)
		floatobj.Ch += result
	}
	floatobj.Ch /= float64(NBlocks)
	outCH <- floatobj
	// wg.Done()
}
func SimulateLink(SNR float64, pdp vlib.VectorF, spcode vlib.VectorC) float64 {
	var cdma core.CDMA
	bitTs := 1.0
	cdma.InitializeChip()

	cdma.SetSpreadCode(spcode, true)
	var mpchannel core.MPChannel
	mpchannel.InitializeChip()
	///
	param := core.NewIIDChannel()
	param.SetPDP(pdp)
	param.Mode = "iid"
	symbolTs := 2 * bitTs
	param.Ts = 5 * symbolTs // bitTs // dataArray.Ts * float64(BlockSize)
	mpchannel.InitParam(param)

	var modem core.Modem

	modem.InitializeChip()
	modem.InitModem(2)

	feedback := make(gocomm.Complex128AChannel, 10)
	mpchannel.SetFeedbackChannel(feedback)
	modem.SetFeedbackChannel(feedback)

	bitsample := vlib.VectorB(sources.RandB(BlockSize))
	bitChannel := gocomm.NewBitChannel()
	var dataArray gocomm.SBitObj
	dataArray.MaxExpected = bitsample.Size()
	dataArray.Message = fmt.Sprintf("Src %f  ", SNR)
	dataArray.Ts = bitTs
	dataArray.MaxExpected = BlockSize
	// fmt.Printf("bits=%v", bitsample)

	go (func(bitChannel gocomm.BitChannel) {
		for i := 0; i < BlockSize; i++ {
			dataArray.Ch = bitsample[i]
			// fmt.Printf("\n Transmitting .. %v", dataArray)
			bitChannel <- dataArray

			// modem.ModulateFn(dataArray)
			// symCH <- dataArray
			dataArray.TimeStamp += dataArray.Ts
		}
	})(bitChannel)

	go modem.Modulate(bitChannel)
	// fmt.Printf("Generated Bits", samples)
	symCH := chipset.ToComplexCH(modem.PinByName("outputPin0"))
	go cdma.Spread(symCH)

	OutCH := chipset.ToComplexACH(cdma.PinByID(2))
	ch1 := gocomm.NewComplex128Channel()
	go gocomm.ComplexA2Complex(OutCH, ch1)
	/// Add Multipath channel
	go mpchannel.Channel(ch1)
	ch2 := chipset.ToComplexCH(mpchannel.PinByName("outputPin0"))
	/// Add Noise
	noiseDb := 1.0 / dsp.InvDb(SNR)
	// fmt.Printf("\n %f AWGN Noise : %f", SNR, noiseDb)
	var awgn channel.ChannelEmulator
	awgn.SetNoise(0, noiseDb)
	awgn.InitializeChip()

	go awgn.AWGNChannel(ch2)
	ch3 := chipset.ToComplexCH(awgn.PinByID(1))

	// go awgn. (1/noiseDb, samples)

	go cdma.DeSpread(ch3)
	ch4 := chipset.ToComplexCH(cdma.PinByID(3))
	//go sink.CROremote(cdma.PinByID(3))

	go modem.DeModulate(ch4)
	ch5 := chipset.ToComplexCH(modem.PinByID(3)) // .Channel.(gocomm.Complex128Channel)
	// gocomm.WGroup.Add(1)

	// gocomm.WGroup.Add(1)
	// go gocomm.SinkComplex(ch4, "")
	// ///
	// bitsample
	var BER core.BER
	BER.TrueBits = bitsample
	BER.InitializeChip()
	ch6 := gocomm.NewBitChannel()
	BER.Reset()
	go BER.BERCount(ch6)
	go (func(ch6 gocomm.BitChannel) {
		count := 1
		for i := 0; i < count; i++ {
			data := <-ch5
			count = data.MaxExpected
			//fmt.Printf("\n sym received %#v", data)
			rxbits := gocomm.Complex2Bits(data)
			// fmt.Printf("\n bits converted received %#v", rxbits)
			ch6 <- rxbits[0]
			ch6 <- rxbits[1]
			// dec := gocomm.Complex2Bits(data)
			// result[2*i] = dec[0]
			// result[2*i+1] = dec[1]

			// fmt.Printf("\n %v", gocomm.Complex2Bits(data))

		}

	})(ch6)
	ch7 := BER.PinByName("outputPin0").Channel.(gocomm.FloatChannel)
	count := 1
	// result := vlib.NewVectorB(BlockSize)
	var result float64 = 0
	for i := 0; i < count; i++ {

		err := <-ch7
		count = err.MaxExpected
		// fmt.Printf("\nerr=(%f,%v) %#v", err.TimeStamp, err.Ch, err)
		result = err.Ch
	}
	// err := float64(bitsample.CountErrors(result)) / float64(BlockSize)
	// fmt.Printf("\ntxbits=%v", bitsample)
	// fmt.Printf("\nrxbits=%v", result)
	// fmt.Printf("\nErr=%v %v", bitsample.CountErrors(result), err)
	// gocomm.WGroup.Wait()

	return result
}

//type xyvec map[float64]float64

//for key, value := range m {
//    fmt.Println("Key:", key, "Value:", value)
//}
func Print(xyvec map[float64]float64, xlabel string, ylabel string) {
	x := vlib.NewVectorF(len(xyvec))
	y := vlib.NewVectorF(len(xyvec))
	cnt := 0
	for vx, _ := range xyvec {
		x[cnt] = vx
		cnt++
	}

	sort.Float64s([]float64(x))

	for indx, vx := range x {
		y[indx] = xyvec[vx]
	}

	//keys := []float64(x)
	fmt.Printf("\n%s=%e\n %s=%e", xlabel, x, ylabel, y)
	//fmt.Printf("\n%s=%f", xlabel, x)
	//fmt.Printf("\n%s=%f", ylabel, y)
}
