package main

import (
	"io"
	"os"
	"path/filepath"
)

func CopyFile(src, dst string) error {
	// Open the source file for reading
	// Open returns a file descriptor only if err is nil
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create the destination file for writing
	// Opens the file with RW (read/write) permissions
	// If the file exists, it will be truncated to zero bytes
	// The file's directory must already exist

	// Isolate the path from dst
	dstDir := filepath.Dir(dst)
	// Go documentation uses os.Lstat for soome reason
	// https://pkg.go.dev/os#example-FileMode
	fi, err := os.Stat(src)
	if err != nil {
		return err
	}
	// Try to make the dirs with the same permissions as the source file
	err = os.MkdirAll(dstDir, fi.Mode().Perm())
	if err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy the contents from source to destination
	// The Copy function will write an unlimited number
	// of bytes from the source to the destination.
	// This can cause large files to consume lots of
	// memory. CopyN might be a better option depending
	// on user requirements.
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
