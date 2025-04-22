package main

import (
	"github.com/common-nighthawk/go-figure"
)

func main() {
	// Print a header with go-figure
	header := figure.NewFigure("gosync", "doom", true)
	header.Print()

	// Initialize and execute the CLI
	rootCmd := NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		// Handle any errors from the CLI
		panic(err)
	}
}
