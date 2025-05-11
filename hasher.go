package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// calculateHash calculates the SHA256 hash of a file.
// It reads the file in chunks to handle large files efficiently.
func calculateHash(filePath string) (string, error) {
	// fmt.Printf("Calculating hash for file: %s\n", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	buf := make([]byte, 4096)

	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break // end of file reached
			}
			return "", err // error reading file
		}
		_, err = hasher.Write(buf[:n]) //process only the data that was read
		if err != nil {
			return "", err
		}
	}

	hash := fmt.Sprintf("%x", hasher.Sum(nil))
	return hash, nil
}

// HashFileMeta calculates the SHA256 hash of the file specified in the FileMeta struct.
// It then adds that hash vaule to the struct. A new FileMeta is returned with the hash.
func HashFileMeta(fm FileMeta, root string) (FileMeta, error) {
	// call that file above to actually get the hash

	absolutePath := filepath.Join(root, fm.Path) // Use the absolute path

	hash, err := calculateHash(absolutePath)
	if err != nil {
		return fm, err // returns new FileMeta and error values
	}
	fm.Hash = hash
	return fm, nil
}
