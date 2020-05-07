package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

func main() {

	patternPtr := flag.String("pattern", "", "The pattern to filter the resources")
	modePtr := flag.String("mode", "include", "Either 'exclude' or 'inlcude'")
	// operationPtr := flag.String("operation", "plan", "Either 'plan' or 'apply'")
	flag.Parse()

	fmt.Print("Resources matching the pattern ")
	fmt.Printf(InfoColor, *patternPtr)
	fmt.Print(" will be ")
	fmt.Printf(NoticeColor, *modePtr+"d")
	fmt.Print(" in the terraform plan/apply\n\n")

	app := "terraform"

	cmd := exec.Command(app, "state", "list")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	terraformList := strings.Split(strings.Replace(string(stdout), "\r\n", "\n", -1), "\n")

	matchedResources := []string{}
	nonMatchedResources := []string{}

	for _, s := range terraformList {

		if strings.Contains(s, *patternPtr) {
			matchedResources = append(matchedResources, s)
		} else {
			nonMatchedResources = append(nonMatchedResources, s)
		}
	}

	fmt.Println(matchedResources)
	fmt.Println()
	fmt.Println(nonMatchedResources)

}
