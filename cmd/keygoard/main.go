package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		if len(os.Args) < 3 {
			fmt.Println("Error: project name required")
			fmt.Println("Usage: keygoard init <project-name>")
			os.Exit(1)
		}
		projectName := os.Args[2]
		if err := initProject(projectName); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "version":
		fmt.Println("keygoard v0.1.0")

	case "help":
		printUsage()

	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("keygoard - TinyGo keyboard firmware framework")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  keygoard init <project-name>  Create a new keyboard project")
	fmt.Println("  keygoard version              Show version information")
	fmt.Println("  keygoard help                 Show this help message")
}
