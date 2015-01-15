package main

import (
	"fmt"
	"github.com/wiless/gocomm/chipset"
	"github.com/wiless/gocomm/customchips"
)

func main() {

	var customchip customchips.QtChip
	customchip.InitializeChip()
	var wowchip chipset.Chip
	wowchip = customchip

	fmt.Printf("\n SampleChip = %#v", wowchip)
}
