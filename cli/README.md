# CLI Implementations

## 1. Overview

This directory holds language-specific CLI implementations for creating, inspecting, and modifying NovusPack (`.nvpk`) packages.
Each subdirectory is a separate implementation; the binary name indicates the language.

CLI tools in this directory are intended to be **functional implementations of the public API surface**.
They should expose the same capabilities as the NovusPack API (create, inspect, modify, validate) in a way that is usable both from scripts and from an interactive session.

## 2. Purpose and Requirements

High level purpose and requirements.
Use as a loose guideline until more formal requirements have been defined for a given implementation.

### 2.1 Dual Modality

CLI implementations must support:

- **Command-based usage:** Every capability is invokable as a subcommand (e.g. `nvpkg create`, `nvpkg add`, `nvpkg list`) so that scripts and automation can drive the tool without user interaction.
- **Interactive usage:** A REPL (read-eval-print loop) mode where a user opens a package once and then runs list, add, remove, read, etc. without repeating the package path; changes can be batched in memory and written with a single `write` (or equivalent).

All API functionality that is relevant to package creation, inspection, modification, and validation should be available in both modes.

### 2.2 High-Level Capabilities (from API Surface)

Implementations are expected to cover at least the following, in both command and interactive form where applicable:

- **Create** – Create a new empty package (`create <path>` with optional flags).
- **Info** – Show package metadata and summary (`info <path>`).
- **List** – List package contents – paths, sizes (`list <path>`).
- **Add** – Add files or directories (`add <path> <sources>...`).
- **Remove** – Remove a file, directory, or pattern (`remove <path> <entry>`).
- **Read** – Read a file to stdout or file (`read <path> <internal>`).
- **Extract** – Extract all or a subtree to a directory (`extract <path> [-o dir]`).
- **Header** – Print raw package header (`header <path>`).
- **Validate** – Validate package integrity (`validate <path>`).
- **Comment** – Get or set package comment (`comment` get/set/clear).
- **Identity** – Get or set Vendor ID / App ID (`identity` get/set).
- **Metadata** – Show or manipulate package metadata (`metadata <path>`).
- **Interactive** – REPL: open/write, pwd/cd/ls, and all above (`interactive` / `i`).

Exact subcommand names and flags may vary by implementation; the list reflects the current Go CLI (nvpkg) as the reference.

## 3. Implementations

| Binary | Language | Directory        | README / Description                                  |
| ------ | -------- | ---------------- | ----------------------------------------------------- |
| nvpkg  | Go       | [nvpkg/](nvpkg/) | [nvpkg/README.md](nvpkg/README.md) – Cobra-based CLI. |
| nvpkr  | Rust     | (planned)        | Rust implementation.                                  |
| nvpkz  | Zig      | (planned)        | Zig implementation.                                   |

Additional language implementations may be added under their own subdirectories with corresponding binary names and a README that describes build, commands, and how they map to the API.
