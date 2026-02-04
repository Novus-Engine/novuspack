# CLI vs API Surface: Gaps and Recommendations

This document compares the NovusPack Go API surface with the nvpkg CLI and recommends what to implement next.

## Current CLI Coverage

| CLI command | API used                                         |
| ----------- | ------------------------------------------------ |
| create      | NewPackage, Create (CreateWithOptions via flags) |
| info        | OpenPackage, GetInfo                             |
| list        | OpenPackage, ListFiles                           |
| add         | OpenPackage/NewPackage, AddFile, AddDirectory    |
| remove      | OpenPackage, RemoveFile                          |
| read        | OpenPackage, ReadFile                            |
| extract     | OpenPackage, ListFiles, ReadFile                 |
| header      | OpenPackage (raw header read)                    |
| interactive | All of the above in REPL                         |

## Recommended Additions (by priority)

### 1. High value, low effort

- **Validate** – `nvpkg validate <path>`

  - API: `Package.Validate(ctx)`
  - Use: integrity check before/after operations; CI; debugging.
  - Fits REQ-VALID-002 (package integrity validation).

- **Remove by pattern / directory**

  - API: `RemoveFilePattern(ctx, pattern)`, `RemoveDirectory(ctx, dirPath)`
  - CLI: e.g. `remove <pkg> --pattern "*.tmp"` or `remove <pkg> /dir/` (remove all under path).
  - Use: bulk cleanup without scripting single-file remove.

- **Info: show VendorID / AppID**
  - API: already in `GetInfo()` → `PackageInfo.VendorID`, `AppID`.
  - CLI: extend `nvpkg info` (and interactive `info`) to print VendorID/AppID when non-zero.
  - Use: identity inspection; matches create’s `--vendor-id` / `--app-id`.

### 2. Write and lifecycle options

- **SafeWrite** – `nvpkg write ... --safe` or interactive `write --safe`

  - API: `Package.SafeWrite(ctx, overwrite bool)`.
  - Use: avoid overwriting existing file unless `--overwrite`.

- **Defragment** – `nvpkg defragment <path>` or interactive `defragment`

  - API: `Package.Defragment(ctx)`.
  - Use: compact package after many add/remove cycles.

- **FastWrite** (optional)
  - API: `Package.FastWrite(ctx)`.
  - Use: faster write when full integrity check can be skipped (document trade-off).

### 3. Open mode and safety

- **Open read-only**
  - API: `OpenPackageReadOnly(ctx, path)`.
  - Use: `nvpkg list/read/info/extract/validate <path>` without risk of modification; scripting/CI.
  - Could be global flag `--read-only` or separate commands that default to read-only.

### 4. Comment and identity on existing packages

- **Get/Set comment on open package**

  - API: `GetComment()`, `SetComment(comment)`, `ClearComment()`, `HasComment()`.
  - CLI: interactive `comment` / `comment "..."` / `comment --clear`; optional non-interactive `nvpkg comment <path> [--set "..."] [--clear]`.
  - Use: adjust description without recreate.

- **Get/Set identity on open package**
  - API: `GetVendorID`, `SetVendorID`, `GetAppID`, `SetAppID`, `GetPackageIdentity`, `SetPackageIdentity`, etc.
  - CLI: interactive `vendor-id` / `app-id` get/set; optional non-interactive `nvpkg identity <path> [--vendor-id N] [--app-id N]`.
  - Use: fix or set identity after create.

### 5. Add by pattern

- **AddFilePattern**
  - API: `AddFilePattern(ctx, pattern, options)`.
  - CLI: `add <pkg> "*.json"` or `add <pkg> --pattern "*.json"` (with cwd or base path).
  - Use: bulk add by glob without listing files in shell.

### 6. Advanced (lower priority)

- **GetMetadata** – full metadata dump (e.g. `nvpkg metadata <path>` or `--json` on info).

  - API: `Package.GetMetadata()`.
  - Use: tooling, debugging, integration.

- **Path metadata / hierarchy** – ListPaths, GetPathInfo, GetPathHierarchy, ListDirectories.

  - Use: path-centric view vs file-centric list; advanced tooling.

- **Session base / target path** – SetSessionBase, GetSessionBase, SetTargetPath.

  - Use: path derivation and redirecting write; power users and scripts.

- **Lookup by ID/hash/type** – GetFileByFileID, GetFileByHash, GetFileByChecksum, FindEntriesByTag, FindEntriesByType, GetFileCount.

  - Use: dedup inspection, content-addressable lookup, filtering by type/tag.

- **Path metadata API** – AddPathMetadata, UpdatePathMetadata, GetPathConflicts, AssociateFileWithPath, etc.
  - Use: rich metadata workflows; likely need a separate “metadata” subcommand or interactive verbs.

## Summary

| Priority | Feature                                                                       | API surface                  | Suggested CLI surface                 |
| -------- | ----------------------------------------------------------------------------- | ---------------------------- | ------------------------------------- |
| High     | Validate                                                                      | Validate(ctx)                | `nvpkg validate <path>`               |
| High     | Remove pattern/dir                                                            | RemoveFilePattern, RemoveDir | `remove --pattern` / remove dir path  |
| High     | Info VendorID/AppID                                                           | GetInfo (existing)           | Extend `info` output                  |
| Medium   | SafeWrite                                                                     | SafeWrite(ctx, overwrite)    | `write --safe` [--overwrite]          |
| Medium   | Defragment                                                                    | Defragment(ctx)              | `nvpkg defragment <path>`             |
| Medium   | Open read-only                                                                | OpenPackageReadOnly          | `--read-only` or read-only commands   |
| Medium   | Comment get/set                                                               | Get/Set/ClearComment         | interactive + optional `comment`      |
| Medium   | Identity get/set                                                              | Get/Set VendorID/AppID       | interactive + optional `identity`     |
| Lower    | Add by pattern                                                                | AddFilePattern               | `add --pattern "*.json"`              |
| Lower    | FastWrite, GetMetadata, path/hierarchy, session/target, lookup by ID/hash/tag | Various                      | As needed for tooling and power users |

Implementing **validate**, **remove by pattern/directory**, and **info VendorID/AppID** first gives the best payoff for typical use and aligns with existing API and requirements.
