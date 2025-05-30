package main

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cfg := &CLIConfig{}

	rootCmd := &cobra.Command{
		Use:   "gosync",
		Short: "GoSync is a fast and reliable file synchronization tool.",
		Long:  "GoSync synchronizes files from a source directory to a destination directory with features like dry-run, verbose logging, and concurrent processing.",
		Run: func(cmd *cobra.Command, args []string) {
			RunCLI(*cfg)
		},
	}

	rootCmd.Flags().StringVarP(&cfg.Source, "source", "s", "", "Source directory to sync from (required)")
	rootCmd.Flags().StringVarP(&cfg.Dest, "dest", "d", "", "Destination directory to sync to (required)")
	rootCmd.Flags().BoolVarP(&cfg.DryRun, "dry-run", "r", false, "Preview changes without copying")
	rootCmd.Flags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "Output detailed logs")

	rootCmd.MarkFlagRequired("source")
	rootCmd.MarkFlagRequired("dest")

	return rootCmd
}

func RunCLI(cfg CLIConfig) {
	fmt.Println("🛠 GoSync CLI Starting...")
	fmt.Println("Input Source Directory:", cfg.Source)
	fmt.Println("Input Destination Directory:", cfg.Dest)

	// Convert source and dest to absolute paths
	var err error
	cfg.Source, err = filepath.Abs(cfg.Source)
	if err != nil {
		panic("❌ Invalid source path: " + err.Error())
	}

	cfg.Dest, err = filepath.Abs(cfg.Dest)
	if err != nil {
		panic("❌ Invalid destination path: " + err.Error())
	}

	fmt.Println("Resolved Source Path:", cfg.Source)
	fmt.Println("Resolved Destination Path:", cfg.Dest)

	if cfg.DryRun {
		fmt.Println("Running in dry-run mode (no files will be copied).")
	} else {
		fmt.Println("Running in normal mode (files will be copied).")
	}

	if cfg.Verbose {
		fmt.Println("Verbose logging is enabled.")
	}

	// replacing temporary simulation with calls to GatherFiles (walker.go)
	// slight changes so gatherfiles returns both src and dst
	var src map[string]FileMeta
	var dst map[string]FileMeta
	srcFiles, dstFiles, err := GatherFiles(cfg.Source, cfg.Dest, cfg.Verbose)
	if err != nil {
		fmt.Println("Error gathering files:", err)
		return
	}

	src = srcFiles
	dst = dstFiles

	// Compare and detect changes
	toCopy := CompareFiles(src, dst, cfg.DryRun, cfg.Verbose)

	// If not dry-run, perform actual copying
	if !cfg.DryRun {
		const numWorkers = 10 // Number of concurrent copy operations
		var wg sync.WaitGroup
		tasks := make(chan FileToCopyInfo, numWorkers)

		// Start worker goroutines
		for i := 0; i < numWorkers; i++ {
			go func() {
				for task := range tasks {
					err := CopyFile(task.SrcPath, task.DstPath)
					if err != nil {
						fmt.Printf("❌ Failed to copy to %s: %v\n", task.DstPath, err)
					} else if cfg.Verbose {
						fmt.Printf("✅ Copied to: %s\n", task.DstPath)
					}
					wg.Done()
				}
			}()
		}

		// Distribute copy tasks
		for _, file := range toCopy {
			srcPath := filepath.Join(cfg.Source, file.Path)
			dstPath := filepath.Join(cfg.Dest, file.Path)

			wg.Add(1)
			tasks <- FileToCopyInfo{SrcPath: srcPath, DstPath: dstPath}
		}

		// Close tasks channel
		close(tasks)
		wg.Wait()

		fmt.Println("✅ Sync operation completed.")
	}
}
