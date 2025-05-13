package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// walkDir function walks the directory at the given root path and sends FileMeta structs to the provided channel.
func walkDir(root string, filesChan chan<- FileMeta, errChan chan<- error, verbose bool, wg *sync.WaitGroup) {
	defer wg.Done() // Signal when this walkDir is done (to GatherFiles)

	var hashWg sync.WaitGroup
	concurrency := runtime.NumCPU() * 4
	jobs := make(chan FileMeta, concurrency)
	results := make(chan FileMeta, concurrency)

	// start worker goroutines to calculate hashes concurrently
	for i := 0; i < concurrency; i++ {
		hashWg.Add(1)
		go ConcurrentHashFileMetaWorker(jobs, results, &hashWg, root)
	}

	// Forward results in a separate goroutine (but track when it completes)
	var forwardWg sync.WaitGroup
	forwardWg.Add(1)
	go func() {
		defer forwardWg.Done() // This is crucial - signals when forwarding is done
		for result := range results {
			filesChan <- result
		}
	}()

	// Walk the directory and send jobs
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error walking path %s: %v\n", path, err)
			errChan <- err
			return nil
		}
		if d.IsDir() {
			return nil
		}
		fileInfo, err := d.Info()
		if err != nil {
			fmt.Printf("Error getting file info for %s: %v\n", path, err)
			errChan <- err
			return nil
		}
		// getting relative path
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			fmt.Printf("Error getting relative path for %s: %v\n", path, err)
			errChan <- err
			return nil
		}
		fileMeta := FileMeta{
			Path: relPath,
			Size: fileInfo.Size(),
			Mode: fileInfo.Mode(),
		}
		if verbose {
			fmt.Printf("[Verbose] File found: %s, size: %d, mode: %s\n", fileMeta.Path, fileMeta.Size, fileMeta.Mode)
		}
		jobs <- fileMeta
		return nil
	})

	if err != nil {
		errChan <- err
	}

	// Close jobs when walking is done
	close(jobs)

	// Wait for hash workers to finish
	hashWg.Wait()

	// Close results once all hashing is done
	close(results)

	// Wait for forwarding to finish before returning
	forwardWg.Wait()
}

// ConcurrentHashFileMetaWorker is a worker goroutine that calculates the hash of a file
func ConcurrentHashFileMetaWorker(jobs <-chan FileMeta, results chan<- FileMeta, wg *sync.WaitGroup, root string) {
	defer wg.Done()
	for fm := range jobs {
		hashedFm, err := HashFileMeta(fm, root) // functions in hasher.go handle this
		if err != nil {
			fmt.Printf("Error calculating hash for %s: %v\n", fm.Path, err)
			hashedFm.Hash = ""
			results <- hashedFm
			continue
		}
		results <- hashedFm
	}
}

// GatherFiles walks the directory and returns a map of FileMeta keyed by path.
// It uses walkDir to do the walking and handles errors from walkDir.
func GatherFiles(sourceRoot, destRoot string, verbose bool) (srcFiles, dstFiles map[string]FileMeta, err error) {
	srcFilesChan := make(chan FileMeta)
	dstFilesChan := make(chan FileMeta)
	errChan := make(chan error)

	// Use a WaitGroup to track when walkDir functions complete
	var walkWg sync.WaitGroup
	walkWg.Add(2)

	// Start the directory walks
	go walkDir(sourceRoot, srcFilesChan, errChan, verbose, &walkWg)
	go walkDir(destRoot, dstFilesChan, errChan, verbose, &walkWg)

	srcFilesMap := make(map[string]FileMeta)
	dstFilesMap := make(map[string]FileMeta)
	errors := make([]error, 0) // slice to collect errors

	// Collect errors in a goroutine
	var errDone = make(chan struct{})
	go func() {
		for err := range errChan {
			errors = append(errors, err)
		}
		close(errDone)
	}()

	// Function to read from a filesChan and store into a map
	readFilesChan := func(filesChan <-chan FileMeta, filesMap map[string]FileMeta) {
		for fileMeta := range filesChan {
			filesMap[fileMeta.Path] = fileMeta
		}
	}

	// Read from channels in separate goroutines
	var readWg sync.WaitGroup
	readWg.Add(2)
	go func() {
		readFilesChan(srcFilesChan, srcFilesMap)
		readWg.Done()
	}()
	go func() {
		readFilesChan(dstFilesChan, dstFilesMap)
		readWg.Done()
	}()

	// Wait for all walk operations to complete
	walkWg.Wait()
	fmt.Println("All walks complete. Closing channels...")

	// Safe to close the channels now as all producers are done
	close(srcFilesChan)
	close(dstFilesChan)

	// Wait for reads to complete
	readWg.Wait()

	// Close error channel
	close(errChan)
	<-errDone

	// Return the maps and any errors
	if len(errors) > 0 {
		return srcFilesMap, dstFilesMap, fmt.Errorf("multiple errors: %v", errors)
	}
	return srcFilesMap, dstFilesMap, nil
}
