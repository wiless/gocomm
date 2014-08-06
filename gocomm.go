package gocomm

type ChannelDataStruct interface {
	GetMaxExpected() int
}

type SBitChannel struct {
	Ch          uint8
	MaxExpected int
	Message     string
}

func (s SBitChannel) GetMaxExpected() int {
	return s.MaxExpected
}

func (s SBitChannelA) GetMaxExpected() int {
	return s.MaxExpected
}

func (s SComplex128Channel) GetMaxExpected() int {
	return s.MaxExpected
}

func (s SComplex128ChannelA) GetMaxExpected() int {
	return s.MaxExpected
}

type SBitChannelA struct {
	Ch          []uint8
	MaxExpected int
	Message     string
}

type SComplex128Channel struct {
	Ch          complex128
	MaxExpected int
	Message     string
}

type SComplex128ChannelA struct {
	Ch          []complex128
	MaxExpected int
	Message     string
}

type BitChannel chan SBitChannel

type BitChannelA chan SBitChannelA

type Complex128Channel chan SComplex128Channel

type Complex128ChannelA chan SComplex128ChannelA

func NewBitChannel() BitChannel {
	return make(BitChannel, 1)
}

func NewBitAChannel() BitChannelA {
	return make(BitChannelA, 1)
}

func NewComplex128Channel() Complex128Channel {
	return make(Complex128Channel, 1)
}

func NewComplex128AChannel() Complex128ChannelA {
	return make(Complex128ChannelA, 1)
}
