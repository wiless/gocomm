package main

import (
	"fmt"
	"reflect"
	"wiless/gocomm/modem"
	"wiless/vlib"
)

func main() {
	var gvec1 vlib.GIntVector
	var gvec2 vlib.GDoubleVector
	//gvec1 = make(vlib.GIntVector, 10)
	// var gvec2 vlib.GDoubleVector
	// var anyVector vlib.VectorIface

	gvec1.SetSize(10)
	gvec2.SetSize(13)
	// gvec1 = make([]int, 100)

	// anyVector{gvec1}.

	fmt.Printf("\n gvec1 = %v", gvec1)
	fmt.Printf("\n gvec2 = %v", gvec2)
	fmt.Printf("\n gvec2 = %v", vlib.GetSize(gvec2))
	fmt.Printf("\n gvec2 = %v", vlib.GetSize(gvec1))

	var v modem.Modem
	// var fm modem.ModemIterface
	// fm = v
	cntr := reflect.TypeOf(v).NumMethod()
	fmt.Printf("\n No. of methods %v", cntr)
	for i := 0; i < cntr; i++ {
		fmt.Printf("\n No. of methods %v", reflect.TypeOf(v).Method(i).Name)
	}
	fmt.Printf("\n =============================================")
	cntr = reflect.TypeOf(&v).NumMethod()
	fmt.Printf("\n No. of methods %v", cntr)
	for i := 0; i < cntr; i++ {
		fmt.Printf("\n No. of methods %v", reflect.TypeOf(&v).Method(i).Name)
	}
	// fmt.Printf("Hello %v", gvec2)
	//gvec1.Size()
	var nano car = car{5}
	var donkey animal
	var what somethingMakesCommon
	what = &donkey
	fmt.Printf("\n==============================\n")
	fmt.Printf("\nAnimal %v", donkey)
	fmt.Printf("\nCar %v", nano)
	fmt.Printf("\nSomeInterface %v", what)
	fmt.Printf("\n==============================\n")
	what.Bark(2)
	cntr = reflect.TypeOf(what).NumMethod()
	fmt.Printf("\n No. of methods %v", cntr)
	obj := reflect.ValueOf(what)

	for i := 0; i < cntr; i++ {
		methodstr := reflect.TypeOf(what).Method(i).Name
		nargs := reflect.TypeOf(what).Method(i).Type.NumIn()
		fmt.Printf("\n Name of the method %v", methodstr)

		fmt.Printf("\n No. of Arguments of the method %v is  %#v", methodstr, nargs-1)
		fmt.Printf("\n No. of Arguments of the method %v is  %v", methodstr, reflect.TypeOf(what).Method(i))
		//fmt.Printf("\nIs valid to call %v", reflect.ValueOf(what).Method(i).IsValid())
		fmt.Printf(" ... Szzz trying to call ")
		//	in := make([]reflect.Value, nargs)
		nargs = obj.Method(i).Type().NumIn()
		args := make([]reflect.Value, nargs)
		// args[0].SetInt(5)
		for j := 0; j < nargs; j++ {
			var para int = 5
			args[j] = reflect.ValueOf(para)

			// fmt.Printf("\n ARG %v", reflect.TypeOf(what).Method(i).Type.In(j))
			// // reflect.TypeOf(what).Method(i).Type.In(j) = 5
			// args[j] = reflect.ValueOf(what).Method(i)
		}
		outvals := reflect.ValueOf(what).Method(i).Call(args)
		fmt.Printf(".. Worked ?? Output of %s", methodstr)
		for indx, val := range outvals {
			fmt.Printf("\n OUTPUT %d : %v", indx, val.Interface())
		}

	}
	// fmt.Printf(
}

type somethingMakesCommon interface {
	Bark(int) []string
}

type car struct {
	wheels int
}

func (c car) String() string {
	return fmt.Sprintf("I can Drive on RoAD %d tyres", c.wheels)
}

func (c *car) Bark(times int) []string {
	for i := 0; i < times; i++ {
		fmt.Printf("\nBeep..")
	}
	return []string{"Tata", "Motors"}
}

func (c car) Sleep() {
	fmt.Printf("----I can sleep by turn off---")
}
func (a animal) Sleep() {
	fmt.Printf("---I can sleep---")
}

func (a *animal) Bark(times int) []string {
	for i := 0; i < times; i++ {
		fmt.Printf("\n Yonkee...")
	}
	return []string{"meow"}
}

type animal struct {
	legs int
	tail float32
}
