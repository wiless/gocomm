package main

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	// "strings"
	// "text/template"
	"time"
	"wiless/gocomm/channel"
	"wiless/gocomm/chipset"
	"wiless/gocomm/modem"
	// "wiless/gocomm/sink"
	"wiless/gocomm/sources"
)

func CreateChipFrom(chipname string, nameit string) chipset.Chip {
	var result chipset.Chip
	switch chipname {
	case "Modem":
		fmt.Printf("\n Found Modem")
		var tmp modem.Modem
		tmp.SetName(nameit)
		/// Initialize the Chips
		///
		result = tmp

		break
	case "ChannelEmulator":
		var tmp channel.ChannelEmulator
		/// Initialize the Chips
		///

		result = tmp
		break
		fmt.Printf("\n Found Cool")
	case "BitSource":
		var tmp sources.BitSource
		/// Initialize the Chips
		///
		result = tmp
		break
		fmt.Printf("\n Found Bitsource")
	case "":
	default:
		fmt.Printf("\n Unknown chip set Type %v", chipname)
	}
	return result
}

type JsonCircuit struct {
	Name  string
	Chips []ChipDesc
	Links []LinkDesc
}

type ChipDesc struct {
	Type string
	Name string
}
type LinkDesc struct {
	Name                  string
	SourceChipName        string
	SourcePinName         string
	DestinationChipName   string
	DestinationModuleName string
	DestinationPinName    string
}

func main() {

	t := time.Now()
	var chips [4]chipset.Chip
	var circuit JsonCircuit
	// circuit.Chips = make([]ChipDesc, 4)
	// databytes, err := json.MarshalIndent(circuit, "", "\t")
	bytearray, ferr := ioutil.ReadFile("connection.json")
	fmt.Printf("\n Reading File %v, \n Error : %v \n ", "connection.json", ferr)
	err := json.Unmarshal(bytearray, &circuit)
	if err != nil {
		fmt.Print(err)
		fmt.Fprintf(os.Stdout, "%s", string(bytearray))
	}
	fmt.Printf("\n %#v", circuit)
	chips[0] = CreateChipFrom("Modem", "txmodem")
	chips[1] = CreateChipFrom("Modem", "rxmodem")
	chips[2] = CreateChipFrom("ChannelEmulator", "fading")
	chips[3] = CreateChipFrom("BitSource", "source")

	for i := 0; i < 4; i++ {
		fmt.Printf("\n CHIP %v", chips[i])
		for k := 0; k < chips[i].InPinCount()+chips[i].OutPinCount(); k++ {
			fmt.Printf("\n \t %d  PINS  : %v", k, chips[i].PinByID(k))
		}

	}

	fmt.Println("\n", time.Now())
	fmt.Printf("\nTime Elaspsed %v \n", time.Since(t))
}
