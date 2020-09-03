package main

import (
	"fmt"
	"gtarsum/version"
	"log"
	"os"
	"strings"
)

func main() {
	l := len(os.Args) - 1
	if l != 1 {
		log.Fatalf("One, and only one filename is expected (instead of %d)", l)
	}
	f := os.Args[1]
	fl := strings.ToLower(f)
	if fl == "-v" || fl == "--version" || fl == "version" {
		fmt.Println(version.String())
		os.Exit(0)
	}
}
