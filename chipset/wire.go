package chipset

import (
	"fmt"

	"reflect"
	// "wiless/gocomm"
)

var wireIDCounter int

type Wire struct {
	SourceChip      Chip
	DestinationChip Chip
	id              int
}

type WireError struct {
	id  int
	msg string
}

func (w *WireError) Error() string {
	return fmt.Sprintf("Wire %d : %v", w.id, w.msg)
}

func NewWire(SourceChip Chip, DestinationChip Chip) Wire {
	var result Wire
	result.SourceChip = SourceChip
	result.DestinationChip = DestinationChip
	return result

}

func (w *Wire) Join(SourceChip Chip, DestinationChip Chip) {
	w.SourceChip = SourceChip
	w.DestinationChip = DestinationChip
}

func (w *Wire) ConnectAuto(SourceChip Chip, DestinationChip Chip, moduleID int) (success bool, outPinName string) {

	pins := w.PossibleSourcePinsModule(SourceChip, DestinationChip, moduleID)
	if len(pins) == 0 {
		return false, ""
	}
	// fmt.Printf("Trying to connect to one of %v ", pins)
	w.SourceChip = SourceChip
	w.DestinationChip = DestinationChip
	var pid int = 0
	if len(pins) > 1 {
		var pid int
		fmt.Printf("%s : %s Compatible Output Pins : ", (w.SourceChip).Name(), (w.DestinationChip).Module(moduleID).Name)
		for i := 0; i < len(pins); i++ {
			fmt.Printf(" %s,", (w.SourceChip).PinOut(i).Name)
		}

		fmt.Scanf("\n Possibly more outpins are compatible %d", pid)
	}

	return w.ConnectToModule(moduleID, pins[pid])

}

func (w *Wire) ConnectPins(srcPinName string, moduleName string) (success bool, outPinName string) {

	success = false
	outPinName = ""
	if w.SourceChip == nil && w.DestinationChip == nil {
		fmt.Printf("Both Source and Distination not Set")
		return success, outPinName
	}

	module := w.DestinationChip.ModuleByName(moduleName)
	srcPin := w.SourceChip.PinByName(srcPinName)
	if srcPin.Name == "" {
		fmt.Printf("\n Could not Find Pin %v  in %v", srcPinName, w.SourceChip.Name())
		return success, outPinName
	}

	if module.Name == "" {
		fmt.Printf("\n Could not Find Module %v  in %v", moduleName, w.DestinationChip.Name())
		return success, outPinName
	}

	fmt.Printf("\nWire : Calling ConnectToModule : %v (%v) ,%v (%v)", module.Id, module.Name, srcPin.Id, srcPin.Name)

	return w.ConnectToModule(module.Id, srcPin.Id)
	// return success, outPinName

}

func (w *Wire) PossibleSourcePinsModule(SourceChip Chip, DestinationChip Chip, moduleID int) (matchPins []int) {
	inpins := (DestinationChip).Module(moduleID).InPins
	incnt := len(inpins)
	// incnt := (*DestinationChip).InPinCount()
	outcnt := (SourceChip).OutPinCount()

	matchPins = make([]int, 0, outcnt)
	var connectable bool = true
	if incnt == 0 && outcnt == 0 {
		return matchPins
	}

	// var connectable bool=false
	// var found bool=false
	for j := 0; j < incnt; j++ {

		for i := 0; i < outcnt; i++ {
			// fmt.Printf("\n Checking %v and %v connectable for %v ", (*SourceChip).PinOut(i).Name, (*DestinationChip).PinIn(inpins[j]).Name, (*DestinationChip).Module(moduleID).Name)
			// fmt.Printf("\n Data Types %v and %v connectable for %v ", (*SourceChip).PinOut(i).DataType, (*DestinationChip).PinIn(inpins[j]).DataType, (*DestinationChip).Module(moduleID).Name)
			connectable = connectable && ((SourceChip).PinOut(i).DataType == (DestinationChip).PinIn(inpins[j]).DataType)
			// fmt.Printf("\n connect = %v", connectable)
			if connectable {
				matchPins = append(matchPins, i)
				// fmt.Printf("\n Pins %v and %v , Type (%v) can be connected", (*SourceChip).PinOut(i).Name, (*DestinationChip).PinIn(inpins[j]).Name, (*SourceChip).PinOut(i).DataType)
			}
		}
	}

	return matchPins
}

func (w *Wire) IsModuleConnectable(SourceChip Chip, DestinationChip Chip, moduleID int) (matches int) {

	result := w.PossibleSourcePinsModule(SourceChip, DestinationChip, moduleID)

	return len(result)

}

func (w *Wire) IsConnectable(SourceChip Chip, DestinationChip Chip) (matches int) {
	incnt := (DestinationChip).InPinCount()
	outcnt := (SourceChip).OutPinCount()
	matches = 0
	var connectable bool = false
	if incnt == 0 && outcnt == 0 {
		return matches
	}

	// var connectable bool=false
	// var found bool=false

	for i := 0; i < outcnt; i++ {
		for j := 0; j < incnt; j++ {
			connectable = ((SourceChip).PinOut(i).DataType == (DestinationChip).PinIn(j).DataType)
			if connectable {
				matches++
			}
		}
	}

	return matches

}

func (w *Wire) Connect(SourceChip Chip, DestinationChip Chip, moduleID int, srcPinID int) {

	if (DestinationChip).InPinCount() == 0 || (SourceChip).OutPinCount() == 0 {
		fmt.Printf("Not sufficient pins to connect between them")
	}
	// fmt.Printf("\n Source : %#v", *SourceChip)
	// fmt.Printf("\n Destination %#v", *DestinationChip)
	w.SourceChip = SourceChip
	w.DestinationChip = DestinationChip

	// fmt.Printf("\n Check if SOURCE chip is ready to be read %v", (w.SourceChip).PinOut(1))

	fmt.Println("\n============")
	for i := 0; i < (w.SourceChip).OutPinCount(); i++ {
		// fmt.Printf("\n Check if SOURCE chip PIN %d is ready to be read %v", i, (w.SourceChip).PinOut(i))
		fmt.Printf("\n Check if DEST chip PIN %d is ready  to be WRITE %v", i, (w.DestinationChip).PinIn(i))

		count := (w.DestinationChip).ModulesCount()
		for k := 0; k < count; k++ {
			fmt.Printf("\n Available modules with %s are %s ", (w.DestinationChip).Name(), (w.DestinationChip).Module(k).Name)
		}

	}

	// fmt.Printf("\n Check if DEST chip is ready to be WRITE %v", (w.DestinationChip).PinIn(1))

	/// DEFAULT CONNECT

	fmt.Printf("\n==========  Inspect DEST CHIP ")
	fmt.Printf("\n WIRE : Will try connect to IN from %#v", (w.DestinationChip).PinIn(0))
	// fmt.Printf("\n WIRE : Will try to trigger from %#v", (w.DestinationChip).Module(0))
	// fmt.Printf("\n WIRE : Will try to Fns from %#v", (w.DestinationChip).Module(0).Function.Call(make([]reflect.Value, 0)))

	fmt.Println("==================")
	// w.ConnectPins(0, 1)
	w.ConnectToModule(moduleID, srcPinID)
}

func (w *Wire) JoinPins(pvsChipOut int, nextChipIn int) {
	srcType := (w.SourceChip).PinOut(pvsChipOut).DataType
	destType := (w.DestinationChip).PinIn(nextChipIn).DataType

	fmt.Printf("\n %v == %v ? ", srcType.Name(), destType.Name())

	if srcType == destType {
		// fmt.Printf("\n READY TO CONNECT PIN MATCHES")
		// fmt.Printf("Output Available from Dest Chip at")
		// fmt.Print("\n Channel is available %v", (w.DestinationChip).PinOut())

	}
}

func (w *Wire) ConnectToModule(moduleId int, srcPinId int) (success bool, outPinName string) {
	moduleName := w.DestinationChip.Module(moduleId).Name
	srcType := w.SourceChip.PinByID(srcPinId).DataType
	outpins := w.DestinationChip.Module(moduleId).OutPins
	inpins := w.DestinationChip.Module(moduleId).InPins

	if len(inpins) == 0 {
		fmt.Println("%v:%v does not expect Input from pins", w.DestinationChip.Name(), moduleName)
		return success, outPinName
	}
	// else {
	// 	fmt.Printf("\n %v : Following Pins found as Input %v \n ================", moduleName, w.DestinationChip.PinByID(inpins[0]))
	// }

	// fmt.Printf("\n is this valid %v:%v", w.SourceChip.Name(), w.SourceChip.PinByID(srcPinId))

	nextChipIn := inpins[0]
	destType := (w.DestinationChip).PinByID(nextChipIn).DataType

	if srcType == destType {
		// fmt.Printf("\n READY TO CONNECT PIN MATCHES")
		outPinName = (w.DestinationChip).PinByID(outpins[0]).Name
		// (w.SourceChip).PinOut(srcPinId).Channel
		// temp := (w.DestinationChip).PinIn(nextChipIn)
		// temp.Channel = (w.SourceChip).PinOut(srcPinId).Channel
		// w.SourceChip PinOut(srcPinId) = temp
		nargs := (w.DestinationChip).Module(moduleId).Function.Type().NumIn()
		if nargs > 1 {
			fmt.Printf("\n seems like %v need %v args", moduleName, nargs)

		}
		fmt.Printf("\n===== EXECUTING CHIP %v:%v : O/p = %v @ %v ", (w.DestinationChip).Name(), moduleName, (w.DestinationChip).PinByID(outpins[0]).Name, (w.DestinationChip).PinByID(outpins[0]).Channel)

		in := make([]reflect.Value, nargs)
		for i := 0; i < len(in); i++ {
			// in[0] = make()
			// ch :=
			in[0] = reflect.ValueOf((w.SourceChip).PinByID(srcPinId).Channel)
		}
		// fmt.Printf("\n Will pass this argument %v to %v", in[0], moduleName)
		// (w.DestinationChip)
		go (w.DestinationChip).Module(moduleId).Function.Call(in)
		// fmt.Printf("\n After executing")
		// fmt.Printf("\n SUMMA TEST is available %#v", (w.DestinationChip).PinIn(inpins[0]))
		success = true

	} else {
		fmt.Printf("\n PIN Compability Failed %v == %v ? ", srcType.Name(), destType.Name())
	}
	return success, outPinName
}
