# GoSync

GoSync is a fast and reliable command-line tool written in Go that synchronizes files from a source directory to a destination directory. It detects additions or modifications and copies files accordingly, with optional features like dry-run mode, verbose logging, and concurrent file processing.

## Features

- 🔍 Detects file additions and modifications
- 🧪 Dry-run mode to preview changes without writing
- 📣 Verbose logging for detailed output
- ⚡ Concurrent file hashing and copying

## Usage

```bash
gosync --source /path/to/source --dest /path/to/dest [--dry-run] [--verbose]
```
| Flag        | Description                          |
|-------------|--------------------------------------|
| `--source`  | Source directory to sync from        |
| `--dest`    | Destination directory to sync to     |
| `--dry-run` | Preview changes without copying      |
| `--verbose` | Output detailed logs                 |

## Directory Structure and Component Descriptions
```
/gosync/
├── main.go            # Program entry point
├── cli.go             # CLI argument parsing and config setup
├── walker.go          # Walks source/dest directories and gathers metadata
├── hasher.go          # Calculates file hashes, optionally concurrent
├── comparator.go      # Compares file metadata to detect changes
├── copier.go          # Manages file copying with a worker pool
├── types.go           # Shared structs and constants
├── go.mod             # Go module definition
└── README.md          # Project documentation
```

## How It Works

GoSync works by:

1. Walking both source and destination directories concurrently
2. Calculating SHA256 hashes of all files for comparison
3. Detecting which files are new or modified
4. Using a worker pool to copy files efficiently

The tool uses Go's concurrency model with goroutines and channels to maximize performance on multi-core systems.
