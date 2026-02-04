# nvpkg

## 1. Overview

nvpkg is a command-line tool for creating, inspecting, and modifying NovusPack (`.nvpk`) packages.
It is built with [Cobra](https://github.com/spf13/cobra) and uses the [NovusPack Go API](../../api/go/README.md) as its only project dependency.

The CLI supports creating empty packages, adding files or directories, listing contents, reading or extracting files from a package, removing entries, and inspecting package headers.

## 2. Requirements

- Go 1.25 or later
- The [api/go](../../api/go/) module (satisfied via `replace` in [go.mod](go.mod) when building from this repository)

## 3. Installation and Build

Build the binary from the `nvpkg` directory or from the repository root.

### 3.1 Build from Source

From the `nvpkg` directory:

#### 3.1.1 Release Build (Smallest Size; Requires [UPX](https://upx.github.io/) on PATH)

```bash
make build
```

Output: `nvpkg` (or `nvpkg.exe` on Windows).
Uses `CGO_ENABLED=0`, `-ldflags="-s -w"`, `-trimpath`, then `upx --best`.
Install UPX if needed (e.g. `apt install upx-ucl`).

#### 3.1.2 Development Build (with Debug Symbols, for Debugging and Stack Traces)

```bash
make build-dev
```

Output: `nvpkg-dev` (or `nvpkg-dev.exe` on Windows).

### 3.2 Build from Repository Root

From the repository root:

```bash
make build-nvpkg     # release binary (nvpkg; requires upx)
make build-dev-nvpkg # development binary (nvpkg-dev)
```

## 4. Commands and Usage

All commands take a package path (path to a `.nvpk` file) where applicable.
Internal paths inside a package use a leading slash (e.g. `/config.json`).

### 4.1 Global Help

```bash
./nvpkg --help
./nvpkg <command> --help
```

### 4.2 Create

Create a new empty NovusPack package at the given path.
The file is written immediately; the package contains no entries until you add them with `add`.

Usage:

```text
nvpkg create <package path> [flags]
```

Flags:

| Flag          | Type   | Description     |
| ------------- | ------ | --------------- |
| `--comment`   | string | Package comment |
| `--vendor-id` | uint32 | Vendor ID       |
| `--app-id`    | uint64 | Application ID  |

Examples:

```bash
./nvpkg create myapp.nvpk
./nvpkg create myapp.nvpk --comment "My application assets"
./nvpkg create myapp.nvpk --vendor-id 1 --app-id 100
```

### 4.3 Info

Show metadata and summary for an existing package (file count, sizes, Vendor ID/App ID when set, comment if set).

Usage:

```text
nvpkg info <package path>
```

Example:

```bash
./nvpkg info myapp.nvpk
```

### 4.4 List

List all files in a package.
Output is one line per file: display path, size, stored size.

Usage:

```text
nvpkg list <package path>
```

Example:

```bash
./nvpkg list myapp.nvpk
```

### 4.5 Add

Add files or directories to a package.
If the package path does not exist, a new package is created and then the sources are added.
If it exists, the package is opened and the sources are added.
After adding, the package is written to disk.

Usage:

```text
nvpkg add <package path> <file or dir> [file or dir ...] [flags]
```

Flags:

| Flag   | Type   | Description                              |
| ------ | ------ | ---------------------------------------- |
| `--as` | string | Store under this path (single file only) |

Examples:

```bash
./nvpkg add myapp.nvpk config.json
./nvpkg add myapp.nvpk config.json --as /config/app.json
./nvpkg add myapp.nvpk ./assets ./data
```

### 4.6 Remove

Remove a file, a directory (all files under a path), or files matching a pattern from a package.
The package is written back to disk after the removal.

Usage:

```text
nvpkg remove <package path> <internal path or pattern> [flags]
```

Flags:

| Flag        | Description                                            |
| ----------- | ------------------------------------------------------ |
| `--pattern` | Treat second argument as a glob pattern (e.g. `*.tmp`) |

- Single file: `remove <pkg> /path/to/file`
- Directory (path ending with `/`): `remove <pkg> /path/to/dir/` removes all files under that path
- Pattern: `remove <pkg> "*.tmp" --pattern`

Examples:

```bash
./nvpkg remove myapp.nvpk /config/old.json
./nvpkg remove myapp.nvpk /cache/
./nvpkg remove myapp.nvpk "*.tmp" --pattern
```

### 4.7 Read

Read a file from a package by its internal path.
Output goes to stdout unless `--output` / `-o` is set.

Usage:

```text
nvpkg read <package path> <internal path> [flags]
```

Flags:

| Flag             | Type   | Description                     |
| ---------------- | ------ | ------------------------------- |
| `-o`, `--output` | string | Write to file instead of stdout |

Examples:

```bash
./nvpkg read myapp.nvpk /config.json
./nvpkg read myapp.nvpk /config.json -o config.json
```

### 4.8 Extract

Extract all or a subtree of files from a package to a directory.
Without an internal path, extracts every file.
With an internal path (e.g. `/docs`), extracts only that file or directory subtree.

Usage:

```text
nvpkg extract <package path> [internal path] [flags]
```

Flags:

| Flag             | Type   | Description                          |
| ---------------- | ------ | ------------------------------------ |
| `-o`, `--output` | string | Directory to extract into (required) |

Examples:

```bash
./nvpkg extract myapp.nvpk -o ./out
./nvpkg extract myapp.nvpk /docs -o ./docs
```

### 4.9 Header

Print the raw package header (magic, format version, index start, flags).
Does not open the full package; only reads the header from disk.

Usage:

```text
nvpkg header <package path>
```

Example:

```bash
./nvpkg header myapp.nvpk
```

### 4.10 Validate

Validate package integrity (header, index, and optional content checks).

Usage:

```text
nvpkg validate <package path>
```

Example:

```bash
./nvpkg validate myapp.nvpk
```

### 4.11 Interactive

Run nvpkg in a read-eval-print loop (REPL).
Use `open <path>` to set the current package; then `list`, `add`, `remove`, and `read` use that path without repeating it.
Package changes (add, remove) stay in-memory until you run `write`; only `write` persists to disk.
A current working directory (cwd) is maintained; `pwd`, `cd`, and `ls` navigate the local filesystem.
Paths support `~` for home (e.g. `open ~/pkg.nvpk`, `cd ~`).
Add takes sources first and optional package path last when no package is open (e.g. `add f1 f2 pkg.nvpk`).
Use `--as <internal path>` when adding a single file to set its path inside the package.
Paths given to `add` are resolved against cwd when relative.
Up/down arrow in interactive mode browses command history.
Tab completes command names and, for `open`, `cd`, `ls`, `create`, and `add`, completes paths from the current working directory.
Use `validate` to check package integrity; `help` for in-session commands, `quit` or `exit` to leave.

Usage:

```text
nvpkg interactive
nvpkg i
```

In-session commands:

| Command                          | Description                                                                   |
| -------------------------------- | ----------------------------------------------------------------------------- |
| `open <path>`                    | Set current package (opens and keeps in memory)                               |
| `close`                          | Clear current package (closes without writing)                                |
| `write`                          | Persist current package to disk (changes stay in-memory until write)          |
| `pwd`                            | Print current working directory                                               |
| `cd [dir]`                       | Change directory (no arg => home)                                             |
| `ls [dir]`                       | List local dir: size, mod date, name (default: cwd)                           |
| `create <path>`                  | Create empty package (flags: `--comment`, etc.)                               |
| `info [path]`                    | Show package info (path optional if open)                                     |
| `list [path]`                    | List contents (path optional if open)                                         |
| `header [path]`                  | Print raw header (path optional if open)                                      |
| `add [src]... [path]`            | Add file(s)/dir(s); path optional if package open; `--as <path>` for one file |
| `remove [path] <internal path>`  | Remove entry                                                                  |
| `read [path] <internal path>`    | Read file; `-o file` to write to file                                         |
| `extract [path] [internal path]` | Extract all or subtree; `-o dir` required                                     |
| `help`                           | Show in-session help                                                          |
| `quit`, `exit`, `q`              | Exit                                                                          |

Example:

```bash
./nvpkg interactive
nvpkg> pwd
/home/user/project
nvpkg> ls
assets/
config.json
nvpkg> open myapp.nvpk
Current package: myapp.nvpk
nvpkg [myapp.nvpk]> add config.json
nvpkg [myapp.nvpk]> list
/config.json  42  40
nvpkg [myapp.nvpk]> quit
```

## 5. Make Targets

The [Makefile](Makefile) provides the same conventions as the root and [api/go](../../api/go/Makefile) Makefiles.

| Target                 | Description                                                               |
| ---------------------- | ------------------------------------------------------------------------- |
| `make test`            | Run all unit tests                                                        |
| `make coverage`        | Run tests with coverage; fail if cmd coverage below COVERAGE_MIN (82%)    |
| `make coverage-90`     | Same as coverage with 90% minimum                                         |
| `make coverage-html`   | Generate HTML coverage report                                             |
| `make coverage-report` | Print coverage summary to the terminal                                    |
| `make lint`            | Run gofmt, go vet, and golangci-lint (see [.golangci.yml](.golangci.yml)) |
| `make build`           | Build release binary: ldflags -s -w + UPX --best (nvpkg; requires upx)    |
| `make build-dev`       | Build development binary with debug symbols (nvpkg-dev)                   |
| `make ci`              | Run tidy, verify, coverage, build, and lint                               |
| `make tidy`            | Run go mod tidy                                                           |
| `make clean`           | Remove binaries (nvpkg, nvpkg-dev) and coverage artifacts                 |

From the repository root you can run `make test-nvpkg`, `make ci-nvpkg`, `make lint-nvpkg`, `make coverage-nvpkg`, `make build-nvpkg`, `make build-dev-nvpkg`, and related targets that delegate into this directory.

## 6. Functionality Scripts

The [scripts/](scripts/) directory contains bash scripts that exercise the CLI.

| Script                                                   | Purpose                                            |
| -------------------------------------------------------- | -------------------------------------------------- |
| [run_all.sh](scripts/run_all.sh)                         | Run all checks (requires built nvpkg or NVPKG set) |
| [check_create.sh](scripts/check_create.sh)               | Create empty package and run info                  |
| [check_header_info.sh](scripts/check_header_info.sh)     | Create package and run header and info             |
| [check_add_list_read.sh](scripts/check_add_list_read.sh) | Add file, list, read (skips if add fails)          |
| [check_remove.sh](scripts/check_remove.sh)               | Add file, remove it, list (skips if add fails)     |

Temporary files are written under the repository `tmp/` directory.
Set `NVPKG` to the path of the binary if it is not at `../nvpkg` relative to the script directory.

## 7. Testing and Coverage

Unit tests live in [cmd/](cmd/) with the pattern `*_test.go`.
Coverage is measured only for the `cmd` package; the minimum is enforced by `make coverage` (default 82%, target 90% once the API path-metadata write is complete).

Run tests:

```bash
make test
```

Run tests with coverage and enforce minimum:

```bash
make coverage
```

### 7.1 TTY test harness (pty)

Interactive mode uses two code paths:

- **Non-TTY (tests):** When `InteractiveStdin` is set (e.g. by tests), the REPL uses a `bufio.Scanner` over that reader. Most tests use this path.
- **TTY:** When `InteractiveStdin` is nil, the REPL uses `runInteractiveWithLiner`, which uses [liner](https://github.com/peterh/liner) for readline-style history and completion. That path is only hit when stdin is a real terminal.

**TestRunInteractive_WithPty** exercises the TTY path by running `nvpkg interactive` in a **subprocess** with a pseudoterminal ([github.com/creack/pty](https://github.com/creack/pty)):

1. The test starts `go run . interactive` with stdin/stdout/stderr attached to the pty slave.
2. The child process sees a TTY and uses `runInteractiveWithLiner`.
3. The test writes scripted input to the pty master (`help\n`, `quit\n`) and reads output from the master.
4. It asserts the output contains the prompt and help text.

The test is skipped on Windows and when pty is unsupported. Because the liner path runs in the child process, it does not increase coverage of the test binary; it validates that the TTY path works when run with a real pty.

## 8. Linting

Linting follows the same strategy as [api/go](../../api/go/): gofmt, go vet, and golangci-lint with [.golangci.yml](.golangci.yml).
Enabled linters include contextcheck, dupl, gocyclo, gocritic, goconst, and gocognit.

Run lint:

```bash
make lint
```

## 9. Related Documentation

- [api/go README](../../api/go/README.md) – NovusPack Go API used by nvpkg
- [Root Makefile](../../Makefile) – Repository-wide targets including nvpkg
- [.github/copilot-instructions.md](../../.github/copilot-instructions.md) – Project and markdown standards
