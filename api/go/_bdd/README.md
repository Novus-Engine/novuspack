# BDD Test Infrastructure

This directory contains the Behavior-Driven Development (BDD) test infrastructure for NovusPack, using [godog](https://github.com/cucumber/godog) (Cucumber for Go).

## Reorganization Status

The BDD step definitions have been reorganized from a flat structure into a feature-based organization:

- **Core domain**: Split from `core_steps.go` (12,649 lines) into 5 files (~12,806 lines total)
- **File management domain**: Split from `file_mgmt_steps.go` (3,415 lines) into 7 files (~3,542 lines total)
- **Compression domain**: Split from `compression_steps.go` (2,072 lines) into 5 files (~961 lines total)
- **File format domain**: Split from `file_format_steps.go` (1,968 lines) into 5 files (~2,040 lines total)

All reorganized domains are fully functional and registered in `hooks.go`.

## Structure

```plaintext
_bdd/
├── suite_test.go          # Main test suite entry point
├── support/
│   ├── hooks.go           # Scenario initialization and step registration
│   ├── world.go           # Test world/context management
│   ├── common_steps.go    # Shared/common step definitions
│   └── testdata/          # Test fixtures and sample data
└── steps/
    ├── core/              # Core domain steps (split from core_steps.go)
    │   ├── package_lifecycle.go
    │   ├── package_operations.go
    │   ├── package_info.go
    │   ├── package_properties.go
    │   └── generic_patterns.go
    ├── file_mgmt/         # File management domain steps (split from file_mgmt_steps.go)
    │   ├── file_addition.go
    │   ├── file_removal.go
    │   ├── file_extraction.go
    │   ├── file_queries.go
    │   ├── file_tags.go
    │   ├── file_source.go
    │   └── patterns.go
    ├── compression/       # Compression domain steps (split from compression_steps.go)
    │   ├── operations.go
    │   ├── types.go
    │   ├── configuration.go
    │   ├── streaming.go
    │   └── patterns.go
    └── file_format/        # File format domain steps (split from file_format_steps.go)
        ├── header.go
        ├── file_entry.go
        ├── file_index.go
        ├── signatures.go
        └── parsing.go
```

## Feature Discovery

Feature files are automatically discovered by godog from the [`features/`](../../../features) directory at the repository root.

The BDD test infrastructure is located at [`api/go/_bdd/`](../_bdd), so the feature files are accessed via a relative path.

Configuration in [`suite_test.go`](suite_test.go):

```go
var opt = godog.Options{
    Format: "pretty",
    Paths:  []string{"../../../features"},  // From api/go/_bdd/ to features/
}
```

Godog recursively scans the [`features/`](../../../features) directory and executes all `.feature` files it finds.

## Step Registration

All step definition functions are registered in [`support/hooks.go`](support/hooks.go) within the `InitializeScenario` function.

The registration ensures that step definitions from all domains are available when scenarios execute.

### Registration Order

Steps are registered in a specific order to ensure proper pattern matching:

1. **Common/shared steps** - Registered first via `RegisterCommonSteps()`
2. **Domain-specific steps** - Registered before generic patterns (ensures specific handlers match first)
3. **Generic/consolidated patterns** - Registered last (e.g., `RegisterCoreGenericPatterns()`)

This order ensures that specific step patterns take precedence over generic consolidated patterns.

### Domain-to-Step-File Mapping

#### Reorganized Domains (Feature-Based Structure)

| Feature Directory | Step Registration Function                | Step File                            | Domain Tag            | Phase | Status      |
| ----------------- | ----------------------------------------- | ------------------------------------ | --------------------- | ----- | ----------- |
| `core/`           | `RegisterCoreLifecycleSteps()`            | `steps/core/package_lifecycle.go`    | `@domain:core`        | 1     | ✅ Complete |
| `core/`           | `RegisterCoreOperationsSteps()`           | `steps/core/package_operations.go`   | `@domain:core`        | 1     | ✅ Complete |
| `core/`           | `RegisterCoreInfoSteps()`                 | `steps/core/package_info.go`         | `@domain:core`        | 1     | ✅ Complete |
| `core/`           | `RegisterCorePropertiesSteps()`           | `steps/core/package_properties.go`   | `@domain:core`        | 1     | ✅ Complete |
| `core/`           | `RegisterCoreGenericPatterns()`           | `steps/core/generic_patterns.go`     | `@domain:core`        | 1     | ✅ Complete |
| `file_mgmt/`      | `RegisterFileMgmtAdditionSteps()`         | `steps/file_mgmt/file_addition.go`   | `@domain:file_mgmt`   | 2     | ✅ Complete |
| `file_mgmt/`      | `RegisterFileMgmtRemovalSteps()`          | `steps/file_mgmt/file_removal.go`    | `@domain:file_mgmt`   | 2     | ✅ Complete |
| `file_mgmt/`      | `RegisterFileMgmtExtractionSteps()`       | `steps/file_mgmt/file_extraction.go` | `@domain:file_mgmt`   | 2     | ✅ Complete |
| `file_mgmt/`      | `RegisterFileMgmtQuerySteps()`            | `steps/file_mgmt/file_queries.go`    | `@domain:file_mgmt`   | 2     | ✅ Complete |
| `file_mgmt/`      | `RegisterFileMgmtTagSteps()`              | `steps/file_mgmt/file_tags.go`       | `@domain:file_mgmt`   | 2     | ✅ Complete |
| `file_mgmt/`      | `RegisterFileMgmtSourceSteps()`           | `steps/file_mgmt/file_source.go`     | `@domain:file_mgmt`   | 2     | ✅ Complete |
| `file_mgmt/`      | `RegisterFileMgmtPatterns()`              | `steps/file_mgmt/patterns.go`        | `@domain:file_mgmt`   | 2     | ✅ Complete |
| `compression/`    | `RegisterCompressionOperationsSteps()`    | `steps/compression/operations.go`    | `@domain:compression` | 3     | ✅ Complete |
| `compression/`    | `RegisterCompressionTypesSteps()`         | `steps/compression/types.go`         | `@domain:compression` | 3     | ✅ Complete |
| `compression/`    | `RegisterCompressionConfigurationSteps()` | `steps/compression/configuration.go` | `@domain:compression` | 3     | ✅ Complete |
| `compression/`    | `RegisterCompressionStreamingSteps()`     | `steps/compression/streaming.go`     | `@domain:compression` | 3     | ✅ Complete |
| `compression/`    | `RegisterCompressionPatternsSteps()`      | `steps/compression/patterns.go`      | `@domain:compression` | 3     | ✅ Complete |
| `file_format/`    | `RegisterFileFormatHeaderSteps()`         | `steps/file_format/header.go`        | `@domain:file_format` | 2     | ✅ Complete |
| `file_format/`    | `RegisterFileFormatEntrySteps()`          | `steps/file_format/file_entry.go`    | `@domain:file_format` | 2     | ✅ Complete |
| `file_format/`    | `RegisterFileFormatIndexSteps()`          | `steps/file_format/file_index.go`    | `@domain:file_format` | 2     | ✅ Complete |
| `file_format/`    | `RegisterFileFormatSignatureSteps()`      | `steps/file_format/signatures.go`    | `@domain:file_format` | 2     | ✅ Complete |
| `file_format/`    | `RegisterFileFormatParsingSteps()`        | `steps/file_format/parsing.go`       | `@domain:file_format` | 2     | ✅ Complete |

#### Pending Reorganization (Monolithic Files - Backed Up)

| Feature Directory      | Step Registration Function          | Step File (backed up)              | Domain Tag                    | Phase | Status     |
| ---------------------- | ----------------------------------- | ---------------------------------- | ----------------------------- | ----- | ---------- |
| `basic_ops/`           | `RegisterBasicOpsSteps()`           | `basic_ops_steps.go.bak`           | `@domain:basic_ops`           | 1     | ⏳ Pending |
| `dedup/`               | `RegisterDedupSteps()`              | `dedup_steps.go.bak`               | `@domain:dedup`               | 4     | ⏳ Pending |
| `file_types/`          | `RegisterFileTypesSteps()`          | `file_types_steps.go.bak`          | `@domain:file_types`          | 2     | ⏳ Pending |
| `generics/`            | `RegisterGenericsSteps()`           | `generics_steps.go.bak`            | `@domain:generics`            | 5     | ⏳ Pending |
| `metadata/`            | `RegisterMetadataSteps()`           | `metadata_steps.go.bak`            | `@domain:metadata`            | 4     | ⏳ Pending |
| `metadata_system/`     | `RegisterMetadataSystemSteps()`     | `metadata_system_steps.go.bak`     | `@domain:metadata_system`     | 5     | ⏳ Pending |
| `security/`            | `RegisterSecurityValidationSteps()` | `security_validation_steps.go.bak` | `@domain:security_validation` | 4     | ⏳ Pending |
| `security_encryption/` | `RegisterSecurityEncryptionSteps()` | `security_encryption_steps.go.bak` | `@domain:security_encryption` | 4     | ⏳ Pending |
| `signatures/`          | `RegisterSignaturesSteps()`         | `signatures_steps.go.bak`          | `@domain:signatures`          | 3     | ⏳ Pending |
| `streaming/`           | `RegisterStreamingSteps()`          | `streaming_steps.go.bak`           | `@domain:streaming`           | 3     | ⏳ Pending |
| `testing/`             | `RegisterTestingSteps()`            | `testing_steps.go.bak`             | `@domain:testing`             | 5     | ⏳ Pending |
| `validation/`          | `RegisterValidationSteps()`         | `validation_steps.go.bak`          | `@domain:validation`          | 5     | ⏳ Pending |
| `writing/`             | `RegisterWritingSteps()`            | `writing_steps.go.bak`             | `@domain:writing`             | 3     | ⏳ Pending |

**Note:** All step files include domain tags in package comments and function documentation for better organization and filtering. Reorganized domains use feature-based file structure for improved maintainability.

## Feature File Statistics

**Total Feature Files:** 898

**Testable Feature Files:** ~855 (43 stub files excluded)

**Stub/Placeholder Files:** 43

### Features per Domain

| Domain                 | Feature Count | Testable Files |
| ---------------------- | ------------- | -------------- |
| `file_mgmt/`           | 139           | ~130           |
| `compression/`         | 135           | ~131           |
| `metadata/`            | 95            | ~90            |
| `security/`            | 96            | ~90            |
| `file_format/`         | 74            | ~73            |
| `writing/`             | 51            | ~49            |
| `streaming/`           | 49            | ~46            |
| `signatures/`          | 46            | ~44            |
| `core/`                | 42            | ~39            |
| `basic_ops/`           | 82            | ~79            |
| `file_types/`          | 30            | ~29            |
| `generics/`            | 23            | ~22            |
| `testing/`             | 11            | ~11            |
| `dedup/`               | 14            | ~14            |
| `validation/`          | 6             | ~6             |
| `security_encryption/` | 3             | ~3             |
| `metadata_system/`     | 2             | ~2             |

## Stub Files

Some feature files are kept for reference but contain no testable scenarios.
These files document requirements and higher-level scenarios that are intentionally excluded from BDD execution.

All files in this list are tagged with `@skip` to exclude them from BDD test execution by default.
The BDD system also supports `@wip` tags for work-in-progress features that should be excluded from test execution.

### Reference Files by Domain

- `basic_ops/`: 3 files
  - `basic_operations_definitions.feature`
  - `basic_operations_error_handling.feature`
  - `basic_operations_parameter_specification.feature`
- `compression/`: 4 files
  - `compression_definitions_compression.feature`
  - `compression_error_handling.feature`
  - `compression_parameter_specification.feature`
  - `compression_system_behavior.feature`
- `core/`: 2 files
  - `core_definitions_error_error.feature`
  - `core_error_handling.feature`
- `file_format/`: 1 file
  - `file_format_definitions_structure_format.feature`
- `file_mgmt/`: 6 files
  - `fileentry_directory_association_methods.feature`
  - `file_management_definitions.feature`
  - `file_management_error_handling.feature`
  - `file_management_parameter_specification.feature`
  - `file_management_system_behavior.feature`
  - `file_management_usage_examples.feature`
  - `file_query_operations_information.feature`
- `file_types/`: 1 file
  - `file_type_system_definitions_type_type.feature`
- `generics/`: 1 file
  - `generic_concurrency_patterns_operation.feature`
- `metadata/`: 6 files
  - `directory_metadata_system.feature`
  - `metadata_definitions_format_structure.feature`
  - `metadata_structures.feature`
  - `metadata_usage_examples.feature`
  - `package_information_metadata_information_metadata.feature`
  - `special_metadata_file_types.feature`
- `security/`: 8 files
  - `comment_security_and_injection_prevention.feature`
  - `ml_kem_crystals_kyber_key_exchange.feature`
  - `package_signing_system_signing.feature`
  - `quantum_safe_encryption_system.feature`
  - `security_definitions.feature`
  - `security_layers_and_architecture.feature`
  - `security_structures_validation_structure.feature`
  - `security_usage_examples.feature`
  - `signature_types_and_algorithms.feature`
- `signatures/`: 2 files
  - `signature_definitions_validation_validation.feature`
  - `signature_error_handling.feature`
- `streaming/`: 3 files
  - `streaming_core_types_and_structures_types.feature`
  - `streaming_definitions_streaming.feature`
  - `streaming_usage_examples.feature`
- `writing/`: 2 files
  - `writing_definitions.feature`
  - `writing_usage_examples.feature`

## Requirements Coverage

All feature files must include:

- `@spec(...)` tag linking to the technical specification section
- `@REQ-XXX-NNN` tag linking to the requirement ID
- `@domain:xxx` tag indicating the feature domain

See [`docs/requirements/traceability.md`](../../../docs/requirements/traceability.md) for the complete mapping of specifications to requirements to features.

## Running BDD Tests

From the `api/go/` directory:

```bash
make bdd
```

Or run tests directly:

```bash
go test ./_bdd -v -tags=bdd
```

### Running Domain-Specific Tests

Filter tests by domain tags using the `bdd-domain` target:

```bash
# Run only file_format domain tests
make bdd-domain BDD_DOMAIN='@domain:file_format'

# Run only core domain tests
make bdd-domain BDD_DOMAIN='@domain:core'

# Run only compression domain tests
make bdd-domain BDD_DOMAIN='@domain:compression'
```

Available domain tags:

- `@domain:basic_ops`
- `@domain:core`
- `@domain:file_format`
- `@domain:file_mgmt`
- `@domain:file_types`
- `@domain:compression`
- `@domain:signatures`
- `@domain:streaming`
- `@domain:dedup`
- `@domain:metadata`
- `@domain:metadata_system`
- `@domain:security_validation`
- `@domain:security_encryption`
- `@domain:generics`
- `@domain:validation`
- `@domain:testing`
- `@domain:writing`

You can also use godog tag expressions directly:

```bash
go test -tags=bdd ./_bdd -args --godog.tags='@domain:file_format && ~@skip'
```

From the repository root:

```bash
make bdd-go-v1
```

Or use godog directly (from `api/go/_bdd/`):

```bash
godog
```

## Linting

The project includes a BDD lint script that verifies all feature files have the required tags.

From the `api/go/` directory:

```bash
make bdd-lint
```

Or run directly:

```bash
bash ../../scripts/bdd-lint.sh
```

### Requirement Tagging

Sub-features of variable-length data parsing use `REQ-FILEFMT-015`.
Compression-related features use `REQ-FILEFMT-018` where appropriate.

## Test World

The [`support/world.go`](support/world.go) file provides a test world context that:

- Manages temporary directories for test data
- Provides path resolution utilities
- Handles cleanup after scenarios

Each scenario gets a fresh world instance via Before hooks.

## Step Implementation Status

### Current Status (2025-11-30)

**Step Registration Status:** ✅ 100% Coverage Achieved - All Step Patterns Registered

- **Total Unique Step Patterns:** 11,734 (from feature files)
- **Step Registrations:** 2,725 (consolidated using regex patterns)
- **Step Functions:** 2,581
- **Remaining Undefined:** ✅ **0 steps** (100% registration coverage - verified 2025-12-09)
- **Verification:** All 4 reorganized domains (core, file_mgmt, compression, file_format) tested and confirmed 0 undefined steps
- **Refactoring Approach:** ✅ Using regex patterns to consolidate similar steps (77% reduction: 11,734 → 2,725)
- **Domain Tags:** ✅ All step files include domain and phase tags
- **Compilation Status:** ✅ All step files compile successfully
- **Catch-All Pattern:** ❌ Removed (as requested - using specific patterns instead)
- **Test Execution:** ✅ All steps registered - ready for functional implementation

### Implementation Progress

- ✅ Phase 0: Domain tags added to all step files (@domain:xxx, @phase:N)
- ✅ Phase 0: Step registration complete (2,725 registrations using consolidated patterns)
- ✅ Phase 0: Fixed duplicate function declarations (resolved naming conflicts)
- ✅ Phase 1: Common steps infrastructure implemented (error verification, validation, context management)
- ✅ Phase 2: Core package operations steps improved (creation, opening, closing, validation, info)
- ✅ Phase 3: Basic operations steps improved (error handling, lifecycle, resource management)
- ✅ Phase 4: Test world enhanced (context storage, package metadata, improved cleanup)
- ✅ Phase 1-11: Pattern consolidation complete (generic method calls, properties, types, domain-specific, numeric, bit flags, quoted strings, prepositions/conjunctions)
- ✅ Phase 11 (2025-11-21): High-value consolidations (146 patterns → 9 patterns: Option, Type, Version, Entries, Structure, Hash, Method Call, Remains, Updates/Removes/Gets/Sets)
- ✅ Phase 12 (2025-11-21): Individual pattern registration (381 unique patterns registered)
- ✅ Phase 12 (2025-11-30): Final verification - 0 undefined steps achieved

### Recent Updates

#### Domain Tagging System

- All step files include package-level domain tags
- Each `Register*Steps` function includes domain, phase, and tag documentation
- Tags align with feature file `@domain:xxx` tags for better organization

#### Step Registration Refactoring

- ✅ **Refactored to use regex patterns** - Consolidating similar steps using regex with capture groups and optional groups
- ✅ **Reduced code duplication** - Multiple step variations now handled by single functions
- ✅ **Examples implemented** - File path and FileEntry variations consolidated in `file_mgmt_steps.go`
- ✅ **Documentation updated** - README and STEP_CHECKLIST now include regex pattern guidelines
- ✅ **Compilation successful** - All refactored steps compile correctly
- ✅ **Pattern consolidation complete** - 10 phases of pattern consolidation implemented

#### Pattern Consolidation Phases

- ✅ **Phase 1:** "a/an" Type Instance Patterns (e.g., `^a PackageComment$`, `^an AppID value$`)
- ✅ **Phase 2:** Method Fails Patterns (e.g., `^AddSignature fails$`)
- ✅ **Phase 3:** Type Implementation Patterns (e.g., `^AES256GCMFileHandler implementation$`)
- ✅ **Phase 4:** Enhanced Capitalized Action Patterns (e.g., `^AES support provides...$`)
- ✅ **Phase 5:** Lowercase Action Patterns (e.g., `^compression fails$`, `^configuration enables...$`)
- ✅ **Phase 6:** Numeric Start Patterns (e.g., `^a (\d+)-bit X$`, `^X value (\d+)x(\d+)$`)
- ✅ **Phase 7:** Bit Indicates Patterns (e.g., `^bit (\d+) indicates...$`)
- ✅ **Phase 8:** Quoted String Patterns (e.g., patterns with quoted strings)
- ✅ **Phase 9:** Complex Preposition/Conjunction Patterns (`of`, `that`, `when`, `before/after`, `using`, `during`, `favor`, `and`, `or`)
- ✅ **Phase 10:** All function implementations added and compiling

#### Additional Pattern Consolidation

- ✅ **Phase 1:** Escaped Character Patterns (e.g., `^an io\.Reader$`, `^bit (\d+) \((\d+) of features\)...$`, `^X X\.(\d+)\/PKCS#(\d+)$`)
- ✅ **Phase 2:** Two-Word Capitalized Patterns (e.g., `^Asset Metadata contains...$`, `^API definitions reference...$`)
- ✅ **Phase 3:** Two-Word Lowercase End Patterns (e.g., `^compression configuration$`, `^archival storage requirements$`)
- ✅ **Phase 4:** Verified existing pattern matches (most patterns already covered by previous phases)
- ✅ **Phase 5:** Complex Multi-Word Patterns (e.g., `^a probe result indicating type "([^"]*)"$`)

### Step Registration Completion Status

**Status:** ✅ **COMPLETE** - All step patterns registered (2025-11-30)

- ✅ Extracted all undefined steps from godog test output (11,734 unique step patterns identified)
- ✅ Generated step registration code (2,725 registrations using consolidated patterns)
- ✅ Registered all steps in appropriate domain step files (12 consolidation phases complete)
- ✅ Fixed duplicate function declarations
- ✅ Removed catch-all pattern (replaced with specific patterns for better test fidelity)
- ✅ Final verification: 0 undefined steps - 100% registration coverage achieved

#### Phase 1-5: Functional Implementation

- Many step implementations still contain TODOs and will need API integration once the NovusPack API is implemented
- Step functions currently use placeholder logic that will be replaced with actual API calls

Step implementations are added incrementally as features are developed.

Each step file follows the pattern:

```go
package steps

import (
    "github.com/cucumber/godog"
)

func RegisterXxxSteps(ctx *godog.ScenarioContext) {
    // Steps will be implemented incrementally
}
```

## Step Implementation Roadmap

Step definitions should be implemented in the following priority order:

### Phase 1: Core Infrastructure (Priority: High)

1. **`core/`** (42 files, ~39 testable)

   - Foundation for all other domains
   - Package creation, opening, closing
   - Error handling infrastructure
   - Context management

2. **`basic_ops/`** (82 files, ~79 testable)

   - Basic package operations
   - Lifecycle management
   - Validation and defragmentation

### Phase 2: File Operations (Priority: High)

1. **`file_format/`** (74 files, ~73 testable)

   - Package file format structure
   - File entry parsing
   - Header and metadata parsing

2. **`file_types/`** (30 files, ~29 testable)

   - File type detection
   - Type registration and mapping

3. **`file_mgmt/`** (139 files, ~130 testable)

   - File addition, removal, extraction
   - File queries and search
   - Tag management

### Phase 3: Package Operations (Priority: Medium)

1. **`compression/`** (135 files, ~131 testable)

   - Package compression/decompression
   - Compression type selection
   - Streaming compression

2. **`signatures/`** (46 files, ~44 testable)

   - Digital signature operations
   - Signature validation
   - Signature management

3. **`streaming/`** (49 files, ~46 testable)

   - Streaming file operations
   - Buffer pool management
   - Backpressure handling

4. **`writing/`** (51 files, ~49 testable)

   - Package writing operations
   - SafeWrite and FastWrite
   - Write strategies

### Phase 4: Advanced Features (Priority: Medium)

1. **`security/`** (96 files, ~90 testable)

   - Security validation
   - Encryption type management
   - Security status checking

2. **`metadata/`** (95 files, ~90 testable)

   - Metadata management
   - Tag operations
   - Package information

3. **`dedup/`** (14 files, ~14 testable)

   - Deduplication operations
   - Content block management

### Phase 5: Supporting Features (Priority: Low)

1. **`generics/`** (23 files, ~22 testable)

   - Generic helper functions
   - Concurrency patterns

2. **`validation/`** (6 files, ~6 testable)

   - Input validation
   - Package integrity validation

3. **`testing/`** (11 files, ~11 testable)

   - Testing infrastructure
   - Test coverage requirements

4. **`metadata_system/`** (2 files, ~2 testable)

   - Metadata system operations

5. **`security_encryption/`** (3 files, ~3 testable)

   - Encryption-specific operations
   - Key management

### Common Step Patterns

Many step patterns are shared across domains:

- **Context management**: "Given a valid context", "Given a context with timeout"
- **Package state**: "Given an open NovusPack package", "Given an existing package file"
- **File operations**: "Given a file at path", "When file is added", "Then file exists in package"
- **Error handling**: "Then error is returned", "And error indicates", "And error type is"
- **Validation**: "Then validation passes", "And package is valid", "And structure is correct"

These common patterns should be implemented in shared helper functions or a common steps file.

### Step Registration Best Practices

#### Use Regex Patterns to Consolidate Similar Steps

Instead of creating individual functions for every step variation, use regex patterns with capture groups and optional groups to consolidate similar steps:

##### Good Example - Consolidated Pattern

```go
// Matches multiple variations with a single function
ctx.Step(`^a file path(?: containing only whitespace| in the package| with redundant separators| with relative references)?$`, aFilePathWithVariation)

func aFilePathWithVariation(ctx context.Context, variation string) error {
    // variation will be empty for "a file path", or contain the variation text
    switch variation {
    case " containing only whitespace":
        // Handle whitespace case
    case " in the package":
        // Handle package case
    // ... etc
    default:
        // Handle basic "a file path" case
    }
    return nil
}
```

##### Bad Example - Individual Functions

```go
// Don't do this - creates unnecessary duplication
ctx.Step(`^a file path$`, aFilePath)
ctx.Step(`^a file path containing only whitespace$`, aFilePathContainingOnlyWhitespace)
ctx.Step(`^a file path in the package$`, aFilePathInThePackage)
ctx.Step(`^a file path with redundant separators$`, aFilePathWithRedundantSeparators)
ctx.Step(`^a file path with relative references$`, aFilePathWithRelativeReferences)
```

#### Use Capture Groups for Parameters

Extract parameters from step text using capture groups:

```go
// Extract file path from quoted string
ctx.Step(`^a file at path "([^"]*)"$`, aFileAtPath)

func aFileAtPath(ctx context.Context, path string) error {
    // path contains the captured value
    return nil
}
```

#### Use Alternation for Multiple Options

Use `(?:option1|option2|option3)` for multiple valid options:

```go
// Matches "with" or "without"
ctx.Step(`^a FileEntry (?:with|without) encryption key$`, aFileEntryWithEncryptionKey)

func aFileEntryWithEncryptionKey(ctx context.Context, hasKey string) error {
    // hasKey will be "with" or "without"
    return nil
}
```

#### When to Use Individual Functions

Create separate functions when:

- The logic is significantly different between variations
- The step text is completely unrelated
- The function signature would become too complex with parameters

#### Pattern Matching Guidelines

1. **Use optional groups `(?:...)?` for variations**: `^a file path(?: in the package)?$`
2. **Use capture groups `(...)` for parameters**: `^a file at path "([^"]*)"$`
3. **Use alternation `(?:a|b|c)` for multiple options**: `^a FileEntry (?:with|without) encryption key$`
4. **Anchor patterns**: Always use `^` at start and `$` at end to ensure exact matches
5. **Escape special characters**: Use `\[` for literal brackets, `\.` for literal dots, etc.

#### Examples from Common Steps

See [`support/common_steps.go`](support/common_steps.go) for examples:

- Parameter extraction: `ctx.Step(`^a file at path "([^"]\*)"$`, aFileAtPath)`
- Alternation: `ctx.Step(`^a file with (?:corrupted|invalid) (?:package )?format$`, aFileWithInvalidFormat)`
- Context variations: `ctx.Step(`^a context for (?:package creation|package operations)$`, aContextForPackageOperations)`

## Adding New Features

When adding new feature files:

1. Create the `.feature` file in the appropriate domain directory under [`features/`](../../../features)
2. Include `@spec(...)` and `@REQ-XXX-NNN` tags
3. For reorganized domains, add step implementations to the appropriate file in the domain subdirectory:
   - Core domain: `steps/core/` (package_lifecycle.go, package_operations.go, etc.)
   - File management: `steps/file_mgmt/` (file_addition.go, file_removal.go, etc.)
   - Compression: `steps/compression/` (operations.go, types.go, etc.)
   - File format: `steps/file_format/` (header.go, file_entry.go, etc.)
4. Ensure the corresponding step registration function is called in [`hooks.go`](support/hooks.go)
5. Update [`docs/requirements/traceability.md`](../../../docs/requirements/traceability.md) if adding new requirements

### File Organization Guidelines

- **Target file size**: 200-800 lines per file (maintainability)
- **Feature-based grouping**: Group related steps by functionality, not just domain
- **Registration functions**: One per file, named `Register<Domain><Feature>Steps()`
- **Package declarations**: Use domain package name (e.g., `package core`, `package file_mgmt`)
- **Build tags**: All step files must include `//go:build bdd`

## Verification Checklist

- [x] All feature directories have corresponding step registration functions
- [x] All step definition files exist and match registration function names
- [x] All step files include domain and phase tags
- [x] Godog is configured to discover all feature files from `features/` directory
- [x] All step registrations are called in `InitializeScenario`
- [x] All 898 feature files have both `@spec(...)` and `@REQ-` tags
- [x] All 43 stub files are tagged with `@skip` to exclude from test execution
- [x] Traceability mapping covers all feature directories
- [x] All undefined steps registered (using consolidated regex patterns - 2,725 registrations)
- [x] **100% step registration coverage achieved (0 undefined steps as of 2025-11-30)**
- [ ] Step definitions implemented for testable feature files (~855 files)
- [x] Step registration refactored to use regex patterns (complete - 12 phases of consolidation implemented)
- [x] Catch-all pattern removed (replaced with specific patterns)
- [x] **Final validation complete (2025-11-30): 0 undefined steps verified**
