# Verche

## Overview

**Verche** is a command-line tool designed to help manage and update the `go.mod` file in Go projects. It ensures your dependencies are up to date and compliant with your project's requirements.

---

## Features

1. **Manual and Automatic Modes**:
   - **Manual Mode**: Prompts you to review changes to dependencies manually.
   - **Automatic Mode**: Updates all dependencies to the latest permissible versions.

2. **Dependency Analysis**:
   - Reads the `go.mod` file line by line.
   - Identifies outdated dependencies and fetches the latest versions using an artifact repository.

3. **Go Version Management**:
   - Updates the Go version in your `go.mod` file based on the specified version.

4. **File Management**:
   - Uses a temporary file (`verched_go.mod`) for intermediate updates to prevent data loss.
   - Overwrites `go.mod` with validated updates.

5. **Cleanup and Validation**:
   - Runs `go mod tidy` to ensure dependency coherence after updates.

---

## Installation

### Prerequisites
- **Go**: Version 1.21 or higher.

### Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/s4bb4t/verche.git
   cd verche
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the executable:
   ```bash
   cd cmd/checker
   go build -o verche
   ```

4. Add the executable to your system's PATH:
   ```bash
   sudo mv verche /usr/local/bin/
   ```

---

## Usage

Run Verche with the required flags:

```bash
verche -in <project_directory> -v <golang_version> -m <method>
```

### Flags:
- **`-in`** (required): Path to the project directory containing `go.mod`.
- **`-v`** (optional): Go version to set in `go.mod` (default: `1.23.0`).
- **`-m`** (optional): Method of operation (`manual` or `auto`).

### Example Commands:
1. Update dependencies automatically:
   ```bash
   verche -in ./my_project -v 1.23.0 -m auto
   ```

2. Review updates manually:
   ```bash
   verche -in ./my_project -v 1.23.0 -m manual
   ```

---

## Internals

### Configuration (`pkg/config`)
- Defines project paths and the Go version.
- Loads flags and validates input.
- Constructs file paths for `go.mod` and `verched_go.mod`.

### Update Process (`pkg/updater`)
1. **File Parsing**:
   - Reads `go.mod` line by line.
   - Uses regex to extract dependency names and versions.

2. **Dependency Resolution**:
   - Fetches metadata for the latest permissible versions of dependencies.
   - Resolves conflicts based on semantic versioning.

3. **File Writing**:
   - Updates the `go.mod` file with the latest versions.
   - Uses a temporary file to ensure safe operations.

4. **Cleanup**:
   - Runs `go mod tidy` to finalize changes.

### Error Handling
- Detects and reports issues during file access or dependency resolution.
- Provides detailed error messages for troubleshooting.

---

## Common Issues

### Missing Flags
- Ensure the `-in` flag is provided:
  ```bash
  verche -in ./my_project
  ```

### Permissions
- Ensure you have write permissions for the project directory:
  ```bash
  chmod +w ./my_project/go.mod
  ```

### Dependency Resolution
- If a dependency cannot be resolved, verify its name and repository status.

---

## Contribution

For any questions, contact the maintainer at **s4bb4t@yandex.ru**.