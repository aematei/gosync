package main

import (
	"fmt"
)

// compares files in source and destination maps.
// It returns a list of files from source that are new or modified.
func CompareFiles(src, dst map[string]FileMeta, dryRun, verbose bool) []FileMeta {
	var toCopy []FileMeta

	for path, srcMeta := range src {
		dstMeta, exists := dst[path]

		if !exists || srcMeta.Size != dstMeta.Size {
			toCopy = append(toCopy, srcMeta)

			if dryRun {
				fmt.Printf("[Dry Run] Would copy: %s (Size: %d bytes)\n", path, srcMeta.Size)
			} else if verbose {
				fmt.Printf("[Verbose] Scheduled for copy: %s\n", path)
			}
		}
	}

	if verbose && len(toCopy) == 0 {
		fmt.Println("[Verbose] No changes detected. Source and destination are in sync.")
	}

	return toCopy
}
