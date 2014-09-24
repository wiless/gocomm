package customchips
/// Function 
func (m *OFDM) fft(inTxsymbols gocomm.Complex128Channel,){ 
/// Read your data from Input channel(s) 
/// And write it to OutputChannels defined below
	outTxsymbols:=m.PinByName("outTxsymbols").Channel(gocomm.gocomm.Complex128Channel),

}
package customchips
/// Function 
func (m *OFDM) ifft(inRxsymbols gocomm.Complex128Channel,){ 
/// Read your data from Input channel(s) 
/// And write it to OutputChannels defined below
	outRxsymbols:=m.PinByName("outRxsymbols").Channel(gocomm.gocomm.Complex128Channel),

}
