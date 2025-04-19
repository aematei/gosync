package main

import (
	"github.com/common-nighthawk/go-figure"
)

func main() {
	// Print a header with go-figure
	header := figure.NewFigure("gosync", "doom", true)
	header.Print()

	// Parse CLI args
	// TODO: Implement parseArgs in cli.go and Config in types.go
	config := parseArgs()
}
