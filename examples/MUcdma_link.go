package main

import (
	"encoding/json"
	"fmt"
	"github.com/wiless/gocomm"
	"github.com/wiless/gocomm/channel"
	"github.com/wiless/gocomm/chipset"
	"github.com/wiless/gocomm/core"
	"github.com/wiless/gocomm/dsp"
	"github.com/wiless/gocomm/sources"
	"github.com/wiless/vlib"
	"io/ioutil"
	"log"
	"runtime"
	"sort"
	"sync"
	"time"
)

func main() {

	fmt.Print("GOPROCS:=", runtime.GOMAXPROCS(6))
	runtime.SetCPUProfileRate(-1)
	start := time.Now()

	user1 := NewSetup()
	user2 := NewSetup()
	user3 := NewSetup()
	// user4 := NewSetup()
	// user5 := NewSetup()
	// fmt.Printf("\nLink %v", user1)
	// fmt.Printf("\nLink %v", user2)
	data, err := ioutil.ReadFile("settings.json")
	if err != nil {
		log.Print("Unable to Read File : ", err)
	}
	result := chipset.GetMetaInfo(data, "Modem1")
	fmt.Print("Found Setting : ", result, "len = ", len(result))
	var jsons string = `{"NBlocks":100,"snr":"0:2:16","SF":1}`
	var mymodem core.Modem
	mymodem.SetName("Modem2")
	mymodem.SetJson(data)

	fmt.Print("SOMETHING", string(mymodem.GetJson()))
	user1.Set(jsons)
	user2.Set(jsons)
	user3.Set(jsons)
	// user3.Set(jsons)
	// user4.Set(jsons)
	// user5.Set(jsons)
	fmt.Printf("Starting simulation ...")

	fmt.Printf("\n started user 1")
	go user1.Run()
	fmt.Printf("\n started user 2")
	go user2.Run()
	// go user3.Run()
	// go user4.Run()

	// 	user1.Run()
	// 	user2.Run()
	// 	user3.Run()
	// 	user4.Run()
	fmt.Printf("\n started user 3")
	user3.Run()

	// time.Sleep(10 * time.Second)
	fmt.Print("\n Elapsed : ", time.Since(start))

}

var setupid string

// var (
// 	BlockSize  int
// 	NBlocks    = 200
// 	snr_ber    map[float64]float64
// 	s.snr_block  map[float64]float64
// 	snr_ber1   map[float64]float64
// 	s.snr_block1 map[float64]float64
// 	snr        vlib.VectorF
// 	SF         int
// )

func init() {
	setupid = chipset.RandIndentifier(8)
}

var settingVariable map[string]interface{}

type setup struct {
	uid int

	name     string
	chips    []chipset.Chip
	settings string
	// steupSetting settingVariable
	wg *sync.WaitGroup

	BlockSize  int
	NBlocks    int
	snr_ber    map[float64]float64
	snr_block  map[float64]float64
	snr_ber1   map[float64]float64
	snr_block1 map[float64]float64
	snr        vlib.VectorF
	SF         int
	Results    string
}

func NewSetup() *setup {
	var result = new(setup)
	result.name = chipset.RandIndentifier(8)
	result.SetDefaults()

	return result
}

func (s *setup) String() string {
	return fmt.Sprintf("id:%s", s.name)
}

func (s *setup) SetDefaults() {
	s.NBlocks = 200
	s.BlockSize = 1000
}

func (s *setup) Set(jsons string) {
	s.settings = jsons
	var v map[string]interface{}
	json.Unmarshal([]byte(jsons), &v)
	fmt.Printf("\nJSON Object : %v", v)

	s.NBlocks = int(v["NBlocks"].(float64))
	s.SF = int(v["SF"].(float64))
	s.snr = vlib.ToVectorF(v["snr"].(string))
	s.wg = new(sync.WaitGroup)

}

func (s *setup) Get() string {

	return s.settings
}

func (s *setup) Run() {
	s.snr_ber = make(map[float64]float64)
	s.snr_block = make(map[float64]float64)
	s.snr_ber1 = make(map[float64]float64)
	s.snr_block1 = make(map[float64]float64)
	hn := vlib.NewOnesF(1)

	spcode := vlib.NewOnesC(s.SF)
	// snr = vlib.VectorF{0, 10, 100}
	// snr = vlib.ToVectorF("0:2:20")
	fmt.Printf("\nSNR : %v", s.snr)
	fmt.Printf("\nhn : %v", hn)
	// linkresult := make(map[float64]float64, 1)
	// linkresult1 := make(map[float64]float64, 1)

	outCH := make([]gocomm.FloatChannel, s.snr.Size())
	// outCH1 := make([]gocomm.FloatChannel, snr.Size())
	fmt.Printf("\n SNR ")
	for i := 0; i < len(s.snr); i++ {
		// blocks := s.snr_block[snr[i]] + 1
		fmt.Printf("%2.2f ", s.snr[i])
	}

	fmt.Printf("\n")
	for i := 0; i < s.snr.Size(); i++ {
		// s.wg.Add(1)
		outCH[i] = gocomm.NewFloatChannel()
		// outCH1[i] = gocomm.NewFloatChannel()
		s.wg.Add(1)
		go s.SimulateLinkFn(s.snr[i], hn, spcode, outCH[i], 0)
		// go s.SimulateLinkFn(snr[i]-2, hn, spcode, outCH1[i], 1)

	}

	s.wg.Wait()

	log.Printf("\nSNR : %v", s.snr)
	log.Printf("\nBER : %v", s.snr_ber)
	log.Printf("\nhn : %v", hn)
	s.Results = string(Print(s.snr_ber, "", ""))
	log.Printf("Result %v", s.Results)
}

func (s *setup) SimulateLinkFn(SNR float64, pdp vlib.VectorF, spcode vlib.VectorC, outCH gocomm.FloatChannel, uid int) {
	var floatobj gocomm.SFloatObj
	floatobj.Ch = 0

	for j := 0; j < s.NBlocks; j++ {
		s.snr_block[SNR] = float64(j)
		result := s.SimulateLink(SNR, pdp, spcode, uid)
		floatobj.Ch += result
	}
	floatobj.Ch /= float64(s.NBlocks)
	outCH <- floatobj
	s.wg.Done()
}
func (s *setup) SimulateLink(SNR float64, pdp vlib.VectorF, spcode vlib.VectorC, uid int) float64 {
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

	bitsample := vlib.VectorB(sources.RandB(s.BlockSize))
	bitChannel := gocomm.NewBitChannel()
	var dataArray gocomm.SBitObj
	dataArray.MaxExpected = bitsample.Size()
	dataArray.Message = fmt.Sprintf("Src %f  ", SNR)
	dataArray.Ts = bitTs
	dataArray.MaxExpected = s.BlockSize
	// fmt.Printf("bits=%v", bitsample)

	go (func(bitChannel gocomm.BitChannel) {
		for i := 0; i < s.BlockSize; i++ {
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
	ch7 := chipset.ToFloatCH(BER.PinByName("outputPin0"))
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
	if uid == 0 {
		s.snr_ber[SNR] += result
	} else {
		s.snr_ber1[SNR] += result
	}

	// snr_ber[SNR] /= s.snr_block[SNR]
	s.updateTable()
	// fmt.Printf("\r BLOCK %v \n BER : %v", s.snr_block, snr_ber)
	return result
}

func Print(xyvec map[float64]float64, xlabel string, ylabel string) []byte {
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
	type temp struct {
		SNR vlib.VectorF
		BER vlib.VectorF
	}
	data := temp{x, y}
	//keys := []float64(x)
	fmt.Printf("\n%s=%1.2e\n %s=%1.2e", xlabel, x, ylabel, y)
	result, err := json.Marshal(data)
	if err != nil {
		return nil
	} else {
		return result
	}

	//fmt.Printf("\n%s=%f", xlabel, x)
	//fmt.Printf("\n%s=%f", ylabel, y)
}

func (s *setup) updateTable() {
	// fmt.Printf("\r === ADB === %v", snr_ber[0])
	str := ""
	// str1 := ""
	for _, val := range s.snr {
		blocks := s.snr_block[val] + 1
		// blocks1 := s.s.snr_block1[snr[i]-2] + 1
		// str += fmt.Sprintf("%2.2e (%d) ", snr_ber[snr[i]]/blocks, int(s.snr_block[snr[i]])+1)
		str += fmt.Sprintf("%2.2e ", s.snr_ber[val]/blocks)
		// str1 += fmt.Sprintf("%2.2e ", snr_ber1[snr[i]]/blocks1)
	}
	// fmt.Printf("\r BER(block) : %s \t %s", str, str1)
	fmt.Printf("\r BER (%s) : %s ", s.name, str)
}
