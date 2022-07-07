package output

import (
	"encoding/json"
	"fmt"
	"github.com/deepfence/IOCScanner/core"
	"github.com/fatih/color"
	// "github.com/fatih/color"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	Indent = "  " // Indentation for Json printing
)

type IOCFound struct {
	LayerID          string   `json:"Image Layer ID,omitempty"`
	RuleName         string   `json:"Matched Rule Name,omitempty"`
	StringsToMatch   []string `json:"Matched Part,omitempty"`
	CategoryName     []string `json:"Category,omitempty"`
	Severity         string   `json:"Severity,omitempty"`
	SeverityScore    float64  `json:"Severity Score,omitempty"`
	CompleteFilename string   `json:"Full File Name,omitempty"`
	Meta             []string `json:"rule meta"`
}

type IOCOutput interface {
	WriteIOC(string) error
}

type JsonDirIOCOutput struct {
	Timestamp time.Time
	DirName   string `json:"Directory Name"`
	IOC       []IOCFound
}

type JsonImageIOCOutput struct {
	Timestamp   time.Time
	ImageName   string `json:"Image Name"`
	ImageId     string `json:"Image ID"`
	ContainerId string `json:"Container ID"`
	IOC         []IOCFound
}

func (imageOutput *JsonImageIOCOutput) SetImageName(imageName string) {
	imageOutput.ImageName = imageName
}

func (imageOutput *JsonImageIOCOutput) SetImageId(imageId string) {
	imageOutput.ImageId = imageId
}

func (imageOutput *JsonImageIOCOutput) SetTime() {
	imageOutput.Timestamp = time.Now()
}

func (imageOutput *JsonImageIOCOutput) SetIOC(IOC []IOCFound) {
	imageOutput.IOC = IOC
}

func (imageOutput JsonImageIOCOutput) WriteIOC(outputFilename string) error {
	err := printIOCToJsonFile(imageOutput, outputFilename)
	return err
}

func (dirOutput *JsonDirIOCOutput) SetDirName(dirName string) {
	dirOutput.DirName = dirName
}

func (dirOutput *JsonDirIOCOutput) SetTime() {
	dirOutput.Timestamp = time.Now()
}

func (dirOutput *JsonDirIOCOutput) SetIOC(IOC []IOCFound) {
	dirOutput.IOC = IOC
}

func (dirOutput JsonDirIOCOutput) WriteIOC(outputFilename string) error {
	err := printIOCToJsonFile(dirOutput, outputFilename)
	return err
}

func printIOCToJsonFile(IOCJson interface{}, outputFilename string) error {
	file, err := json.MarshalIndent(IOCJson, "", Indent)
	if err != nil {
		core.GetSession().Log.Error("printIOCToJsonFile: Couldn't format json output: %s", err)
		return err
	}
	err = ioutil.WriteFile(outputFilename, file, os.ModePerm)
	if err != nil {
		core.GetSession().Log.Error("printIOCToJsonFile: Couldn't write json output to file: %s", err)
		return err
	}

	return nil
}

func (imageOutput JsonImageIOCOutput) PrintJsonHeader() {
	fmt.Fprintf(os.Stdout, "{\n")
	fmt.Fprintf(os.Stdout, Indent+"\"Timestamp\": \"%s\",\n", time.Now().Format("2006-01-02 15:04:05.000000000 -07:00"))
	fmt.Fprintf(os.Stdout, Indent+"\"Image Name\": \"%s\",\n", imageOutput.ImageName)
	fmt.Fprintf(os.Stdout, Indent+"\"Image ID\": \"%s\",\n", imageOutput.ImageId)
	fmt.Fprintf(os.Stdout, Indent+"\"IOC\": [\n")
}

func (imageOutput JsonImageIOCOutput) PrintJsonFooter() {
	printJsonFooter()
}

func (dirOutput JsonDirIOCOutput) PrintJsonHeader() {
	fmt.Fprintf(os.Stdout, "{\n")
	fmt.Fprintf(os.Stdout, Indent+"\"Timestamp\": \"%s\",\n", time.Now().Format("2006-01-02 15:04:05.000000000 -07:00"))
	fmt.Fprintf(os.Stdout, Indent+"\"Directory Name\": \"%s\",\n", dirOutput.DirName)
	fmt.Fprintf(os.Stdout, Indent+"\"IOC\": [\n")
}

func (dirOutput JsonDirIOCOutput) PrintJsonFooter() {
	printJsonFooter()
}

func printJsonFooter() {
	fmt.Fprintf(os.Stdout, "\n"+Indent+"]\n")
	fmt.Fprintf(os.Stdout, "}\n")
}

func PrintColoredIOC(IOCs []IOCFound, isFirstIOC *bool, fileScore float64, severity string) {
	for _, IOC := range IOCs {
		printColoredIOCJsonObject(IOC, isFirstIOC, fileScore, severity)
		*isFirstIOC = false
	}
}

// Function to print json object with the matches IOC string in color
// @parameters
// IOC - Structure with details of the IOC found
// isFirstIOC - indicates if some IOC are already printed, used to properly format json
func printColoredIOCJsonObject(IOC IOCFound, isFirstIOC *bool, fileScore float64, severity string) {
	Indent3 := Indent + Indent + Indent

	if *isFirstIOC {
		fmt.Fprintf(os.Stdout, Indent+Indent+"{\n")
	} else {
		fmt.Fprintf(os.Stdout, ",\n"+Indent+Indent+"{\n")
	}

	if IOC.LayerID != "" {
		fmt.Fprintf(os.Stdout, Indent3+"\""+color.BlueString("Image Layer ID")+"\": %s,\n", jsonMarshal(IOC.LayerID))
	}
	fmt.Fprintf(os.Stdout, Indent3+"\""+color.BlueString("Matched Rule Name")+"\": %s,\n", jsonMarshal(IOC.RuleName))
	fmt.Fprintf(os.Stdout, Indent3+"\""+color.BlueString("Strings to match are")+"\":\n")
	for _, c := range IOC.StringsToMatch {
		if len(c) > 0 {
			fmt.Fprintf(os.Stdout, Indent3+Indent3+"\""+c+"\",\n")
		}
	}
	summary := ""
	for _, c := range IOC.CategoryName {
		if len(c) > 0 {
			str := []string{"The file", IOC.CompleteFilename, "has a", c, "match."}
			summary = strings.Join(str, " ")
		}
	}

	//fmt.Fprintf(os.Stdout, Indent3+"\"String to Match\": %s,\n", IOC.StringsToMatch)
	//fmt.Fprintf(os.Stdout, Indent3+"\"File Match Severity\": %s,\n", jsonMarshal(severity))
	//fmt.Fprintf(os.Stdout, Indent3+"\"File Match Severity Score\": %.2f,\n", fileScore)
	fmt.Fprintf(os.Stdout, Indent3+"\""+color.BlueString("Category")+"\": %s,\n", IOC.CategoryName)
	fmt.Fprintf(os.Stdout, Indent3+"\""+color.BlueString("File Name")+"\": %s,\n", jsonMarshal(IOC.CompleteFilename))
	for _, c := range IOC.Meta {
		var metaSplit = strings.Split(c, " : ")
		if len(metaSplit) > 1 {
			fmt.Fprintf(os.Stdout, Indent3+"\""+color.BlueString(metaSplit[0])+"\": "+metaSplit[1])
			if metaSplit[0] == "description" {
				str := []string{"The file has a rule match that ", strings.Replace(metaSplit[1], "\n", "", -1) + "."}
				summary = summary + strings.Join(str, " ")
			} else {
				if len(metaSplit[0]) > 0 {
					str := []string{"The matched rule file's ", metaSplit[0], " is", strings.Replace(metaSplit[1], "\n", "", -1) + "."}
					summary = summary + strings.Join(str, " ")
				}
			}
		}
	}
	if len(summary) > 0 {
		fmt.Fprintf(os.Stdout, Indent3+"\""+color.BlueString("Summary")+"\": %s,\n", color.YellowString(summary))
	}

	fmt.Fprintf(os.Stdout, Indent+Indent+"}\n")
}

func jsonMarshal(input string) string {
	output, _ := json.Marshal(input)
	return string(output)
}

func removeFirstLastChar(input string) string {
	if len(input) <= 1 {
		return input
	}
	return input[1 : len(input)-1]
}
