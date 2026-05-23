package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		fmt.Printf("Error loading ENV file: %v", err)
	}

	if len(os.Args) < 2 {
		fmt.Println("No arguments")
		// printHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "migration":
		migrateCommand()
	case "run":
		runCommand()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		// printHelp()
		os.Exit(1)
	}
}

// TODO:  write printHelp
// func printHelp() {}

func runCommand()     {}
func migrateCommand() {}
