package main

import "os"

type CLIConfig struct {
	Source  string
	Dest    string
	DryRun  bool
	Verbose bool
	Watch   bool
}

type FileMeta struct {
	Path string // relvative to src or dest
	Size int64
	Mode os.FileMode // for file permissions
	Hash string      // SHA256 hash of the file
}

// Struct used to share information with worker goroutines
type FileToCopyInfo struct {
	SrcPath string
	DstPath string
}
