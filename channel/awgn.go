package channel

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"wiless/gocomm"
	"wiless/gocomm/chipset"
	"wiless/gocomm/sources"
)

type ChannelEmulator struct {
	noiseDb float64

	/// Chipset Related
	name          string
	isInitialized bool
	Pins          map[string]chipset.PinInfo
	PinNames      map[int]string
	Modules       map[string]chipset.ModuleInfo
	ModuleNames   map[int]string
}

func (m *ChannelEmulator) AWGNChannel(dummy gocomm.Complex128Channel) {
	fmt.Printf("\n Noise ready to Input %v", dummy)
	outCH := m.Pins["symbolOut"].Channel.(gocomm.Complex128Channel)
	fmt.Printf("\n Output ready to Output %v", outCH)
	var chdataOut gocomm.SComplex128Channel
	var chdataIn gocomm.SComplex128Channel
	samples := 1
	// result := make([]complex64, samples)
	var StdDev float64 = math.Sqrt(math.Pow(10, m.noiseDb*.1))
	var Mean float64 = 0
	var noise complex128

	for i := 0; i < samples; i++ {

		chdataIn = <-dummy
		chdataOut.MaxExpected = chdataIn.MaxExpected
		samples = chdataIn.MaxExpected
		fmt.Printf("\nAWGN expects %d samples @ %v", samples, dummy)
		if Mean != 0 && StdDev != 1 {
			noise = complex128(complex(rand.NormFloat64()*StdDev+Mean, rand.NormFloat64()*StdDev+Mean))
		} else {
			noise = complex128(complex(rand.NormFloat64(), rand.NormFloat64()))
		}
		chdataOut.Ch = chdataIn.Ch + noise
		outCH <- chdataOut
	}

}

func GenerateNoise(noiseDb float64, samples int) []complex64 {

	result := make([]complex64, samples)
	var StdDev float64 = math.Sqrt(math.Pow(10, noiseDb*.1))
	var Mean float64 = 0
	if Mean != 0 && StdDev != 1 {
		for i := 0; i < samples; i++ {

			result[i] = complex64(complex(rand.NormFloat64()*StdDev+Mean, rand.NormFloat64()*StdDev+Mean))
		}
	} else {
		for i := 0; i < samples; i++ {
			result[i] = complex64(complex(rand.NormFloat64(), rand.NormFloat64()))

		}
	}
	return result

}

/// Fading/AWGN Channel that operates on each sample
func (m *ChannelEmulator) FadingChannel(InCH gocomm.Complex128Channel) {
	outCH := m.Pins["symbolOut"].Channel.(gocomm.Complex128Channel)
	NextSize := 1
	N0 := .01 /// 10dB SNR
	var chdataOut gocomm.SComplex128Channel
	var chdataIn gocomm.SComplex128Channel
	for i := 0; i < NextSize; i++ {

		chdataIn = <-InCH
		sample := chdataIn.Ch

		/// Do the processing here
		// gain := 1 //sources.RandNC(1)
		noise := sources.RandNC(N0)
		///Fading
		// psample := sample * gain
		psample := sample + noise
		///
		chdataOut.MaxExpected = chdataIn.MaxExpected
		chdataOut.Ch = psample
		outCH <- chdataOut
	}
	// close(InCH)
}

/// CHIPSET interface

func (m ChannelEmulator) IsInitialized() bool {
	return m.isInitialized
}

func (m *ChannelEmulator) InitModules() {
	var totalModules int = m.ModulesCount()
	m.Modules = make(map[string]chipset.ModuleInfo, totalModules)
	m.ModuleNames = make(map[int]string, totalModules)
	// b := [...]string{"Penn", "Teller"}
	strlist := [...]string{"fadingChannel", "awgn"}
	for i := 0; i < totalModules; i++ {
		m.ModuleNames[i] = strlist[i]
	}

	for i := 0; i < totalModules; i++ {
		var minfo chipset.ModuleInfo
		minfo.Name = m.ModuleNames[i]

		switch minfo.Name {
		case "fadingChannel":
			minfo.Desc = "This emulates a 1-tap fading (multiplicative) channel"
			minfo.InPins = []int{0}
			minfo.OutPins = []int{0}
			method := reflect.ValueOf(m).MethodByName("FadingChannel")
			minfo.Function = method

		case "awgn":
			minfo.Desc = "This emulates additive white noise to the input signal "
			minfo.InPins = []int{0}
			minfo.OutPins = []int{0}
			method := reflect.ValueOf(m).MethodByName("AWGNChannel")
			minfo.Function = method

		}
		m.Modules[minfo.Name] = minfo
	}

}
func (m *ChannelEmulator) InitPins() {
	m.isInitialized = true
	totalpins := m.InPinCount() + m.OutPinCount()
	m.Pins = make(map[string]chipset.PinInfo, totalpins)
	m.PinNames = make(map[int]string, totalpins)
	// b := [...]string{"Penn", "Teller"}
	strlist := [4]string{"symbolIn", "symbolOut"}
	for i := 0; i < totalpins; i++ {
		m.PinNames[i] = strlist[i]
	}

	for i := 0; i < totalpins; i++ {
		var pinfo chipset.PinInfo
		// pinfo.CreateComplex128Channel()
		pinfo.Name = m.PinNames[i]
		if i < m.InPinCount() {
			pinfo.InputPin = true
		} else {
			pinfo.InputPin = false
		}
		m.Pins[m.PinNames[i]] = pinfo

	}

	testcch := gocomm.NewComplex128Channel()

	var dummypin chipset.PinInfo

	/// all Input Pins
	dummypin = m.Pins["symbolIn"]
	dummypin.DataType = reflect.TypeOf(testcch)
	m.Pins["symbolIn"] = dummypin

	/// All output pins
	dummypin = m.Pins["symbolOut"]
	dummypin.DataType = reflect.TypeOf(testcch)
	dummypin.CreateComplex128Channel()
	m.Pins["symbolOut"] = dummypin

}

func (m *ChannelEmulator) InitializeChip() {

	m.InitPins()
	m.InitModules()

}

// PinsIn() int
// 	PinsOut() int
// 	Pin(int) PinInfo
// PinsIn() int
// 	PinsOut() int
// 	Pin(int) PinInfo
// 	ModulesCount() int
// 	Module(int) ModuleInfo
func (m ChannelEmulator) InPinCount() int {
	return 1
}

func (m ChannelEmulator) OutPinCount() int {
	return 1
}
func (m ChannelEmulator) Pin(pid int) chipset.PinInfo {
	// result := new(chipset.PinInfo)
	return m.Pins[m.PinNames[pid]]
	// return result
}

func (m ChannelEmulator) PinByName(pinname string) chipset.PinInfo {
	return m.Pins[pinname]
}

func (m ChannelEmulator) PinIn(pid int) chipset.PinInfo {

	return m.Pins[m.PinNames[pid]]

}
func (m ChannelEmulator) PinOut(pid int) chipset.PinInfo {
	return m.Pins[m.PinNames[pid+m.InPinCount()]]

}

// Has Modulator and Demodulator
func (m ChannelEmulator) ModulesCount() int {
	return 2
}
func (m *ChannelEmulator) ModuleByName(mname string) chipset.ModuleInfo {
	return m.Modules[mname]
}

func (m ChannelEmulator) Module(moduleid int) chipset.ModuleInfo {
	return m.ModuleByName(m.ModuleNames[moduleid])

}

func (m ChannelEmulator) Name() string {
	return "ChannelEmulator"
}
