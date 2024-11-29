# Verche: A Tool for Managing and Updating `go.mod`

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/s4bb4t/verche.git
   cd verche
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the project:
   ```bash
   cd cmd/checker
   go build -o verche.exe
   ```

4. Add `verche` to your `PATH`:
   - Move the binary to a directory already in your `PATH` (e.g., `/usr/local/bin`):
     ```bash
     sudo mv verche /usr/local/bin
     ```
   - **OR**, add the current directory to your `PATH` temporarily:
     ```bash
     export PATH=$PATH:$(pwd)
     ```
   - To make it permanent, add this line to your shell configuration file (`~/.bashrc`, `~/.zshrc`, etc.):
     ```bash
     export PATH=$PATH:/path/to/verche/directory
     ```

### Usage

Run the tool with the required flags:
```bash
verche -in <path_to_project_directory> -v <golang_version>
```

#### Example:
```bash
verche -in ./my-project -v 1.21.0
```

### How It Works

1. **Configuration Loading**:
   - Parses command-line flags (`-in` for project path, `-v` for Go version).
   - Validates input and sets up paths for `go.mod` and a temporary `verched_go.mod`.

2. **Package Analysis and Update**:
   - Reads the `go.mod` file line by line.
   - Identifies packages and versions using regex (`liner` package).
   - Sends requests to the package repository (`handler` package) to fetch metadata for the latest permissible version.
   - Updates the `go.mod` with the latest versions and ensures compatibility.

3. **Go Version Adjustment**:
   - Updates the Go version and toolchain references in `go.mod` based on the provided version.

4. **File Handling**:
   - Writes the updated content to a temporary file.
   - Replaces the original `go.mod` with the updated content.

5. **Dependency Cleanup**:
   - Runs `go mod tidy` in the project directory to finalize updates.