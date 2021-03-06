package core

import (
	"fmt"
	"github.com/wiless/gocomm"
	"github.com/wiless/gocomm/chipset"
	"github.com/wiless/gocomm/dsp"
	"github.com/wiless/gocomm/sources"
	"github.com/wiless/vlib"
	"log"
	"math/cmplx"
	"reflect"
)

func init() {
	log.Println("core::MPChannel")
}

type ChannelParam struct {
	Ts           float64
	TimeStamp    float64
	pdp          vlib.VectorF
	Coeff        vlib.VectorC
	Mode         string
	FilterMemory vlib.VectorC
	PowerInDBm   float64
}

func (c *ChannelParam) SetPDP(pdp vlib.VectorF) {
	c.pdp = pdp

}

func NewIIDChannel() ChannelParam {
	var result ChannelParam
	result.Ts = -1
	result.Mode = "IID"
	result.pdp = vlib.NewOnesF(1)
	result.Coeff = vlib.NewOnesC(1)
	return result
}

func (c *ChannelParam) SetFlatAWGN() {

}

func DefaultChannel() ChannelParam {
	var result ChannelParam
	result.Ts = 0
	result.Mode = "AWGN"
	result.pdp = vlib.NewOnesF(1)
	result.PowerInDBm = 0
	result.Coeff = vlib.NewOnesC(1)
	return result
}

func (m *ChannelParam) InitParam(p ChannelParam) {
	*m = p

	if m.Ts == 0 {
		m.Mode = "AWGN"
	}
	m.Validate()

}
func (m *ChannelParam) Validate() {

	if m.Ts == -1 {
		m.Mode = "IID"
	}
	if m.Ts == 0 {
		m.Mode = "AWGN"
	}
	if m.pdp.Size() == 0 && m.Coeff.Size() == 0 {
		m.pdp = vlib.NewOnesF(1)
		m.Coeff = vlib.NewOnesC(1)
	}
	if m.pdp.Size() == 0 && m.Coeff.Size() != 0 {
		m.pdp = vlib.NewOnesF(m.Coeff.Size())
	}
	m.FilterMemory = vlib.NewOnesC(m.Coeff.Size())
	m.TimeStamp = -1
}

type MPChannel struct {
	// Channel related
	ChannelParam
	// Chipset
	name          string
	isInitialized bool
	Pins          map[string]chipset.PinInfo
	Modules       map[string]chipset.ModuleInfo
	ModuleNames   map[int]string
	PinNames      map[int]string

	FeedbackCH gocomm.Complex128AChannel
}

/// AutoGenerated through script

func (m MPChannel) InPinCount() int {
	return 2
}

func (m MPChannel) OutPinCount() int {
	return 2
}

func (m MPChannel) Pin(pid int) chipset.PinInfo {
	return m.Pins[m.PinNames[pid]]
}

func (m MPChannel) PinIn(pid int) chipset.PinInfo {
	if pid >= m.InPinCount() {
		fmt.Printf("%d > No of Input Pins %d", pid, m.InPinCount())
		var result chipset.PinInfo
		result.Id = -1
		return result
	}

	return m.Pins[m.PinNames[pid]]

}
func (m MPChannel) PinByID(pid int) chipset.PinInfo {

	return m.Pins[m.PinNames[pid]]
}

func (m MPChannel) PinOut(pid int) chipset.PinInfo {
	if pid >= m.OutPinCount() {
		fmt.Printf("%d > No of Output Pins %d", pid, m.OutPinCount())
		var result chipset.PinInfo
		result.Id = -1
		return result
	}
	return m.Pins[m.PinNames[pid+m.InPinCount()]]

}

func (m MPChannel) PinByName(pinname string) chipset.PinInfo {

	// panic(fmt.Sprintf("\nMPChannel:PinByName - %s unknown", pinname))
	return m.Pins[pinname]
}

func (m MPChannel) ModulesCount() int {
	return 2
}
func (m MPChannel) ModuleByName(mname string) chipset.ModuleInfo {
	panic(fmt.Sprintf("\nMPChannel:ModuleByName - %s unknown", mname))
	return m.Modules[mname]
}

func (m MPChannel) Module(moduleid int) chipset.ModuleInfo {
	return m.ModuleByName(m.ModuleNames[moduleid])
}

func (m MPChannel) Name() string {
	return "MPChannel"
}

func (m MPChannel) IsInitialized() bool {
	return m.isInitialized
}

func (m *MPChannel) InitializeChip() {
	m.name = "MPChannel"
	m.InitPins()
	m.InitModules()
}

func (m *MPChannel) InitPins() {
	totalpins := m.InPinCount() + m.OutPinCount()
	m.Pins = make(map[string]chipset.PinInfo, totalpins)
	m.PinNames = make(map[int]string, totalpins)
	strlist := [6]string{"inputPin1", "outputPin0", "outputPin1", "inputPin0", "inputPin2", "outputPin2"}
	for cnt := 0; cnt < len(strlist); cnt++ {
		m.PinNames[cnt] = strlist[cnt]
	}

	/// something try begins
	var pinfo chipset.PinInfo

	pinfo.Name = "inputPin0"
	pinfo.Id = 0
	pinfo.InputPin = true
	pinfo.DataType = reflect.TypeOf((*gocomm.FloatChannel)(nil)).Elem()

	m.Pins["inputPin0"] = pinfo

	pinfo.Name = "inputPin1"
	pinfo.Id = 1
	pinfo.InputPin = true
	pinfo.DataType = reflect.TypeOf((*gocomm.Complex128AChannel)(nil)).Elem()

	m.Pins["inputPin1"] = pinfo

	pinfo.Name = "outputPin0"
	pinfo.Id = 2
	pinfo.InputPin = false
	pinfo.DataType = reflect.TypeOf((*gocomm.Complex128Channel)(nil)).Elem()

	pinfo.CreateChannel()

	m.Pins["outputPin0"] = pinfo

	pinfo.Name = "outputPin1"
	pinfo.Id = 3
	pinfo.InputPin = false
	pinfo.DataType = reflect.TypeOf((*gocomm.FloatChannel)(nil)).Elem()
	pinfo.CreateChannel()

	m.Pins["outputPin1"] = pinfo

	pinfo.Name = "inputPin2"
	pinfo.Id = 4
	pinfo.InputPin = true
	pinfo.DataType = reflect.TypeOf((*gocomm.FloatChannel)(nil)).Elem()

	pinfo.CreateChannel()

	m.Pins["inputPin2"] = pinfo

	pinfo.Name = "outputPin2"
	pinfo.Id = 5
	pinfo.InputPin = false
	pinfo.DataType = reflect.TypeOf((*gocomm.Complex128AChannel)(nil)).Elem()

	pinfo.CreateChannel()

	m.Pins["outputPin2"] = pinfo

	// pinfo.Name = "coeffPin"
	// pinfo.Id = 6
	// pinfo.InputPin = false
	// pinfo.DataType = reflect.TypeOf((*gocomm.Complex128AChannel)(nil)).Elem()
	// pinfo.Channel = make(gocomm.Complex128AChannel, 20) /// Last 20 samples will be queued

	// m.Pins["coeffPin"] = pinfo

	/// something try ends

}

func (m *MPChannel) InitModules() {

	var totalModules int = 1

	/// AUTO CODE
	/// something try begins
	var minfo chipset.ModuleInfo
	m.Modules = make(map[string]chipset.ModuleInfo, totalModules)
	m.ModuleNames = make(map[int]string, totalModules)

	strlist := [2]string{"channel", "channelBlock"}
	for cnt := 0; cnt < len(strlist); cnt++ {
		m.ModuleNames[cnt] = strlist[cnt]
	}
	var temp, otemp []int

	minfo.Name = "channel"
	minfo.Id = 0
	minfo.Desc = ""

	temp = append(temp, m.PinByName("inputPin0").Id)

	otemp = append(otemp, m.PinByName("outputPin0").Id, m.PinByName("coeffPin").Id)

	minfo.InPins = temp
	minfo.OutPins = otemp
	m.Modules["channel"] = minfo

	minfo.Name = "channelBlock"
	minfo.Id = 1
	minfo.Desc = ""

	temp = append(temp, m.PinByName("inputPin1").Id)

	otemp = append(otemp, m.PinByName("outputPin1").Id, m.PinByName("coeffPin").Id)

	minfo.InPins = temp
	minfo.OutPins = otemp
	m.Modules["channelBlock"] = minfo

	// minfo.Name = "feedback"
	// minfo.Id = 2
	// minfo.Desc = "internal Feedback to estimators"

	// temp = append(temp, m.PinByName("inputPin2").Id)
	// otemp = append(otemp, m.PinByName("outputPin2").Id)

	// minfo.InPins = temp
	// minfo.OutPins = otemp
	// m.Modules["feedback"] = minfo
	/// AUTO CODE

	m.isInitialized = true
}

func (m *MPChannel) Channel(inputPin0 gocomm.Complex128Channel) {
	/// Read your data from Input channel(s) [inputPin0]
	/// And write it to OutputChannels  [outputPin0]

	outputPin0 := chipset.ToComplexCH(m.Pins["outputPin0"])
	var IdealChObj gocomm.SComplex128AObj

	iters := 1
	for i := 0; i < iters; i++ {
		chData := <-inputPin0
		iters = chData.MaxExpected
		/// Do process here with chData

		outData := m.ChannelFn(chData)
		/// coeff attempt to send
		IdealChObj.Ch = m.Coeff
		IdealChObj.Ts = chData.Ts
		IdealChObj.TimeStamp = chData.TimeStamp
		IdealChObj.Message = chData.Message
		IdealChObj.MaxExpected = chData.MaxExpected
		//fmt.Printf("\n Want to broadcast %#v to %#v", IdealChObj, coeffPin)
		///

		// select {
		// case coeffPin <- IdealChObj:
		// 	// fmt.Printf("\n%f : sent message %v", IdealChObj.TimeStamp, IdealChObj.Ch)
		// default:
		// 	// fmt.Printf("\n%f: no message sent", IdealChObj.TimeStamp)
		// }
		outputPin0 <- outData
	}

}

func (m *MPChannel) SetFeedbackChannel(feedback gocomm.Complex128AChannel) {
	m.FeedbackCH = feedback

	/// Read your data from Input channel(s) [inputPin0]
	/// And write it to OutputChannels  [outputPin0]

	// outputPin2 := chipset.ToComplexACH(m.Pins["outputPin2"])
	// iters := 1
	// var outDataObj gocomm.SComplex128AObj
	// for i := 0; i < iters; i++ {
	// 	chData := <-inputPin2
	// 	iters = chData.MaxExpected
	// 	/// Do process here with chData
	// 	if chData.TimeStamp <= (m.TimeStamp + m.Ts) {
	// 		outDataObj.Ch = m.Coeff
	// 	} else {
	// 		outDataObj.Ch = vlib.NewOnesC(0)
	// 	}

	// 	outDataObj.Ch = m.Coeff
	// 	outDataObj.Ts = chData.Ts
	// 	outDataObj.TimeStamp = chData.TimeStamp
	// 	outDataObj.Message = chData.Message
	// 	outDataObj.MaxExpected = chData.MaxExpected
	// 	outputPin2 <- outDataObj
	// }

}

func (m *MPChannel) ChannelBlock(inputPin1 gocomm.Complex128AChannel) {
	/// Read your data from Input channel(s) [inputPin1]
	/// And write it to OutputChannels  [outputPin1]

	///	outputPin1:=m.Pins["outputPin1"].Channel.(gocomm.<DataType>)
	// iters := 1
	// for i := 0; i < iters; i++ {
	// 	chData := <-[inputPin1]
	// 	iters = chData.MaxExpected
	// 	/// Do process here with chData

	// 	outData:= ChannelBlockFn(chData)
	// 	outData.MaxExpected= ??
	// 	outputPin1 <- outData

	// }

}

func (m *MPChannel) updateCoeff(timestamp float64) {
	/// first time
	generated := false
	if m.TimeStamp == -1 {
		m.Coeff.Resize(m.pdp.Size())
		for i := 0; i < m.pdp.Size(); i++ {
			m.Coeff[i] = sources.RandNC(m.pdp[i])
		}
		generated = true
		m.TimeStamp = 0 /// Unusuall if inbetween the MPchannel timestamp has got RESET !!

	} else {

		/// Existing channel-coeff is valid till m.Timestamp+m.TS,
		valid := timestamp < (m.TimeStamp + m.Ts)

		if !valid {
			/// TRIGGER NEW COEFF
			m.Coeff = vlib.NewVectorC(m.pdp.Size())
			for i := 0; i < m.pdp.Size(); i++ {
				m.Coeff[i] = sources.RandNC(m.pdp[i])
			}

			m.TimeStamp = timestamp
			generated = true
		}

	}

	/// Write new coeff to feedback channel if new was generated
	if m.FeedbackCH != nil && generated {
		var chdata gocomm.SComplex128AObj
		chdata.Ch = m.Coeff
		chdata.TimeStamp = m.TimeStamp
		chdata.Ts = m.Ts
		// fmt.Printf("\n CH:GENERATED @ %v with Coeff : %v, Gain : %v ", timestamp, m.Coeff[0], cmplx.Abs(m.Coeff[0])*cmplx.Abs(m.Coeff[0]))
		m.FeedbackCH <- chdata
	}

}

func (m *MPChannel) ChannelFn(sample gocomm.SComplex128Obj) (result gocomm.SComplex128Obj) {
	/// Read your data from Input channel(s) [inputPin0]
	/// And write it to OutputChannels  [outputPin0]
	// fmt.Printf("\n Channel Param %#v  \n input = %v", m.ChannelParam, sample)

	if m.Ts == -1 {

		m.Coeff.Resize(m.pdp.Size())
		for i := 0; i < m.pdp.Size(); i++ {
			m.Coeff[i] = sources.RandNC(m.pdp[i])
		}
		if m.FeedbackCH != nil {
			var chdata gocomm.SComplex128AObj
			chdata.Ch = m.Coeff
			chdata.TimeStamp = m.TimeStamp
			chdata.MaxExpected = sample.MaxExpected
			chdata.Ts = m.Ts
			fmt.Printf("\n FAST IID generated for sample@%v with %v@%v", sample.TimeStamp, m.TimeStamp, cmplx.Abs(m.Coeff[0])*cmplx.Abs(m.Coeff[0]))
			m.FeedbackCH <- chdata
		}

	} else {

		m.updateCoeff(sample.TimeStamp)

	}
	if m.FilterMemory.Size() != m.Coeff.Size() {
		m.FilterMemory.Resize(m.Coeff.Size())
	}

	result = sample /// Carefull if not same ChType
	// m.TimeStamp = sample.TimeStamp
	m.FilterMemory = m.FilterMemory.ShiftLeft(sample.Ch)

	//dummy := vlib.ElemMultC(m.Coeff, vlib.Conj(m.Coeff))
	foutput := vlib.DotC(m.Coeff, m.FilterMemory)
	// fmt.Printf("\n CHANNEL @%v: I/P %v - Gain : %v  : O/p : %v", sample.TimeStamp, sample.Ch, cmplx.Abs(m.Coeff[0])*cmplx.Abs(m.Coeff[0]), foutput)
	result.Ch = foutput
	result.Message = sample.Message + " Filter"
	result.Ts = sample.Ts
	result.TimeStamp = sample.TimeStamp
	result.MaxExpected = sample.MaxExpected
	// fmt.Printf("\n I/O (hn =%v) : %#v --> %#v", m.Coeff, sample, result)
	return result

}

func (m *MPChannel) ChannelBlockFn(sample gocomm.SComplex128AObj) (result gocomm.SComplex128AObj) {

	result.Ch = dsp.Conv(sample.Ch, m.Coeff)
	return result
}
