package chipset

import (
	"fmt"
	"reflect"
	"wiless/gocomm"
)

type PinInfo struct {
	Id       int
	Name     string
	Desc     string
	DataType reflect.Type
	InputPin bool
	Channel  interface{}
}

type ModuleInfo struct {
	Id       int
	Name     string
	Desc     string
	InPins   []int
	OutPins  []int
	Function reflect.Value
}

type Chip interface {
	Name() string
	InPinCount() int
	OutPinCount() int
	PinByName(string) PinInfo
	PinByID(int) PinInfo

	ModuleByName(string) ModuleInfo

	PinIn(int) PinInfo
	PinOut(int) PinInfo
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
	p.Channel = reflect.New(p.DataType)
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
