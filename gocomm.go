// Package gocomm provides different channel objects and functionalites related to them
package gocomm

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/wiless/vlib"
)

// Possibly to be deprecated
var WGroup sync.WaitGroup

type ChannelDataStruct interface {
	GetMaxExpected() int
}

type SBitObj struct {
	Ch          uint8 // The uint8 data is here
	MaxExpected int
	Message     string
	Ts          float64
	TimeStamp   float64
}

type SBitAObj struct {
	Ch          []uint8
	MaxExpected int
	Message     string
	Ts          float64
	TimeStamp   float64
}

type SFloatObj struct {
	Ch          float64
	MaxExpected int
	Message     string
	Ts          float64
	TimeStamp   float64
}

type SFloatAObj struct {
	Ch          []float64
	MaxExpected int
	Message     string
	Ts          float64
	TimeStamp   float64
}

type SComplex128Obj struct {
	Ch          complex128
	MaxExpected int
	Message     string
	Ts          float64
	TimeStamp   float64
}

type SComplex128AObj struct {
	Ch          []complex128
	MaxExpected int
	Message     string
	Ts          float64
	TimeStamp   float64
}

func (s *SBitObj) Next(sample uint8) {
	s.Ch = sample
	s.UpdateTimeStamp()
}

// Next function sets the sample given in argument and updates the Timestamp
func (s *SBitAObj) Next(sample []uint8) {
	s.Ch = sample
	s.UpdateTimeStamp()
}

func (s *SComplex128Obj) Next(sample complex128) {
	s.Ch = sample
	s.UpdateTimeStamp()
}

func (s *SComplex128AObj) Next(sample []complex128) {
	s.Ch = sample
	s.UpdateTimeStamp()
}

func (s *SBitObj) UpdateTimeStamp() {
	s.TimeStamp += s.Ts
}

func (s *SBitAObj) UpdateTimeStamp() {
	s.TimeStamp += s.Ts
}

func (s *SComplex128Obj) UpdateTimeStamp() {
	s.TimeStamp += s.Ts
}
func (s *SComplex128AObj) UpdateTimeStamp() {
	s.TimeStamp += s.Ts
}
func (s SBitObj) GetMaxExpected() int {
	return s.MaxExpected
}

func (s SBitAObj) GetMaxExpected() int {
	return s.MaxExpected
}

func (s SComplex128Obj) GetMaxExpected() int {
	return s.MaxExpected
}

func (s SComplex128AObj) GetMaxExpected() int {
	return s.MaxExpected
}
func (s SFloatAObj) GetMaxExpected() int {
	return s.MaxExpected
}

func (s SFloatObj) GetMaxExpected() int {
	return s.MaxExpected
}

type BitChannel chan SBitObj

type BitAChannel chan SBitAObj
type FloatChannel chan SFloatObj
type FloatAChannel chan SFloatAObj

type Complex128Channel chan SComplex128Obj

type Complex128AChannel chan SComplex128AObj

func NewBitChannel() BitChannel {
	return make(BitChannel, 1)
}

func NewBitAChannel() BitAChannel {
	return make(BitAChannel, 1)
}
func NewFloatChannel() FloatChannel {
	return make(FloatChannel, 1)
}

func NewFloatAChannel() FloatAChannel {
	return make(FloatAChannel, 1)
}

func NewComplex128Channel() Complex128Channel {
	return make(Complex128Channel, 1)
}

func NewComplex128AChannel() Complex128AChannel {
	return make(Complex128AChannel, 1)
}

func ChannelDuplexer(InCH Complex128Channel, OutCHA []Complex128Channel) {
	Nchanels := len(OutCHA)
	var chdataIn SComplex128Obj
	var chdataOut SComplex128Obj
	NextSize := 1
	for cnt := 0; cnt < NextSize; cnt++ {
		chdataIn = <-InCH
		data := chdataIn.Ch
		NextSize = chdataIn.MaxExpected

		// fmt.Printf("%d InputDuplexer : %v ", cnt, data)
		for i := 0; i < Nchanels; i++ {
			chdataOut.Ch = data
			chdataOut.MaxExpected = NextSize
			chdataOut.Message = chdataIn.Message
			OutCHA[i] <- chdataOut
		}
	}
	close(InCH)
}

func Bit2BitA(InCH BitChannel, OutCH BitAChannel, veclen int) {

	NextSize := 1
	var Outobjs SBitAObj
	var InObj SBitObj
	for j := 0; j < NextSize; j++ {

		vecdata := make([]uint8, veclen)

		for i := 0; i < veclen; i++ {
			InObj = <-InCH
			vecdata[i] = InObj.Ch
			NextSize = InObj.GetMaxExpected()
			if i == 0 {
				Outobjs.TimeStamp = InObj.TimeStamp
				Outobjs.Ts = InObj.Ts * float64(veclen)
				maxexp := int(float64(InObj.MaxExpected) / float64(veclen))
				Outobjs.MaxExpected = maxexp
				Outobjs.Message = InObj.Message + " S2P"
			}

		}
		Outobjs.Ch = vecdata
		// fmt.Printf("\n Wrote Vector %d of %d", j, NextSize)
		OutCH <- Outobjs
	}
}

/// Converts each Vector Sample to a Sample which can be processed at sample rate
/// This can be considered as Upsampler Each vector at rate Ts , is communicated to the next block at Ts/N samples
func Vector2Sample(uid int, NextSize int, InCH Complex128AChannel, OutCH Complex128Channel) {
	var chdataOut SComplex128Obj
	var chdataIn SComplex128AObj

	cnt := 0

	for i := 0; i < NextSize; i++ {
		chdataIn = <-InCH
		indata := chdataIn.Ch
		veclen := len(indata)
		cnt += veclen

		for indx := 0; indx < veclen; indx++ {
			chdataOut.Ch = indata[indx]
			OutCH <- chdataOut
		}
	}
	fmt.Printf("\n User%d : Closing", uid)

	// close(InCH)
}

func Complex2ComplexA(InCH Complex128Channel, OutCH Complex128AChannel, veclen int) {

	NextSize := 1
	var Outobjs SComplex128AObj
	var InObj SComplex128Obj
	for j := 0; j < NextSize; j++ {

		vecdata := make([]complex128, veclen)

		for i := 0; i < veclen; i++ {
			InObj = <-InCH
			vecdata[i] = InObj.Ch
			NextSize = InObj.GetMaxExpected()
			if i == 0 {
				Outobjs.TimeStamp = InObj.TimeStamp
				Outobjs.Ts = InObj.Ts * float64(veclen)
				maxexp := int(float64(InObj.MaxExpected) / float64(veclen))
				Outobjs.MaxExpected = maxexp
				Outobjs.Message = InObj.Message + " S2P"
			}

		}
		Outobjs.Ch = vecdata
		// fmt.Printf("\n Wrote Vector %d of %d", j, NextSize)
		OutCH <- Outobjs
	}

}
func ComplexA2Complex(InCH Complex128AChannel, OutCH Complex128Channel) {
	var chdataOut SComplex128Obj
	var chdataIn SComplex128AObj

	NextSize := 1 // NoOfSuch vectors expected
	for i := 0; i < NextSize; i++ {
		chdataIn = <-InCH
		NextSize = chdataIn.MaxExpected

		indata := chdataIn.Ch
		veclen := len(indata)

		chdataOut.Message = chdataIn.Message + " P2S"
		chdataOut.MaxExpected = veclen * NextSize
		chdataOut.TimeStamp = chdataIn.TimeStamp
		chdataOut.Ts = chdataIn.Ts / float64(veclen)
		// fmt.Printf("\n Received %d of %d block,blocksize = %d, Output Expected %d", i+1, NextSize, veclen, chdataOut.MaxExpected)
		// fmt.Printf("\n Data TimeStamp : %v", chdataIn.TimeStamp)
		for indx := 0; indx < veclen; indx++ {
			chdataOut.Ch = indata[indx]
			OutCH <- chdataOut
			chdataOut.TimeStamp += chdataOut.Ts

		}
	}
	// close(InCH)
}

func (s SComplex128Obj) String() string {
	result := fmt.Sprintf("(%f,%f)\t%% [%d]@%f : %s", s.TimeStamp, s.Ch, s.MaxExpected, s.Ts, s.Message)
	return result
}
func (s SComplex128AObj) String() (result string) {
	result = fmt.Sprintf("(%f,%f)\t%% [%d]@%f : %s\n", s.TimeStamp, s.Ch, s.MaxExpected, s.Ts, s.Message)
	return result
}
func (s SBitObj) String() string {
	result := fmt.Sprintf("(%f,%d)\t%% [%d]@%f : %s", s.TimeStamp, s.Ch, s.MaxExpected, s.Ts, s.Message)
	return result
}
func (s SBitAObj) String() (result string) {
	result = fmt.Sprintf("(%f,%d)\t%% [%d]@%f : %s\n", s.TimeStamp, s.Ch, s.MaxExpected, s.Ts, s.Message)
	return result
}

func ComplexA2ComplexFn(InObj SComplex128AObj) (Outobjs []SComplex128Obj) {

	veclen := len(InObj.Ch)
	Outobjs = make([]SComplex128Obj, veclen)

	maxexp := InObj.MaxExpected * veclen
	Ts := InObj.Ts / float64(veclen)
	for i := 0; i < veclen; i++ {
		Outobjs[i].Ch = InObj.Ch[i]
		Outobjs[i].Message = InObj.Message
		Outobjs[i].MaxExpected = maxexp
		Outobjs[i].Ts = Ts
		Outobjs[i].TimeStamp = InObj.TimeStamp + float64(i)*Ts
	}
	return Outobjs
}

func Complex2ComplexAFn(InObj []SComplex128Obj) (Outobjs SComplex128AObj) {

	veclen := len(InObj)
	samples := vlib.NewVectorC(veclen)

	for i := 0; i < veclen; i++ {
		samples[i] = InObj[i].Ch
	}
	maxexp := int(float64(InObj[0].MaxExpected) / float64(veclen))
	Outobjs.Ch = samples
	Outobjs.Message = InObj[0].Message
	Outobjs.MaxExpected = maxexp
	Outobjs.Ts = InObj[0].Ts * float64(veclen)
	Outobjs.TimeStamp = InObj[0].TimeStamp

	return Outobjs
}

func SinkComplex(ch Complex128Channel, mode string) {
	cnt := 1
	if mode != "" {
		fmt.Printf("\n%s=[", mode)
		for i := 0; i < cnt; i++ {
			result := <-ch
			cnt = result.MaxExpected
			fmt.Printf("(%v,%v);", result.TimeStamp, result.Ch)
		}
		fmt.Printf("]")
	} else {
		for i := 0; i < cnt; i++ {
			result := <-ch
			cnt = result.MaxExpected
			fmt.Printf("\n Sink %d (of %d): %v", i+1, cnt, result)

		}
	}

	WGroup.Done()
}

func SinkComplexA(ch Complex128AChannel) {
	cnt := 1
	for i := 0; i < cnt; i++ {
		result := <-ch
		cnt = result.MaxExpected
		fmt.Printf("\n Sink %d (of %d): %v", i+1, cnt, result)

	}
	WGroup.Done()
}

func ComplexA2Bits(obj []SComplex128Obj) vlib.VectorB {
	count := len(obj)
	result := vlib.NewVectorB(count * 2)
	for i := 0; i < count; i++ {
		result[2*i] = uint8(real(obj[i].Ch))
		result[2*i+1] = uint8(imag(obj[i].Ch))
	}
	return result
}

func Complex2Bits(obj SComplex128Obj) (outobj []SBitObj) {

	result := make([]SBitObj, 2)

	for i := 0; i < 2; i++ {

		result[i].MaxExpected = obj.MaxExpected * 2
		result[i].Ts = obj.Ts / 2.0
		result[i].TimeStamp = obj.TimeStamp + float64(i)*result[i].Ts
		result[i].Message = obj.Message
	}
	result[0].Ch = uint8(real(obj.Ch))
	result[1].Ch = uint8(imag(obj.Ch))
	// fmt.Printf("\n bit encoded received %#v", result)
	// }
	return result
}

type generic interface{}

func ToInt(value generic) (result int) {
	switch reflect.TypeOf(value).String() {
	case "float64":
		return int(value.(float64))
	case "int":
		return value.(int)
	default:
		return 0
	}

	fmt.Print("Type is ", reflect.TypeOf(value))
	sol := int(value.(float64))
	defer func() {
		result = 0
	}()

	fmt.Print("Solution is ", sol)
	return result
}
