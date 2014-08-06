package modem

import (
	"fmt"
	"math"
	"math/cmplx"
	"reflect"
	"strconv"

	// "time"
	"wiless/gocomm"
	"wiless/gocomm/chipset"
	"wiless/gocomm/sources"
)

type ModemIterface interface {
	// SetConstellationTable(ctable map[string]complex128)
	// GetConstellationTable() map[string]complex128
	//MinDistanceKey(symbol complex128) string
	String() string
}

type Modem struct {
	bitsPerSymbol      int
	size               int
	name               string
	offset             float64
	keys               []string
	constellationTable map[string]complex128
	Constellation      []complex128
	isInitialized      bool
	Pins               map[string]chipset.PinInfo
	PinNames           map[int]string
	Modules            map[string]chipset.ModuleInfo
	ModuleNames        map[int]string
}

func NewModem(size int) Modem {
	var result Modem
	result.Init(size)
	result.InitializeChip()
	return result
}

type MPSKModem struct {
	Modem
}

func (m *Modem) SetConstellationTable(ctable map[string]complex128) {
	m.constellationTable = ctable
	m.size = len(ctable)
	m.bitsPerSymbol = int(math.Log2(float64(m.size)))
	fmt.Printf("%v", m)
}

func (m *Modem) GetConstellationTable() map[string]complex128 {

	return m.constellationTable
}

func (m *Modem) PrintTable() {
	fmt.Printf("%v", m.constellationTable)

}

func toStr(bits []uint8) string {
	var result string
	size := len(bits)
	for i := 0; i < size; i++ {
		result += fmt.Sprintf("%d", bits[i])
		// result += strconv.Itoa(int(bits[i]))
		// result += strconv.FormatInt(int64(bits[i]), 10)
	}
	return result
}

func (m Modem) MinDistanceKey(symbol complex128) string {
	var minDistance float64 = 100000.0
	result := ""
	for i := 0; i < m.size; i++ {
		currentDistance := cmplx.Abs(m.Constellation[i] - symbol)
		if minDistance >= currentDistance {
			minDistance = currentDistance
			result = m.keys[i]
		}

	}
	return result
}

func (m *Modem) DeModulateBits(symbol complex128) []uint8 {
	blength := m.bitsPerSymbol
	var result = make([]uint8, blength)
	/// Actual Modulation happens here
	N := m.bitsPerSymbol
	cnt := 0

	var num uint8 = 0
	str := m.MinDistanceKey(symbol)
	// fmt.Printf("\nRx %v %d %v ", symbols[i], len(str), str)

	for bit := 0; bit < N; bit++ {
		num = 0
		if str[bit] == '1' {
			num = 1
		}
		result[cnt] = num
		cnt++
	}

	return result
}

func (m *Modem) GetOutputBlockSize(inputPutBlockSize int) int {
	return inputPutBlockSize / m.bitsPerSymbol
}

func (m *Modem) ModulateBlock(OutBlockSize int, bitchan gocomm.BitChannel, symbolChannel gocomm.Complex128Channel) {

	for i := 0; i < OutBlockSize; i++ {
		var chdataIn gocomm.SBitChannel
		var chdataOut gocomm.SComplex128Channel

		// fmt.Printf("\n MaxSymbols expected is %d , message = %v", bitchan.MaxExpected, bitchan.Message)

		var bits []uint8
		bits = make([]uint8, m.bitsPerSymbol)
		length := m.bitsPerSymbol
		N := m.bitsPerSymbol

		for i := 0; i < length; i++ {
			chdataIn = <-bitchan
			chdataOut.MaxExpected = chdataIn.MaxExpected / 2
			chdataOut.Message = chdataIn.Message
			bits[i] = chdataIn.Ch
			OutBlockSize = chdataIn.MaxExpected

		}

		key := toStr(bits[0:N])
		chdataOut.Ch = m.constellationTable[key]

		symbolChannel <- chdataOut
	}
	// close(bitchan.Ch)
}

func (m *Modem) DeModulateBlock(OutBlockSize int, InCH gocomm.Complex128Channel, outCH gocomm.Complex128Channel) {
	var chdataIn gocomm.SComplex128Channel
	var chdataOut gocomm.SComplex128Channel

	for i := 0; i < OutBlockSize; i++ {
		chdataIn = <-InCH
		symbol := chdataIn.Ch
		OutBlockSize = chdataIn.MaxExpected

		bits := m.DeModulateBits(symbol)
		result := complex(float64(bits[0]), float64(bits[1]))

		chdataOut.Ch = result
		chdataOut.MaxExpected = OutBlockSize
		chdataOut.Message = chdataIn.Message
		// fmt.Printf("\n Writing demod bits to  %v ", outCH.Ch)
		outCH <- chdataOut
	}
	// close(InCH.Ch)
}

// func (m *Modem) DeModulate(symbolChannel gocomm.Complex128Channel, outCH gocomm.Complex128Channel) {
// 	var bits []uint8
// 	bits = make([]uint8, m.bitsPerSymbol)
// 	length := m.bitsPerSymbol
// 	N := m.bitsPerSymbol

// 	for i := 0; i < length; i++ {
// 		bits[i] = <-symbolChannel
// 	}
// 	key := toStr(bits[0:N])
// 	symbolChannel <- m.constellationTable[key]

// }

func (m *Modem) Modulate(bitchan gocomm.BitChannel, symbolChannel gocomm.Complex128Channel) {
	// fmt.Printf("\n Reading bits from  %v", bitchan.Ch)
	var chdataIn gocomm.SBitChannel
	var chdataOut gocomm.SComplex128Channel

	// fmt.Printf("\n MaxSymbols expected is %d , message = %v", bitchan.MaxExpected, bitchan.Message)

	var bits []uint8
	bits = make([]uint8, m.bitsPerSymbol)
	length := m.bitsPerSymbol
	N := m.bitsPerSymbol

	for i := 0; i < length; i++ {
		chdataIn = <-bitchan
		chdataOut.MaxExpected = chdataIn.MaxExpected / 2
		bits[i] = chdataIn.Ch
	}

	key := toStr(bits[0:N])
	chdataOut.Ch = m.constellationTable[key]
	symbolChannel <- chdataOut

	// fmt.Printf("\n Writing symbols to  %v", symbolChannel.Ch)

}

func (m *Modem) ModulateBits(bits []uint8) []complex128 {
	length := len(bits)
	slength := length / m.bitsPerSymbol
	var result = make([]complex128, slength)
	/// Actual Modulation happens here
	N := m.bitsPerSymbol

	cnt := 0
	for i := 0; i < length; i += N {
		key := toStr(bits[i : i+2])
		result[cnt] = m.constellationTable[key]
		cnt++

	}
	return result
}

// func (m *Modem) ModulateBitsCH(bits []uint8) []complex128 {
// 	length := len(bits)
// 	slength := length / m.bitsPerSymbol
// 	var result = make([]complex128, slength)
// 	/// Actual Modulation happens here
// 	N := m.bitsPerSymbol

// 	cnt := 0

// 	for i := 0; i < length; i += N {
// 		symch := make(chan [10]complex128, 1)
// 		go m.GenerateSymbolCH(i, bits[i:i+2], symch)

// 		for symbol := range symch {
// 			result[cnt] = symbol
// 			fmt.Printf("\n Received %d %v", cnt, symbol)
// 			cnt++
// 		}
// 	}

// 	return result
// }

func (m *Modem) GenerateSymbolCH(index int, bits []uint8, txch chan complex128) {
	key := toStr(bits)

	txsymbol := m.constellationTable[key]
	fmt.Printf("\n Begin Tx %d %v", index, txsymbol)
	// if index == 2 {
	// 	fmt.Printf("............Putting to Transmitter sleep")
	// 	time.Sleep(2 * time.Second)
	// 	fmt.Printf("Awake")
	// }
	fmt.Printf("\n Sent to Tx Channel %d %v", index, txsymbol)
	txch <- txsymbol
	fmt.Printf("\n Ends Tx Routine %d %v", index, txsymbol)
}
func (m *Modem) DemodSymbolCH(index int, rxsymbolchIn chan complex128, detBitsChOut chan []uint8) {

	rxsymbol := <-rxsymbolchIn
	fmt.Printf("\nReceived Symbols %d %v", index, rxsymbol)
	close(rxsymbolchIn)
	fmt.Printf("\nClosed Tx channel %d %v", index, rxsymbol)
	str := m.MinDistanceKey(rxsymbol)
	rxbits := make([]uint8, 2)

	N := m.bitsPerSymbol
	var num uint8
	for bit := 0; bit < N; bit++ {
		num = 0
		if str[bit] == '1' {
			num = 1
		}
		rxbits[bit] = num
	}
	fmt.Printf("\nWriting to Rx Channel..%d", index)
	// if index == 3 {
	// 	fmt.Printf("............Putting to Receiver %d sleep", index)
	// 	time.Sleep(2 * time.Second)

	// }
	detBitsChOut <- rxbits
	fmt.Printf("\n.. Closing Rx routine %d ", index)
	close(detBitsChOut)

}
func (m *Modem) Init(wordlength int) {
	m.bitsPerSymbol = wordlength
	switch wordlength {
	case 1:
		m.name = "BPSK"
	case 2:
		m.name = "QPSK"
		m.offset = math.Pi / 4.0
	case 3:
		m.name = "8PSK"
	case 4:
		m.name = "16PSK"

	case 256:
		m.name = "256QAM"

		n := int64(123)
		str := strconv.FormatInt(n, 2)
		fmt.Println("\n")
		var bitvec = make([]int, 8)
		for indx, _ := range str {
			bitvec[indx], _ = strconv.Atoi(string(str[indx]))

		}

		//Symbol=	{(1-2bi)[8-(1-2bi+2)[4-(1-2bi+4)[2-(1-2bi+6)]]]                +j(1-2bi+1)[8-(1-2bi+3)[4-(1-2bi+5)[2-(1-2bi+7)]]]}
		// realvalue:=(1-2*bitvec[0])*(8-(1-2*bitvec[2])*(4-(1-2*bitvec[4])*(2-(1-2*bitvec[6]))));
		// realvalue:=(1-2*bitvec[0])*(-8-(1-2*bitvec[2])*(4-(1-2*bitvec[4])*(2-(1-2*bitvec[6]))));
		realvalue := (1.0 - 2.0*bitvec[0]) * (8 - (1-2*bitvec[2])*(4-(1-2*bitvec[4])*(2-(1-2*bitvec[6]))))
		imagvalue := (1 - 2*bitvec[1]) * (8 - (1-2*bitvec[3])*(4-(1-2*bitvec[2]+5)*(2-(1-2*bitvec[7]))))
		fmt.Printf("%f+j%f", realvalue, imagvalue)
		// symbol := complex(float64(realvalue), float64(imagvalue))

	}
	// var i float64 = 0
	var length int = int(math.Exp2(float64(m.bitsPerSymbol)))
	m.size = length
	m.Constellation = make([]complex128, length)
	m.constellationTable = make(map[string]complex128, length)
	m.keys = make([]string, length)
	m.keys = sources.GrayCode(length)

	for i := 0; i < (len(m.Constellation)); i++ {
		var angle = float64(length-i) * 2 * math.Pi / float64(length)
		value := complex(math.Cos(angle+m.offset), math.Sin(angle+m.offset))
		m.Constellation[int(i)] = value
		// key := strconv.FormatInt(int64(i), 2)
		// if len(key) < m.bitsPerSymbol {
		// key = "0" + key
		// }
		key := m.keys[i]
		m.constellationTable[key] = value
		// fmt.Print("\n Init ", key, value, m.constellationTable[key])

		// m.keys[i] = key
	}

}

func (m *Modem) Print() {
	fmt.Printf("\n Map Table : %v", m.constellationTable)
}
func (m Modem) String() string {
	return fmt.Sprintf("A %s Modem with %d bits/Symbol", m.name, m.bitsPerSymbol)
}

func shift(data []complex128, offset complex128) []complex128 {
	var result = make([]complex128, len(data))
	for i := 0; i < len(data); i++ {
		result[i] = (0.707/2.0)*data[i] + offset
	}
	return result
}

// PinsIn() int
// 	PinsOut() int
// 	Pin(int) PinInfo
// PinsIn() int
// 	PinsOut() int
// 	Pin(int) PinInfo
// 	ModulesCount() int
// 	Module(int) ModuleInfo
func (m Modem) InPinCount() int {
	return 2
}

func (m Modem) OutPinCount() int {
	return 2
}
func (m Modem) Pin(pid int) chipset.PinInfo {
	// result := new(chipset.PinInfo)
	return m.Pins[m.PinNames[pid]]
	// return result
}

func (m Modem) PinIn(pid int) chipset.PinInfo {
	if pid >= m.InPinCount() {
		fmt.Printf("%d > No of Input Pins %d", pid, m.InPinCount())
		var result chipset.PinInfo
		result.Id = -1
		return result
	}

	return m.Pins[m.PinNames[pid]]

}
func (m Modem) PinByID(pid int) chipset.PinInfo {

	return m.Pins[m.PinNames[pid]]
}

func (m Modem) PinOut(pid int) chipset.PinInfo {
	if pid >= m.OutPinCount() {
		fmt.Printf("%d > No of Output Pins %d", pid, m.OutPinCount())
		var result chipset.PinInfo
		result.Id = -1
		return result
	}
	return m.Pins[m.PinNames[pid+m.InPinCount()]]

}

func (m Modem) PinByName(pinname string) chipset.PinInfo {
	return m.Pins[pinname]
}

// Has Modulator and Demodulator
func (m Modem) ModulesCount() int {
	return 2
}
func (m Modem) ModuleByName(mname string) chipset.ModuleInfo {
	return m.Modules[mname]
}

func (m Modem) Module(moduleid int) chipset.ModuleInfo {
	return m.ModuleByName(m.ModuleNames[moduleid])
	// fmt.Printf("\n Method Name %v", method.Type.String())
	// var err bool
	// method, _ := reflect.TypeOf(m).MethodByName("DeModulateBlock")
	// minfo.Function = method.Func
	// fmt.Printf("\n Module (%d) %v", moduleid, minfo.Function)
	// m.DeModulateBlock(OutBlockSize, InCH, outCH)
	// return minfo
}

func (m Modem) Name() string {
	return "MoDem"
}

func (m *Modem) SayModulate(dummy gocomm.BitChannel) {

	// fmt.Printf("\n BLAH BLAH BLAH BLAH BLAH BLAH %v MODULATE RUNNING...", m.Name())
	// maxInputExpected:=
	// fmt.Printf("\n AT MODEM MaxBits  expected is %d", dummy.MaxExpected)
	// temp := 1
	// cnt := 0
	// for i := 0; i < 10; i++ {
	// fmt.Printf("\n%v %d : %s", dummy, i, dummy.Message)
	m.ModulateBlock(1, dummy, m.Pins["symbolOut"].Channel.(gocomm.Complex128Channel))
	// }
	// close(dummy)
	// close()

}
func (m *Modem) SayDemodulate(dummy gocomm.Complex128Channel) {
	// fmt.Printf("\n CHIPSET %v Demodulator RUNNING...", m.Name())
	// fmt.Printf("\n Demodulator will try reading from %v  ", dummy.Ch)
	outCH := m.Pins["bitOut"].Channel.(gocomm.Complex128Channel)

	m.DeModulateBlock(1, dummy, outCH)

}

func (m Modem) IsInitialized() bool {
	return m.isInitialized
}

func (m *Modem) InitModules() {
	var totalModules int = m.ModulesCount()
	m.Modules = make(map[string]chipset.ModuleInfo, totalModules)
	m.ModuleNames = make(map[int]string, totalModules)
	// b := [...]string{"Penn", "Teller"}
	strlist := [...]string{"modulate", "demodulate"}
	for i := 0; i < totalModules; i++ {
		m.ModuleNames[i] = strlist[i]
	}

	for i := 0; i < totalModules; i++ {
		var minfo chipset.ModuleInfo
		minfo.Name = m.ModuleNames[i]

		switch minfo.Name {
		case "modulate":
			minfo.Desc = "This modulate modulates the input bits from Inpin 'bitIn' and writes to 'symbolOut' "
			minfo.Id = 0
			minfo.InPins = []int{m.PinByName("bitIn").Id}
			minfo.OutPins = []int{m.PinByName("symbolOut").Id}
			method := reflect.ValueOf(m).MethodByName("SayModulate")
			minfo.Function = method

		case "demodulate":
			minfo.Desc = "This modulate DeModulates the input Symbols from Inpin 'symbolIn' and writes to 'bitOut' "
			minfo.Id = 1
			minfo.InPins = []int{m.PinByName("symbolIn").Id}
			minfo.OutPins = []int{m.PinByName("bitOut").Id}
			method := reflect.ValueOf(m).MethodByName("SayDemodulate")
			minfo.Function = method

		}
		m.Modules[minfo.Name] = minfo
	}

}
func (m *Modem) InitPins() {
	m.isInitialized = true

	totalpins := m.InPinCount() + m.OutPinCount()
	m.Pins = make(map[string]chipset.PinInfo, totalpins)
	m.PinNames = make(map[int]string, totalpins)
	// b := [...]string{"Penn", "Teller"}
	strlist := [4]string{"bitIn", "symbolIn", "symbolOut", "bitOut"}
	for i := 0; i < totalpins; i++ {
		m.PinNames[i] = strlist[i]

	}

	for i := 0; i < totalpins; i++ {
		var pinfo chipset.PinInfo
		// pinfo.CreateComplex128Channel()
		pinfo.Name = m.PinNames[i]
		pinfo.Id = i
		if i < m.InPinCount() {
			pinfo.InputPin = true
		} else {
			pinfo.InputPin = false
		}
		m.Pins[m.PinNames[i]] = pinfo

	}

	testcch := gocomm.NewComplex128Channel()
	testch := gocomm.NewBitChannel()

	var dummypin chipset.PinInfo

	/// all Input Pins
	dummypin = m.Pins["symbolIn"]
	dummypin.Id = 1
	dummypin.DataType = reflect.TypeOf(testcch)
	m.Pins["symbolIn"] = dummypin

	dummypin = m.Pins["bitIn"]
	dummypin.Id = 0
	dummypin.DataType = reflect.TypeOf(testch)
	m.Pins["bitIn"] = dummypin

	///
	/// All output pins
	dummypin = m.Pins["symbolOut"]
	dummypin.Id = 2
	dummypin.DataType = reflect.TypeOf(testcch)
	dummypin.CreateComplex128Channel()
	dummypin.SourceName = "modulate"
	m.Pins["symbolOut"] = dummypin
	// fmt.Printf("\n %v : PinO : %v , Channel : %v", m.Name(), m.Pins["symbolOut"].Name, m.Pins["symbolOut"].Channel)
	dummypin = m.Pins["bitOut"]
	dummypin.Id = 3
	dummypin.SourceName = "demodulate"
	dummypin.DataType = reflect.TypeOf(testcch)
	dummypin.CreateComplex128Channel()
	m.Pins["bitOut"] = dummypin
	// fmt.Printf("\n %v : PinO : %v , Channel : %v", m.Name(), m.Pins["bitOut"].Name, m.Pins["bitOut"].Channel)

}

func (m *Modem) InitializeChip() {

	m.InitPins()
	m.InitModules()

}
