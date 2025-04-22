package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cfg := &CLIConfig{}

	rootCmd := &cobra.Command{
		Use:   "gosync",
		Short: "GoSync is a fast and reliable file synchronization tool.",
		Long:  "GoSync synchronizes files from a source directory to a destination directory with features like dry-run, verbose logging, live sync, and concurrent processing.",
		Run: func(cmd *cobra.Command, args []string) {
			RunCLI(*cfg)
		},
	}

	rootCmd.Flags().StringVarP(&cfg.Source, "source", "s", "", "Source directory to sync from (required)")
	rootCmd.Flags().StringVarP(&cfg.Dest, "dest", "d", "", "Destination directory to sync to (required)")
	rootCmd.Flags().BoolVarP(&cfg.DryRun, "dry-run", "r", false, "Preview changes without copying")
	rootCmd.Flags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "Output detailed logs")
	rootCmd.Flags().BoolVarP(&cfg.Watch, "watch", "w", false, "Enable live syncing with fsnotify")

	rootCmd.MarkFlagRequired("source")
	rootCmd.MarkFlagRequired("dest")

	return rootCmd
}

func RunCLI(cfg CLIConfig) {
	fmt.Println("🛠 GoSync CLI Starting...")
	fmt.Println("Source Directory:", cfg.Source)
	fmt.Println("Destination Directory:", cfg.Dest)

	if cfg.DryRun {
		fmt.Println("Running in dry-run mode (no files will be copied).")
	} else {
		fmt.Println("Running in normal mode (files will be copied).")
	}

	if cfg.Verbose {
		fmt.Println("Verbose logging is enabled.")
	}

	if cfg.Watch {
		fmt.Println("Live sync is enabled (watch mode).")
	}

	// Stubbed out sync operation
	fmt.Println("🚧 Sync operation not yet implemented.")
}
