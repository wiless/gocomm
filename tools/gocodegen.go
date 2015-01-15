package main

import (
	// "bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"
	"github.com/wiless/gocomm/chipset"
)

var N int

var USERS int

var infilename string
var generateMain bool
var outputdir string
var templatedir string
var packagename string

func init() {
	flag.StringVar(&infilename, "i", "", "input json file for generating code")
	flag.StringVar(&outputdir, "o", "./", "output dir where the generated code is stored")
	flag.StringVar(&templatedir, "t", "./", "dir where the templates are stored")
	flag.StringVar(&packagename, "pkg", "customchips", "dir where the templates are stored")
	flag.BoolVar(&generateMain, "test", false, "generate test code with main")
}

func main() {
	t := time.Now()
	// infilename := "xx.json"
	// args := os.Args
	// if len(args) > 1 {
	// 	infilename = args[1]
	// }
	// fmt.Printf("\n Input args %v", os.Args)
	flag.Parse()
	if infilename == "" {
		flag.PrintDefaults()
		return
	}
	bytearray, ferr := ioutil.ReadFile(infilename)
	fmt.Printf("\n Reading File %v, \n Error : %v \n ", infilename, ferr)
	var newchipset JsonChip
	err := json.Unmarshal(bytearray, &newchipset)
	fmt.Printf("\n Unmarshalling Error : %v", err)
	fmt.Printf("\n JSON Object %#v", newchipset)
	fmt.Printf("\n Validating Object\n")
	newchipset.WriteTestMain = generateMain
	newchipset.Validate()

	//	newchipset.CreateImplemtation()

	fmt.Println("\n", time.Now())
	fmt.Printf("\nTime Elaspsed %v \n", time.Since(t))
	return

}

type JsonChip struct {
	Name    string
	Pins    []JsonPin
	Modules []JsonModule

	InPinCount    int
	OutPinCount   int
	ModuleCount   int
	PinCount      int
	PinNames      string
	ModuleNames   string
	WriteTestMain bool
}

type JsonModule struct {
	Id        int
	Name      string
	Desc      string
	InPins    []string
	OutPins   []string
	Function  string
	TestNames []string
	ArgString string
	// FunctionName string
}

type JsonPin struct {
	Id         int
	Name       string
	Desc       string
	DataType   string
	InputPin   bool
	ModuleName string
}

func createPin(pjson string) chipset.PinInfo {
	var test chipset.PinInfo
	return test
}

func (chip *JsonChip) FillPins(pin chipset.PinInfo) {
	var result JsonPin

	result.Id = pin.Id
	result.Name = pin.Name
	result.Desc = pin.Desc
	result.DataType = pin.DataType.Name()
	result.InputPin = pin.InputPin
	result.ModuleName = pin.SourceName
	chip.Pins = append(chip.Pins, result)
	// return result
}

func (j *JsonChip) FindPin(pinname string) int {
	result := -1
	for i := 0; i < len(j.Pins); i++ {
		if j.Pins[i].Name == pinname {
			return i
		}
	}
	return result
}

/// Checks if the All objects in the Chip are valid and if code can be generate
func (j *JsonChip) Validate() bool {
	var allpinNeededNames []string
	uniquePins := make(map[string]int)
	j.InPinCount = 0
	j.OutPinCount = 0
	j.ModuleCount = len(j.Modules)

	for i := 0; i < len(j.Modules); i++ {
		allpinNeededNames = append(allpinNeededNames, j.Modules[i].InPins...)

		j.InPinCount += len(j.Modules[i].InPins)
		var argstr string
		for cnt := 0; cnt < len(j.Modules[i].InPins); cnt++ {
			val, ok := uniquePins[j.Modules[i].InPins[cnt]]

			if ok {
				uniquePins[j.Modules[i].InPins[cnt]] = val + 1
			} else {
				uniquePins[j.Modules[i].InPins[cnt]] = 1
			}
			found := j.FindPin(j.Modules[i].InPins[cnt])
			if found != -1 {
				argstr += fmt.Sprintf("%s gocomm.%s,", j.Pins[found].Name, j.Pins[found].DataType)
			}
		}
		j.Modules[i].ArgString = strings.TrimRight(argstr, ",")

	}
	for i := 0; i < len(j.Modules); i++ {
		allpinNeededNames = append(allpinNeededNames, j.Modules[i].OutPins...)
		j.OutPinCount += len(j.Modules[i].OutPins)

		for cnt := 0; cnt < len(j.Modules[i].OutPins); cnt++ {
			val, ok := uniquePins[j.Modules[i].OutPins[cnt]]
			if ok {
				uniquePins[j.Modules[i].OutPins[cnt]] = val + 1
			} else {
				uniquePins[j.Modules[i].OutPins[cnt]] = 1

			}

		}
	}

	for cnt := 0; cnt < len(j.Modules); cnt++ {
		j.ModuleNames += fmt.Sprintf("\"%s\",", j.Modules[cnt].Name)
	}
	j.ModuleNames = strings.TrimRight(j.ModuleNames, ",")

	for key, _ := range uniquePins {
		j.PinNames += fmt.Sprintf("\"%s\",", key)
	}

	j.PinNames = strings.TrimRight(j.PinNames, ",")
	fmt.Printf("\n uniquePins : %v", j.PinNames)
	j.PinCount = j.InPinCount + j.OutPinCount

	success := true

	if !success {
		fmt.Printf("\n Some Pins missing : %v", uniquePins)
	} else {

		/// Generate Code for InitPins

		t, terr := template.ParseFiles(templatedir + "./chip_template.txt")
		if terr != nil {
			fmt.Printf("\t Error  %v ", terr)
			return false
		}
		outfilename := strings.ToLower(j.Name + ".go")
		if !j.WriteTestMain {
			outfilename = outputdir + "/" + outfilename
		}
		fmt.Printf("\n===== AUTO GENERATED %s \n", outfilename)
		fd, ferr := os.Create(outfilename)
		packagestr := "package " + packagename
		fd.WriteString(packagestr)

		if ferr == nil {

			t.Execute(fd, j)
		}
		if j.WriteTestMain {
			fmt.Printf("\n Next Step :  > go run %v\n", outfilename)
		} else {
			fmt.Printf("\n Next Step : \n Structure %v in Package : customchips generated .. \n > go run testauto.go", j.Name)
		}

	}
	return true
}

type Args struct {
	ChipName string
	Name     string
	InArgs   []Arg
	OutArgs  []Arg
}
type Arg struct {
	Variable     string
	VariableType string
}

func (j *JsonChip) CreateImplemtation() bool {
	t, terr := template.ParseFiles(templatedir + "./chip_impl_template.txt")
	if terr != nil {
		fmt.Printf("\t Error  %v ", terr)
		return false
	}
	outfilename := strings.ToLower(j.Name + "_impl.go")
	if !j.WriteTestMain {
		outfilename = outputdir + outfilename
	}
	fmt.Printf("\n Create Implementaiton %s", outfilename)

	fd, ferr := os.Create(outfilename)
	if ferr != nil {
		return false
	}

	packagestr := "package " + packagename
	fd.WriteString(packagestr)

	var moduleargs Args
	count := len(j.Modules)
	moduleargs.ChipName = j.Name
	fmt.Printf("\n Module %v", count)

	for i := 0; i < count; i++ {
		fmt.Printf("\n Module %v", j.Modules[i].Name)
		moduleargs.Name = j.Modules[i].Function
		inCNT := len(j.Modules[i].InPins)
		outCNT := len(j.Modules[i].OutPins)
		moduleargs.InArgs = make([]Arg, inCNT)
		moduleargs.OutArgs = make([]Arg, outCNT)
		for cnt := 0; cnt < inCNT; cnt++ {
			indx := j.FindPin(j.Modules[i].InPins[cnt])
			if indx != -1 {
				moduleargs.InArgs[cnt].Variable = j.Pins[indx].Name
				moduleargs.InArgs[cnt].VariableType = "gocomm." + j.Pins[indx].DataType
			} else {
				fmt.Printf("OHOHOHOHO")
			}

		}
		for cnt := 0; cnt < outCNT; cnt++ {
			indx := j.FindPin(j.Modules[i].OutPins[cnt])
			if indx != -1 {
				moduleargs.OutArgs[cnt].Variable = j.Pins[indx].Name
				moduleargs.OutArgs[cnt].VariableType = "gocomm." + j.Pins[indx].DataType
			} else {
				fmt.Printf("OHOHOHOHO")
			}
		}
		fmt.Printf("\n ARG data %v", moduleargs)

		t.Execute(os.Stdout, moduleargs)
		t.Execute(fd, moduleargs)
	}
	return true
}

func (j *JsonChip) CreateStruct() string {
	return "Dump"
}

func (chip *JsonChip) FillModules(module chipset.ModuleInfo) {
	var result JsonModule

	result.Id = module.Id
	result.Name = module.Name
	result.Desc = module.Desc
	{

		strlist := result.InPins
		count := len(module.InPins)
		for i := 0; i < count; i++ {
			if len(chip.Pins) > module.InPins[i] {
				strlist = append(strlist, chip.Pins[module.InPins[i]].Name)
			}

		}
		result.InPins = strlist
	}
	{
		strlist := result.OutPins
		count := len(module.OutPins)
		for i := 0; i < count; i++ {
			if len(chip.Pins) > module.OutPins[i] {
				strlist = append(strlist, chip.Pins[module.OutPins[i]].Name)
			}

		}
		result.OutPins = strlist
	}
	result.Function = module.FunctionName
	chip.Modules = append(chip.Modules, result)
}
