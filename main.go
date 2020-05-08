package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/ryanuber/go-glob"
)

const (
	Blue   = "\033[1;34m%s\033[0m"
	Teal   = "\033[1;36m%s\033[0m"
	Yellow = "\033[1;33m%s\033[0m"
	Red    = "\033[1;31m%s\033[0m"
	Green  = "\033[1;32m%s\033[0m"

	ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

	banner = `          
    _______ _______ _______ 
   |\     /|\     /|\     /|
   | +---+ | +---+ | +---+ |
   | |   | | |   | | |   | |
   | |t  | | |f  | | |t  | |
   | +---+ | +---+ | +---+ |
   |/_____\|/_____\|/_____\| making it easy to do dumb stuff
							
   `
)

// Expects something like:
// module.acq_mart_clean.aws_lambda_function.sf_trigger[0]
func cleanCountResource(s string) string {
	index := strings.Split(s, "[")
	return index[0]
}

// Expects something like:
// # module.acq_mart_clean.aws_lambda_function.sf_trigger[0] will be updated in-place
func stripResource(s string) string {
	words := strings.Fields(strings.TrimSpace(s))
	terraformResource := words[1]
	return cleanCountResource(terraformResource)
}

func unique(s []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func main() {

	patternPtr := flag.String("pattern", "*", "The pattern to filter the resources")
	modePtr := flag.String("mode", "include", "Either 'exclude' or 'inlcude'")
	flag.Parse()

	info, _ := os.Stdin.Stat()

	if info.Mode()&os.ModeCharDevice != 0 {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage: terraform plan | tft -pattern=\"module.*.aws_s3_bucket*\" -mode=exclude")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	var runes []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		runes = append(runes, input)

		// Passthrough to see the terraform plan output
		fmt.Print(string(input))
	}

	re := regexp.MustCompile(ansi)

	inputStringNoColors := re.ReplaceAllString(string(runes), "")

	inputStringArr := strings.Split(strings.Replace(inputStringNoColors, "\r\n", "\n", -1), "\n")
	matchingResources := []string{""}
	nonMatchingResources := []string{""}

	for _, s := range inputStringArr {
		if strings.Contains(s, "#") {

			resource := stripResource(s)
			if glob.Glob(*patternPtr, resource) {
				matchingResources = append(matchingResources, resource)
			} else {
				nonMatchingResources = append(nonMatchingResources, resource)
			}
		}
	}

	// Duplicates can occur with count resources
	matchingResources = unique(matchingResources)
	nonMatchingResources = unique(nonMatchingResources)

	fmt.Print(banner)
	fmt.Print("\nResources matching ")
	fmt.Printf(Blue, *patternPtr)
	fmt.Print(" will be ")
	if *modePtr == "include" {
		fmt.Printf(Green, *modePtr+"d")
	} else {
		fmt.Printf(Red, *modePtr+"d")
	}
	fmt.Print(" when running the following command:")

	if *modePtr == "include" {
		if len(matchingResources) == 1 {
			fmt.Println("\n\n \xF0\x9F\x99\x88 no matches found, do nothing ")
		} else {
			fmt.Print("\n\nterraform plan")
			fmt.Println(strings.Join(matchingResources, " -target="))
		}
	} else {
		if len(nonMatchingResources) == 1 {
			fmt.Println("\n\n \xF0\x9F\x99\x88 no matches found, do nothing ")
		} else {
			fmt.Print("\n\nterraform plan")
			fmt.Println(strings.Join(nonMatchingResources, " -target="))
		}
	}
	fmt.Print("\n")

}
