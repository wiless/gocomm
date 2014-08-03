package sources

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"reflect"
	"strconv"
	"wiless/gocomm"
	"wiless/gocomm/chipset"
)

type Source struct {
	block int
}

func (s *Source) Size() int {
	return s.block
}

func (s *Source) SetSize(size int) {
	s.block = size
}

type BitSource struct {
	Source
	Pins          map[string]chipset.PinInfo
	PinNames      map[int]string
	Modules       map[string]chipset.ModuleInfo
	ModuleNames   map[int]string
	isInitialized bool
}

func (s *BitSource) GenBit(bitChannel gocomm.BitChannel) {
	// fmt.Printf("\n Writing bits to %v", bitChannel)
	bits := RandB(s.Size())
	// var data gocomm.SBitChannel
	var chdata gocomm.SBitChannel
	chdata.MaxExpected = s.Size()

	fmt.Println("\ntxbits=", bits)
	for i := 0; i < s.Size(); i++ {
		chdata.Ch = bits[i]
		// if i == (s.Size() - 1) {
		// 	chdata.MaxExpected += 6
		// 	fmt.Printf("\n extra being pushed, ...")
		// }
		bitChannel <- chdata
	}

	/// Pushing more extra 6 BITS

	// for i := 0; i < 6; i++ {
	// 	chdata.Ch = bits[i]
	// 	bitChannel <- chdata
	// }

	// fmt.Printf("\n SourceBlock closing..")
}

func Randsrc(size int, maxvalue int) []int {
	var result = make([]int, size)
	for i := 0; i < size; i++ {
		result[i] = rand.Intn(maxvalue)
	}
	return result
}

func RandB(size int) []uint8 {
	var result = make([]uint8, size)
	for i := 0; i < size; i++ {
		result[i] = uint8(rand.Intn(2))
	}
	return result
}

func RandReadableChars(size int) []uint8 {

	/// 32 to 126
	var result = make([]uint8, size)
	var startChar byte = 32
	for i := 0; i < size; i++ {
		result[i] = startChar + uint8(rand.Intn(94))
	}
	return result
}

func RandChars(size int) []uint8 {
	var result = make([]uint8, size)

	for i := 0; i < size; i++ {
		result[i] = uint8(rand.Intn(256))

	}
	return result
}

func (m BitSource) PinByName(pinname string) chipset.PinInfo {
	return m.Pins[pinname]
}
func BitsFromMessage(msg string) []uint8 {
	size := len(msg)
	var result = make([]uint8, size*8)

	cnt := 0
	for i := 0; i < size; i++ {
		var val int64 = int64(msg[i])
		bitstr := strconv.FormatInt(val, 2)
		msb := 8 - len(bitstr)

		for j := 0; j < msb; j++ {
			result[cnt] = 0
			cnt++
		}
		for _, val := range bitstr {

			result[cnt] = uint8(val) - '0'
			cnt++

		}
	}
	// fmt.Printf("\n Final Bit stream  : %v", result)
	return result

}

func GrayCode(cnt int) []string {

	KeyWidth := int(math.Log2(float64(cnt)))
	// cnt = cnt - 1

	// KeyWidtKeyWidthh := 3

	// fmt.Printf("\nKey %v \n", KeyWidth)
	i := big.NewInt(0)
	i2 := big.NewInt(0)
	one := big.NewInt(1)
	k := big.NewInt(0)

	keys := make([]string, cnt)
	var indx int = 0
	for j := 0; j < cnt/2; j++ {
		k.Xor(i, i2)
		keys[indx] = fmt.Sprintf("%0*b", KeyWidth, k)
		// fmt.Printf("LENG %s - %d : %v", keys[indx], len(keys[indx]), strings.TrimSpace(keys[indx]))

		indx++
		i.Add(i, one)

		k.Xor(i, i2)
		keys[indx] = fmt.Sprintf("%0*b", KeyWidth, k)
		// fmt.Printf("LENG %s %d %v", keys[indx], len(keys[indx]), strings.TrimSpace(keys[indx]))
		indx++
		// fmt.Printf("%0*b\n", KeyWidth, k)
		i.Add(i, one)
		i2.Add(i2, one)
	}

	return keys
}
func (b BitSource) name() {

}

/// CHIP related interafaces
func (b BitSource) Name() string {
	return "BitSource"
}
func (b BitSource) InPinCount() int {
	return 0
}

func (b BitSource) OutPinCount() int {
	return 1
}

func (b BitSource) PinIn(pid int) chipset.PinInfo {
	return b.Pins[b.PinNames[pid]]

}
func (b BitSource) PinOut(pid int) chipset.PinInfo {
	return b.Pins[b.PinNames[pid+b.InPinCount()]]
}
func (b BitSource) ModulesCount() int {
	return 1
}
func (b BitSource) Module(int) chipset.ModuleInfo {
	var dummy chipset.ModuleInfo
	return dummy
}
func (b BitSource) IsInitialized() bool {
	return b.isInitialized
}

func (b *BitSource) SayHello() string {
	return "\n Hi from " + b.Name()
}
func (m *BitSource) InitModules() {
	var totalModules int = m.ModulesCount()
	m.Modules = make(map[string]chipset.ModuleInfo, totalModules)
	m.ModuleNames = make(map[int]string, totalModules)
	// b := [...]string{"Penn", "Teller"}
	strlist := [...]string{"GenBit"}
	for i := 0; i < totalModules; i++ {
		m.ModuleNames[i] = strlist[i]
	}

	for i := 0; i < totalModules; i++ {
		var minfo chipset.ModuleInfo
		minfo.Name = m.ModuleNames[i]

		switch minfo.Name {
		case "GenBit":
			minfo.Desc = "Generates Uniformly distributed bits 0 and 1 at output pin 'bitOut' "
			minfo.InPins = []int{}
			minfo.OutPins = []int{0}
			method := reflect.ValueOf(m).MethodByName("SayHello")
			minfo.Function = method
			// case "demodulate":
			// 	minfo.InPins = []int{1}
			// 	minfo.OutPins = []int{1}

		}
		m.Modules[minfo.Name] = minfo
	}

}
func (m *BitSource) InitPins() {
	m.isInitialized = true

	totalpins := m.InPinCount() + m.OutPinCount()
	m.Pins = make(map[string]chipset.PinInfo, totalpins)
	m.PinNames = make(map[int]string, totalpins)
	// b := [...]string{"Penn", "Teller"}
	strlist := [1]string{"bitOut"}
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

	testch := gocomm.NewBitChannel()
	var dummypin chipset.PinInfo
	/// All output pins
	dummypin = m.Pins["bitOut"]
	dummypin.Desc = "Output Pin where bits are written"
	dummypin.DataType = reflect.TypeOf(testch)
	dummypin.CreateBitChannel()
	m.Pins["bitOut"] = dummypin

}

func (m *BitSource) InitializeChip() {

	m.InitPins()
	m.InitModules()
	m.isInitialized = true
}
