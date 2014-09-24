package main

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"
	"wiless/gocomm/chipset"
	"wiless/gocomm/sources"
	"wiless/gocomm/sink"

	

)

type JsonCircuit struct {
	Name    string
	ChipsNames []string

	// Pins    []JsonPin
	// Modules []JsonModule

	// InPinCount    int
	// OutPinCount   int
	// ModuleCount   int
	// PinCount      int
	// PinNames      string
	// ModuleNames   string
	// WriteTestMain bool
}

chipset.Chip CreateChipFrom(chipname string){
	var result chipset.Chip
	switch chipname{
		"modem":

	}
}

func main() {

// }st
