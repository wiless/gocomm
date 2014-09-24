package main

import (
	"fmt"
	"wiless/gocomm/chipset"
	"wiless/gocomm/customchips"
)

func main() {

	var customchip customchips.QtChip
	customchip.InitializeChip()
	var wowchip chipset.Chip
	wowchip = customchip

	fmt.Printf("\n SampleChip = %#v", wowchip)
}
