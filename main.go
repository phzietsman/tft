package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	Blue   = "\033[1;34m%s\033[0m"
	Teal   = "\033[1;36m%s\033[0m"
	Yellow = "\033[1;33m%s\033[0m"
	Red    = "\033[1;31m%s\033[0m"
	Green  = "\033[1;32m%s\033[0m"
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
	// Something weird in how its trimming
	terraformResource := words[2]
	return cleanCountResource(terraformResource)
}

func main() {

	patternPtr := flag.String("pattern", "", "The pattern to filter the resources")
	modePtr := flag.String("mode", "include", "Either 'exclude' or 'inlcude'")
	flag.Parse()

	fmt.Println(os.Args)

	info, _ := os.Stdin.Stat()

	if info.Mode()&os.ModeCharDevice != 0 {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage: terraform plan | tft -pattern=aws_s3_bucket -mode=exclude")
		return
	}

	fmt.Print("Resources matching the pattern ")
	fmt.Printf(Blue, *patternPtr)
	fmt.Print(" will be ")
	if *modePtr == "include" {
		fmt.Printf(Green, *modePtr+"d")
	} else {
		fmt.Printf(Red, *modePtr+"d")
	}
	fmt.Print(" in the terraform plan/apply\n\n")

	reader := bufio.NewReader(os.Stdin)
	var runes []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		runes = append(runes, input)
	}

	inputStringArr := strings.Split(strings.Replace(string(runes), "\r\n", "\n", -1), "\n")
	matchingResources := []string{""}
	nonMatchingResources := []string{""}

	for _, s := range inputStringArr {
		if strings.Contains(s, "#") {

			resource := stripResource(s)
			if strings.Contains(resource, *patternPtr) {
				matchingResources = append(matchingResources, resource)
			} else {
				nonMatchingResources = append(nonMatchingResources, resource)
			}
		}
	}

	if *modePtr == "include" {
		fmt.Println(strings.Join(matchingResources, " -target="))
	} else {
		fmt.Println(strings.Join(nonMatchingResources, " -target="))
	}

}
