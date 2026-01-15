# Go-Ten

A Go CLI tool for creating new Go projects with interactive template selection.

## Overview

Go-Ten is a terminal-based project initializer built with the Bubble Tea framework. It provides an interactive menu for selecting project templates and scaffolding new Go projects.

## Current Progress

### ✅ Completed Features
- **Interactive CLI**: Terminal-based user interface with keyboard navigation
- **Project Templates**: Selection from predefined templates (Web API, CLI Tool, gRPC Service, Microservice)
- **Bubble Tea Integration**: Clean, responsive TUI with cursor navigation
- **Proper Architecture**: Separated concerns with `cmd/go-ten/` entry point and `internal/prompts/` package
- **Testing**: Comprehensive unit tests for all model logic
- **Build System**: Ready-to-use with `go run ./cmd/go-ten`

### 🚧 Next Steps
- Template scaffolding implementation
- Custom template support
- Configuration options
- Additional project types

## Usage

```bash
# Run the interactive CLI
go run ./cmd/

```


## Requirements

- Go 1.25.3 or later
- Bubble Tea framework (automatically included via go.mod)