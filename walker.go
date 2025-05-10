package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// walkDir function walks the directory at the given root path and sends FileMeta structs to the provided channel.
func walkDir(root string, filesChan chan<- FileMeta, errChan chan<- error, verbose bool) {
	var wg sync.WaitGroup
	concurrency := runtime.NumCPU() * 4
	jobs := make(chan FileMeta, concurrency)
	results := make(chan FileMeta, concurrency)

	// start worker goroutines to calculate hashes concurrently
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go ConcurrentHashFileMetaWorker(jobs, results, &wg)
	}

	// use a separate goroutine to forward results to filesChan
	go func() {
		for result := range results {
			filesChan <- result
		}
	}()

	// close channels and wait for workers
	defer func() {
		close(jobs)    // close jobs because walk is done
		wg.Wait()      // wait for workers to finish
		close(results) // close channels now that its done. MUST close after workers finish
		// filesChan is closed by GatherFiles
	}()

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
			Path: filepath.Join(root, relPath), // Use the absolute path
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
}

// ConcurrentHashFileMetaWorker is a worker goroutine that calculates the hash of a file
func ConcurrentHashFileMetaWorker(jobs <-chan FileMeta, results chan<- FileMeta, wg *sync.WaitGroup) {
	defer wg.Done()
	for fm := range jobs {
		hashedFm, err := HashFileMeta(fm) // functions in hasher.go handle this
		if err != nil {
			fmt.Printf("Error calculating hash for %s: %v\n", fm.Path, err)
			hashedFm.Hash = ""
			fmt.Printf("ConcurrentHashFileMetaWorker: Sending to results (error): %s\n", fm.Path)
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

	// uses a WaitGroup to wait for both walks to complete
	var wg sync.WaitGroup
	wg.Add(2)

	// starts the directory walk for the source in a separate goroutine
	go func() {
		defer wg.Done()
		walkDir(sourceRoot, srcFilesChan, errChan, verbose)
	}()

	// starts the directory walk for the destination in a separate goroutine
	go func() {
		defer wg.Done()
		walkDir(destRoot, dstFilesChan, errChan, verbose)
	}()

	srcFilesMap := make(map[string]FileMeta)
	dstFilesMap := make(map[string]FileMeta)
	errors := make([]error, 0) // slice to collect errors

	// Function to read from a filesChan and store into a map.
	readFilesChan := func(filesChan <-chan FileMeta, filesMap map[string]FileMeta) {
		for fileMeta := range filesChan {
			filesMap[fileMeta.Path] = fileMeta
		}
	}

	// reaD from the channels in separate goroutines
	go readFilesChan(srcFilesChan, srcFilesMap)
	go readFilesChan(dstFilesChan, dstFilesMap)

	wg.Wait()
	fmt.Println("All walks complete. Closing channels...")
	close(errChan)
	close(srcFilesChan) // close the filesChans here, after all data is read.
	close(dstFilesChan)

	// collect errors
	for err := range errChan {
		errors = append(errors, err)
	}

	// return the maps and the errors.
	if len(errors) > 0 {
		return srcFilesMap, dstFilesMap, fmt.Errorf("multiple errors: %v", errors) //wrap multiple errors
	}
	return srcFilesMap, dstFilesMap, nil
}
