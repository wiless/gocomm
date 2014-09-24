/// Function 
func (m *TurboCodec) encode(uncodedBits gocomm.BitChannel,){ 
/// Read your data from Input channel(s) 
/// And write it to OutputChannels defined below
	codedBits:=m.PinByName("codedBits").Channel(gocomm.gocomm.Complex128Channel),

}
/// Function 
func (m *TurboCodec) decode(codedBits gocomm.Complex128Channel,){ 
/// Read your data from Input channel(s) 
/// And write it to OutputChannels defined below
	decodedBits:=m.PinByName("decodedBits").Channel(gocomm.gocomm.Complex128Channel),

}
