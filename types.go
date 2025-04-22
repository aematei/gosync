package main

type CLIConfig struct {
	Source  string
	Dest    string
	DryRun  bool
	Verbose bool
	Watch   bool
}
