package main

import (
	"io"
	"os"
)

func CopyFile(src, dst string) error {
	// Open the source file for reading
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create the destination file for writing
	// Opens the file with RW (read/write) permissions
	// If the file exists, it will be truncated to zero bytes
	// The file's directory must already exist

	// TODO: Add folder creation with a mutex to support concurrent writes
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
