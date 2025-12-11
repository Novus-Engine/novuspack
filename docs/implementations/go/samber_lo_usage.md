# NovusPack Go Implementation - samber/lo Usage Standards

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. When to Use samber/lo](#1-when-to-use-samberlo)
  - [1.1 High-Value Use Cases](#11-high-value-use-cases)
  - [1.2 When NOT to Use samber/lo](#12-when-not-to-use-samberlo)
  - [1.3 Decision Guidelines](#13-decision-guidelines)
- [2. Code Style Guidelines](#2-code-style-guidelines)
  - [2.1 Import Conventions](#21-import-conventions)
  - [2.2 Function Selection](#22-function-selection)
  - [2.3 Naming and Formatting](#23-naming-and-formatting)
  - [2.4 Integration with Existing Patterns](#24-integration-with-existing-patterns)
- [3. Common Patterns](#3-common-patterns)
  - [3.1 Duplicate Detection](#31-duplicate-detection)
  - [3.2 Aggregation Operations](#32-aggregation-operations)
  - [3.3 Validation Patterns](#33-validation-patterns)
  - [3.4 Search and Filter Operations](#34-search-and-filter-operations)
  - [3.5 Map Operations](#35-map-operations)
- [4. Examples from Codebase](#4-examples-from-codebase)
  - [4.1 Duplicate Detection Example](#41-duplicate-detection-example)
  - [4.2 Aggregation Example](#42-aggregation-example)
  - [4.3 Validation Example](#43-validation-example)
- [5. Best Practices](#5-best-practices)
  - [5.1 Performance Considerations](#51-performance-considerations)
  - [5.2 Error Handling](#52-error-handling)
  - [5.3 Testing Considerations](#53-testing-considerations)
  - [5.4 Code Review Guidelines](#54-code-review-guidelines)

---

## 0. Overview

This document defines the standards and guidelines for using `samber/lo` in the NovusPack Go API implementation.

It provides practical guidance on when and how to leverage `samber/lo` functions to improve code readability and reduce boilerplate while maintaining clarity and performance.

### 0.1 Cross-References

- [Generic Types and Patterns](../../tech_specs/api_generics.md) - Generic types and patterns used in NovusPack
- [samber/lo Documentation](https://lo.samber.dev) - Official library documentation

## 1. When to Use samber/lo

Use `samber/lo` selectively for operations that benefit from declarative, functional-style code.

The library should enhance readability, not obscure simple operations.

### 1.1 High-Value Use Cases

#### 1.1.1 Duplicate Detection

Use `lo.UniqBy()` or `lo.FindDuplicates()` when checking for duplicate values in collections.

This replaces manual map-based tracking and improves code clarity.

#### 1.1.2 Aggregation Operations

Use `lo.SumBy()`, `lo.Reduce()`, or `lo.CountBy()` for accumulating values from collections.

This simplifies manual loops with accumulation variables.

#### 1.1.3 Searching Operations

Use `lo.Find()`, `lo.Contains()`, or `lo.IndexOf()` when searching for specific elements.

This provides early exit behavior and clearer intent than manual loops.

#### 1.1.4 Filtering Collections

Use `lo.Filter()` or `lo.Reject()` when creating new collections based on predicates.

This is more declarative than manual loops with append operations.

#### 1.1.5 Map Transformations

Use `lo.Map()`, `lo.Keys()`, `lo.Values()`, or `lo.Entries()` for map operations.

This simplifies common map manipulation patterns.

### 1.2 When NOT to Use samber/lo

#### 1.2.1 Simple Iteration with Side Effects

Keep explicit `for range` loops when the operation is straightforward and side effects are the primary purpose.

#### 1.2.2 Index-Dependent Operations

Use explicit loops when the index is needed for error messages, logging, or other context.

#### 1.2.3 Complex Error Handling

Use explicit loops when error messages require detailed context including index, position, or other loop state.

#### 1.2.4 Binary I/O Operations

Keep existing binary I/O patterns as they are already well-handled and optimized.

#### 1.2.5 Operations Clearer as Loops

When a simple loop is more readable than a library function, prefer the loop.

### 1.3 Decision Guidelines

Ask these questions when deciding whether to use `samber/lo`:

1. Does the operation match a common pattern (duplicate detection, aggregation, search, filter)?

2. Would the `lo` function make the code more declarative and easier to understand?

3. Is the index needed for error messages or other context?

4. Is the operation simple enough that a loop would be clearer?

If questions 1 and 2 are yes, and questions 3 and 4 are no, use `samber/lo`.

Otherwise, prefer explicit loops.

## 2. Code Style Guidelines

Follow these guidelines to ensure consistent and readable code when using `samber/lo`.

### 2.1 Import Conventions

Import `samber/lo` with the standard alias:

```go
import (
    "github.com/samber/lo"
)
```

Use `lo.` prefix for all function calls to make library usage explicit.

### 2.2 Function Selection

Choose the most specific function for the use case.

Prefer `lo.SumBy()` over `lo.Reduce()` when summing values.

Prefer `lo.Find()` over `lo.Filter()` when searching for a single element.

Prefer `lo.UniqBy()` over manual duplicate detection when checking uniqueness.

### 2.3 Naming and Formatting

Keep predicate functions concise and descriptive.

Use meaningful variable names for results.

Format multi-line function calls with proper indentation:

```go
uniqueIDs := lo.UniqBy(
    entries,
    func(e IndexEntry) uint64 { return e.FileID },
)
```

### 2.4 Integration with Existing Patterns

`samber/lo` functions work seamlessly with existing NovusPack generic types.

#### 2.4.1 With Option[T]

```go
func GetFirstValidEntry(entries []IndexEntry) Option[IndexEntry] {
    entry, found := lo.Find(entries, func(e IndexEntry) bool {
        return e.FileID != 0
    })
    if !found {
        return Option[IndexEntry]{}
    }
    return Option[IndexEntry]{value: entry, set: true}
}
```

#### 2.4.2 With Result[T]

```go
func ValidateUniqueEntries(entries []IndexEntry) Result[[]IndexEntry] {
    unique := lo.UniqBy(entries, func(e IndexEntry) uint64 { return e.FileID })
    if len(unique) != len(entries) {
        return Err[[]IndexEntry](errors.New("duplicate file IDs found"))
    }
    return Ok(entries)
}
```

## 3. Common Patterns

These patterns demonstrate how to use `samber/lo` for common operations in the NovusPack codebase.

### 3.1 Duplicate Detection

#### 3.1.1 Pattern: Check for duplicate FileIDs

```go
func (f *FileIndex) Validate() error {
    unique := lo.UniqBy(f.Entries, func(e IndexEntry) uint64 { return e.FileID })
    if len(unique) != len(f.Entries) {
        return fmt.Errorf("duplicate file IDs found")
    }
    return nil
}
```

#### 3.1.2 Pattern: Find duplicate values

```go
duplicates := lo.FindDuplicates(entries, func(e IndexEntry) uint64 { return e.FileID })
if len(duplicates) > 0 {
    return fmt.Errorf("found %d duplicate entries", len(duplicates))
}
```

### 3.2 Aggregation Operations

#### 3.2.1 Pattern: Sum sizes from multiple collections

```go
func (f *FileEntry) VariableSize() int {
    pathSize := lo.SumBy(f.Paths, func(p PathEntry) int { return p.Size() })
    hashSize := lo.SumBy(f.Hashes, func(h HashEntry) int { return h.Size() })
    optSize := lo.SumBy(f.OptionalData, func(o OptionalDataEntry) int { return o.Size() })
    return pathSize + hashSize + optSize
}
```

#### 3.2.2 Pattern: Count elements matching a condition

```go
compressedCount := lo.CountBy(entries, func(e FileEntry) bool {
    return e.CompressionType != CompressionNone
})
```

### 3.3 Validation Patterns

#### 3.3.1 Pattern: Validate all elements with early exit

```go
func ValidateAllPaths(paths []PathEntry) error {
    invalid, found := lo.Find(paths, func(p PathEntry) bool {
        return p.Validate() != nil
    })
    if found {
        return fmt.Errorf("invalid path found: %w", invalid.Validate())
    }
    return nil
}
```

#### 3.3.2 Pattern: Check if all elements satisfy condition

```go
allValid := lo.EveryBy(entries, func(e IndexEntry) bool {
    return e.FileID != 0
})
if !allValid {
    return fmt.Errorf("some entries have zero FileID")
}
```

### 3.4 Search and Filter Operations

#### 3.4.1 Pattern: Find element by criteria

```go
entry, found := lo.Find(entries, func(e IndexEntry) bool {
    return e.FileID == targetID
})
if !found {
    return fmt.Errorf("entry not found")
}
```

#### 3.4.2 Pattern: Filter collection by predicate

```go
validEntries := lo.Filter(entries, func(e IndexEntry, _ int) bool {
    return e.FileID != 0 && e.Offset > 0
})
```

#### 3.4.3 Pattern: Check if collection contains value

```go
hasEntry := lo.Contains(entries, targetEntry)
```

### 3.5 Map Operations

#### 3.5.1 Pattern: Extract keys from map

```go
keys := lo.Keys(fileMap)
```

#### 3.5.2 Pattern: Extract values from map

```go
values := lo.Values(fileMap)
```

#### 3.5.3 Pattern: Transform map entries

```go
transformed := lo.MapEntries(fileMap, func(k string, v FileEntry) (uint64, string) {
    return v.FileID, k
})
```

## 4. Examples from Codebase

These examples show how to refactor existing code to use `samber/lo`.

### 4.1 Duplicate Detection Example

#### 4.1.1 Before: Manual map tracking

```go
seen := make(map[uint64]bool)
for i, entry := range f.Entries {
    if entry.FileID == 0 {
        return fmt.Errorf("file ID at index %d cannot be zero", i)
    }
    if seen[entry.FileID] {
        return fmt.Errorf("duplicate file ID %d at index %d", entry.FileID, i)
    }
    seen[entry.FileID] = true
}
```

#### 4.1.2 After: Using lo.UniqBy

```go
// Check for zero FileIDs first (still needs index)
for i, entry := range f.Entries {
    if entry.FileID == 0 {
        return fmt.Errorf("file ID at index %d cannot be zero", i)
    }
}

// Check for duplicates
unique := lo.UniqBy(f.Entries, func(e IndexEntry) uint64 { return e.FileID })
if len(unique) != len(f.Entries) {
    return fmt.Errorf("duplicate file IDs found")
}
```

Note: The zero check still uses an explicit loop because it needs the index for error messages.

### 4.2 Aggregation Example

#### 4.2.1 Before: Manual accumulation

```go
size := 0
for _, path := range f.Paths {
    size += path.Size()
}
for _, hash := range f.Hashes {
    size += hash.Size()
}
for _, opt := range f.OptionalData {
    size += opt.Size()
}
```

#### 4.2.2 After: Using lo.SumBy

```go
size := lo.SumBy(f.Paths, func(p PathEntry) int { return p.Size() })
size += lo.SumBy(f.Hashes, func(h HashEntry) int { return h.Size() })
size += lo.SumBy(f.OptionalData, func(o OptionalDataEntry) int { return o.Size() })
```

### 4.3 Validation Example

#### 4.3.1 Before: Full validation loop

```go
for i, path := range f.Paths {
    if err := path.Validate(); err != nil {
        return fmt.Errorf("invalid path at index %d: %w", i, err)
    }
}
```

#### 4.3.2 After: Using lo.Find with early exit

```go
invalid, found := lo.Find(f.Paths, func(p PathEntry) bool {
    return p.Validate() != nil
})
if found {
    return fmt.Errorf("invalid path: %w", invalid.Validate())
}
```

Note: This loses the index information.

If the index is needed for error messages, keep the explicit loop.

## 5. Best Practices

Follow these best practices to ensure effective use of `samber/lo` in the codebase.

### 5.1 Performance Considerations

`samber/lo` functions have minimal performance overhead compared to explicit loops.

Benchmarks show approximately 4% slower execution, which is negligible for most use cases.

The library uses generics with no reflection, ensuring type safety and good performance.

For performance-critical paths, profile before and after refactoring to verify no significant regression.

### 5.2 Error Handling

When using `samber/lo` functions, ensure error messages provide sufficient context.

If index or position information is needed for debugging, prefer explicit loops.

Combine `lo` functions with explicit loops when both declarative operations and detailed error context are needed.

#### 5.2.1 Example: Hybrid approach

```go
// Use lo for duplicate detection
unique := lo.UniqBy(f.Entries, func(e IndexEntry) uint64 { return e.FileID })
if len(unique) != len(f.Entries) {
    // Use explicit loop to find and report specific duplicate
    seen := make(map[uint64]int)
    for i, entry := range f.Entries {
        if prev, exists := seen[entry.FileID]; exists {
            return fmt.Errorf("duplicate file ID %d at indices %d and %d", entry.FileID, prev, i)
        }
        seen[entry.FileID] = i
    }
}
```

### 5.3 Testing Considerations

#### 5.3.1 Testing Code That Uses samber/lo

Test code using `samber/lo` the same way explicit loops are tested.

The `samber/lo` library functions themselves do not need to be tested; they are external dependencies with their own test coverage.

Focus tests on verifying that code using `lo` functions produces the correct results and handles edge cases appropriately.

Ensure test coverage includes edge cases such as empty collections, single elements, and large collections.

When refactoring existing code to use `lo`, verify that behavior remains identical, especially for error messages and edge cases.

#### 5.3.2 Using samber/lo in Test Code

Use `samber/lo` in test code for operations that benefit from declarative patterns, but keep table-driven test iteration as explicit loops.

Use samber/lo in tests for:

- Verifying collections match expected values using `lo.EveryBy()` or `lo.Find()`
- Aggregating test data for assertions (e.g., summing sizes, counting elements)
- Filtering test cases based on conditions when needed
- Checking for duplicates in test data setup

Keep explicit loops in tests for:

- Table-driven test iteration (`for _, tt := range tests`)
- Test setup that requires index-based operations
- Error message construction that needs index context

#### 5.3.3 BDD Step Definitions

BDD step definitions can benefit from `samber/lo` for verification and aggregation operations.

Follow the same guidelines as production code: use `lo` when it improves readability, but keep explicit loops when index information is needed for error messages.

##### 5.3.3.1 Example: Aggregating sizes in BDD steps

```go
// Before: Manual aggregation
pathsSize := 0
for _, p := range entry.Paths {
    pathsSize += p.Size()
}
hashSize := 0
for _, h := range entry.Hashes {
    hashSize += h.Size()
}

// After: Using lo.SumBy
pathsSize := lo.SumBy(entry.Paths, func(p PathEntry) int { return p.Size() })
hashSize := lo.SumBy(entry.Hashes, func(h HashEntry) int { return h.Size() })
```

##### 5.3.3.2 Example: Verifying all elements in BDD steps

```go
// Verify all paths are valid
invalid, found := lo.Find(entry.Paths, func(p PathEntry) bool {
    return p.Validate() != nil
})
if found {
    return fmt.Errorf("invalid path found: %w", invalid.Validate())
}
```

##### 5.3.3.3 When NOT to use samber/lo in BDD steps

Keep explicit loops in BDD steps when:

- Error messages need to include index or position information for debugging
- The step needs to validate and report multiple issues with specific context
- The operation is simple enough that a loop is clearer

BDD steps should prioritize clear, actionable error messages that help identify the specific failing condition.

#### 5.3.4 Test Helpers

Test helper functions can use `samber/lo` for common verification patterns.

Keep helpers focused and avoid over-abstraction.

When test helpers use `lo` functions, ensure they provide clear error messages for debugging test failures.

### 5.4 Code Review Guidelines

When reviewing code that uses `samber/lo`, check:

1. Is the `lo` function the right choice for this operation?

2. Does it improve readability compared to an explicit loop?

3. Are error messages still informative if index information was removed?

4. Is the code consistent with other similar operations in the codebase?

5. Are predicate functions clear and well-named?

If any of these questions raise concerns, suggest using an explicit loop instead.
