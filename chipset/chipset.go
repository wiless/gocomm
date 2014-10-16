package chipset

import (
	"fmt"
	"reflect"
	"wiless/gocomm"
)

type PinInfo struct {
	Id          int
	Name        string
	Desc        string
	DataType    reflect.Type
	InputPin    bool
	NonBlocking bool
	Channel     interface{}
	SourceName  string
}

type ModuleInfo struct {
	Id           int
	Name         string
	Desc         string
	InPins       []int
	OutPins      []int
	FunctionName string
	Function     reflect.Value
}

type Chip interface {
	Name() string
	InPinCount() int
	OutPinCount() int
	PinByName(string) PinInfo
	PinByID(int) PinInfo
	Set(json string)    /// Json formatted parameters for the Chip
	Get() string        /// returns the JSON formatted parameters of the Chip
	Commands() []string /// returns the commands recognized by the Chip
	ModuleByName(string) ModuleInfo

	PinIn(indx int) PinInfo
	PinOut(indx int) PinInfo
	ModulesCount() int
	Module(int) ModuleInfo
	IsInitialized() bool
}

func (p PinInfo) String() string {
	if p.Channel != nil {
		return fmt.Sprintf("Name:%s,Input=%v,Type %v,Desc %v, Channel Type : %v ", p.Name, p.InputPin, p.DataType.Name(), p.Desc, reflect.TypeOf(p.Channel).Name())
	} else {
		return fmt.Sprintf("Name:%s,Input=%v,Type %v,Desc %v, Channel : NIL ", p.Name, p.InputPin, p.DataType.Name(), p.Desc)

	}
}

func (p *PinInfo) CreateChannel() {
	switch p.DataType.Name() {
	case "BitChannel":
		p.Channel = gocomm.NewBitChannel()
	case "Complex128Channel":
		p.Channel = gocomm.NewComplex128Channel()
	case "Complex128AChannel":
		p.Channel = gocomm.NewComplex128AChannel()
	case "BitAChannel":
		p.Channel = gocomm.NewBitAChannel()
	case "FloatAChannel":
		p.Channel = gocomm.NewFloatAChannel()
	case "FloatChannel":
		p.Channel = gocomm.NewFloatChannel()

	default:
		fmt.Printf("PinInfo::CreateChannel():Unkown Type %v needs channel", p.DataType)
		p.Channel = reflect.New(p.DataType)
		fmt.Printf("\n atteptemd Channel =  %v ", p.Channel)
	}

}

func (p *PinInfo) CreateBitChannel() {
	p.Channel = gocomm.NewBitChannel()
}

func (p *PinInfo) CreateBitAChannel() {
	p.Channel = gocomm.NewBitAChannel()
}

func (p *PinInfo) CreateComplex128Channel() {
	p.Channel = gocomm.NewComplex128Channel()

}
func (p *PinInfo) CreateComplex128AChannel() {
	p.Channel = gocomm.NewComplex128AChannel()
}

func Sink(pin PinInfo) {

	fmt.Printf("\n======Sink DataOut from Pin %v =========== \n\n", pin)
	count := 1
	switch pin.DataType.Name() {
	case "FloatChannelA":
		for i := 0; i < count; i++ {
			// fmt.Printf("\n Status of Channel %d = %#v ", i, pin.Channel)
			ddata := <-ToFloatACH(pin)
			// fmt.Printf(" SPECIAL MESSAGE %s", ddata.Message)
			if ddata.Message == "" {
				fmt.Printf("\nSinkPin : %s - Read Bit %d = %v ", pin.Name, i, ddata.Ch)
			} else {
				fmt.Printf("\nSinkPin : %s - Read Bit %d = %v : %s", pin.Name, i, ddata.Ch, ddata.Message)
			}

			count = ddata.MaxExpected
			// ddata := choutData.Ch
			// max = choutData.MaxExpected
			// fmt.Printf(" %d %d", uint8(real(ddata)), uint8(imag(ddata)))
			// fmt.Printf("\n %d @ max Symbols limit  = %d %s ", i, max, choutData.Message)

		}
	case "FloatChannel":
		for i := 0; i < count; i++ {
			// fmt.Printf("\n Status of Channel %d = %#v ", i, pin.Channel)
			ddata := <-pin.Channel.(gocomm.FloatChannel)
			// fmt.Printf(" SPECIAL MESSAGE %s", ddata.Message)
			if ddata.Message == "" {
				fmt.Printf("\nSinkPin : %s - Read Bit %d = %v ", pin.Name, i, ddata.Ch)
			} else {
				fmt.Printf("\nSinkPin : %s - Read Bit %d = %v : %s", pin.Name, i, ddata.Ch, ddata.Message)
			}

			count = ddata.MaxExpected
			// ddata := choutData.Ch
			// max = choutData.MaxExpected
			// fmt.Printf(" %d %d", uint8(real(ddata)), uint8(imag(ddata)))
			// fmt.Printf("\n %d @ max Symbols limit  = %d %s ", i, max, choutData.Message)

		}
	case "BitChannelA":
		for i := 0; i < count; i++ {
			// fmt.Printf("\n Status of Channel %d = %#v ", i, pin.Channel)
			ddata := <-ToBitACH(pin) //pin.Channel.(gocomm.BitChannelA)
			// fmt.Printf(" SPECIAL MESSAGE %s", ddata.Message)
			if ddata.Message == "" {
				fmt.Printf("\nSinkPin : %s - Read Bit %d = %v ", pin.Name, i, ddata.Ch)
			} else {
				fmt.Printf("\nSinkPin : %s - Read Bit %d = %v : %s", pin.Name, i, ddata.Ch, ddata.Message)
			}

			count = ddata.MaxExpected
			// ddata := choutData.Ch
			// max = choutData.MaxExpected
			// fmt.Printf(" %d %d", uint8(real(ddata)), uint8(imag(ddata)))
			// fmt.Printf("\n %d @ max Symbols limit  = %d %s ", i, max, choutData.Message)

		}
	case "BitChannel":
		for i := 0; i < count; i++ {
			// fmt.Printf("\n Status of Channel %d = %#v ", i, pin.Channel)
			ddata := <-pin.Channel.(gocomm.BitChannel)
			// fmt.Printf(" SPECIAL MESSAGE %s", ddata.Message)
			if ddata.Message == "" {
				fmt.Printf("\nSinkPin : %s - Read Bit %d = %v ", pin.Name, i, ddata.Ch)
			} else {
				fmt.Printf("\nSinkPin : %s - Read Bit %d = %v : %s", pin.Name, i, ddata.Ch, ddata.Message)
			}

			count = ddata.MaxExpected
			// ddata := choutData.Ch
			// max = choutData.MaxExpected
			// fmt.Printf(" %d %d", uint8(real(ddata)), uint8(imag(ddata)))
			// fmt.Printf("\n %d @ max Symbols limit  = %d %s ", i, max, choutData.Message)

		}
	case "Complex128Channel":
		for i := 0; i < count; i++ {
			ddata := <-pin.Channel.(gocomm.Complex128Channel)
			count = ddata.MaxExpected
			// fmt.Printf(" SPECIAL MESSAGE %s", ddata.Message)
			if ddata.Message == "" {
				fmt.Printf("\nPin : %s - Read Complex (%d of %d)  = %v ", pin.Name, i, count, ddata.Ch)
			} else {
				fmt.Printf("\nPin : %s - Read Complex (%d of %d) = %v : %s", pin.Name, i, count, ddata.Ch, ddata.Message)
			}

			// ddata := choutData.Ch
			// max = choutData.MaxExpected
			// fmt.Printf(" %d %d", uint8(real(ddata)), uint8(imag(ddata)))
			// fmt.Printf("\n %d @ max Symbols limit  = %d %s ", i, max, choutData.Message)

		}
	case "Complex128ChannelA":
		for i := 0; i < count; i++ {
			ddata := <-pin.Channel.(gocomm.Complex128AChannel)
			count = ddata.MaxExpected
			// fmt.Printf(" SPECIAL MESSAGE %s", ddata.Message)
			if ddata.Message == "" {
				fmt.Printf("\nPin : %s - Read Complex (%d of %d)  = %v ", pin.Name, i, count, ddata.Ch)
			} else {
				fmt.Printf("\nPin : %s - Read Complex (%d of %d) = %v : %s", pin.Name, i, count, ddata.Ch, ddata.Message)
			}

			// ddata := choutData.Ch
			// max = choutData.MaxExpected
			// fmt.Printf(" %d %d", uint8(real(ddata)), uint8(imag(ddata)))
			// fmt.Printf("\n %d @ max Symbols limit  = %d %s ", i, max, choutData.Message)

		}
	default:
		fmt.Printf("\n Unknown Data type")
	}

}

func ToComplexCH(pin PinInfo) gocomm.Complex128Channel {
	return pin.Channel.(gocomm.Complex128Channel)
}
func ToComplexACH(pin PinInfo) gocomm.Complex128AChannel {

	// if pin == 0 {
	// 	fmt.Printf("gocomm:ToComplexACH PinInfo is NULL")
	// 	return make(gocomm.Complex128AChannel)
	// }
	return pin.Channel.(gocomm.Complex128AChannel)
}
func ToBitCH(pin PinInfo) gocomm.BitChannel {
	return pin.Channel.(gocomm.BitChannel)
}
func ToBitACH(pin PinInfo) gocomm.BitChannel {
	return pin.Channel.(gocomm.BitChannel)
}
func ToFloatCH(pin PinInfo) gocomm.FloatChannel {
	return pin.Channel.(gocomm.FloatChannel)
}
func ToFloatACH(pin PinInfo) gocomm.FloatAChannel {
	return pin.Channel.(gocomm.FloatAChannel)
}
