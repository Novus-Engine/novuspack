# Generic Types and Patterns Requirements

## Core Generic Types

- REQ-GEN-001: Generic helpers meet type safety and behavior guarantees. [api_generics.md#1-core-generic-types](../tech_specs/api_generics.md#1-core-generic-types)
- REQ-GEN-025: PathEntry paths MUST be stored with leading slash [type: constraint]. [api_generics.md#131-pathentry-structure](../tech_specs/api_generics.md#131-pathentry-structure)
- REQ-GEN-026: Path storage rules require leading slash for all paths [type: constraint]. [api_generics.md#133-path-storage-rules](../tech_specs/api_generics.md#133-path-storage-rules)
- REQ-GEN-027: Path validation enforces leading slash requirement [type: constraint]. [api_generics.md#134-validation-rules](../tech_specs/api_generics.md#134-validation-rules)
- REQ-GEN-028: Path display conversion strips leading slash for user output [type: constraint]. [api_core.md#23-path-display-and-extraction](../tech_specs/api_core.md#23-path-display-and-extraction)
- REQ-GEN-029: API methods returning paths for display MUST strip leading slash [type: constraint]. [api_core.md#23-path-display-and-extraction](../tech_specs/api_core.md#23-path-display-and-extraction)
- REQ-GEN-003: Generic types (Option, Result, Collection, Strategy, Validator) provide reusable abstractions [type: architectural]. [api_generics.md#11-option-type](../tech_specs/api_generics.md#11-option-type)
- REQ-GEN-009: Option type usage examples demonstrate optional value patterns [type: documentation-only]. [api_generics.md#117-option-type-usage-examples](../tech_specs/api_generics.md#117-option-type-usage-examples)
- REQ-GEN-010: Result type provides error handling type. [api_generics.md#12-result-type](../tech_specs/api_generics.md#12-result-type)
- REQ-GEN-011: Result type usage examples demonstrate error handling patterns [type: documentation-only]. [api_generics.md#127-result-type-usage-examples](../tech_specs/api_generics.md#127-result-type-usage-examples)
- REQ-GEN-030: PathEntry type provides path entry structure for file and directory paths [type: architectural]. [api_generics.md#13-pathentry-type](../tech_specs/api_generics.md#13-pathentry-type)
- REQ-GEN-031: Collection operations provide functional collection manipulation [type: architectural]. [api_generics.md#14-collection-operations](../tech_specs/api_generics.md#14-collection-operations), [api_generics.md#153-aggregation-operations](../tech_specs/api_generics.md#153-aggregation-operations)
- REQ-GEN-032: Data structure operations provide generic data structure utilities [type: architectural]. [api_generics.md#15-data-structure-operations](../tech_specs/api_generics.md#15-data-structure-operations)
- ~~REQ-GEN-012: Collection interface provides collection operations~~ [type: obsolete] (deprecated - replaced by REQ-GEN-031). [api_generics.md#14-collection-operations](../tech_specs/api_generics.md#14-collection-operations)
- ~~REQ-GEN-013: Basic data structures provide reusable data structures~~ [type: obsolete] (deprecated - replaced by REQ-GEN-032). [api_generics.md#15-data-structure-operations](../tech_specs/api_generics.md#15-data-structure-operations)
- REQ-GEN-014: Strategy interface provides strategy pattern support [type: architectural]. [api_generics.md#16-strategy-interface](../tech_specs/api_generics.md#16-strategy-interface)
- REQ-GEN-015: Validator interface provides validation pattern support. [api_generics.md#17-validator-interface](../tech_specs/api_generics.md#17-validator-interface)

## Generic Function Patterns

- REQ-GEN-004: Generic patterns (Collection Operations, Validation Functions, Factory Functions) provide common patterns [type: architectural]. [api_generics.md#21-collection-operations-generic-function-patterns](../tech_specs/api_generics.md#21-collection-operations-generic-function-patterns)
- REQ-GEN-007: Generic validator functions validate input before processing. [api_generics.md#22-validation-functions](../tech_specs/api_generics.md#22-validation-functions), [api_generics.md#223-composevalidators-function](../tech_specs/api_generics.md#223-composevalidators-function)
- REQ-GEN-020: Factory functions provide type-safe factory patterns. [api_generics.md#23-factory-functions](../tech_specs/api_generics.md#23-factory-functions)
- REQ-GEN-019: Generic function patterns provide function pattern support [type: architectural]. [api_generics.md#2-generic-function-patterns](../tech_specs/api_generics.md#2-generic-function-patterns)

## Generic Concurrency Patterns

- REQ-GEN-016: Generic concurrency patterns provide concurrent operation support [type: architectural]. [api_generics.md#18-generic-concurrency-patterns](../tech_specs/api_generics.md#18-generic-concurrency-patterns)
- REQ-GEN-017: Generic concurrency methods provide concurrent method support. [api_generics.md#19-generic-concurrency-methods](../tech_specs/api_generics.md#19-generic-concurrency-methods)

## Generic Configuration Patterns

- REQ-GEN-018: Generic configuration patterns provide configuration pattern support [type: architectural]. [api_generics.md#110-generic-configuration-patterns](../tech_specs/api_generics.md#110-generic-configuration-patterns)

## Error Handling

- REQ-GEN-002: Error wrapper patterns are consistent [type: architectural]. [api_generics.md#33-error-handling](../tech_specs/api_generics.md#33-error-handling), [api_generics.md#333-error-transformation](../tech_specs/api_generics.md#333-error-transformation)
- REQ-GEN-008: Context errors propagated correctly through generic error handling. [api_generics.md#33-error-handling](../tech_specs/api_generics.md#33-error-handling)

## Type Constraints and Validation

- REQ-GEN-006: Generic type parameters validated (type constraints enforced at compile/runtime) [type: constraint].
  [api_generics.md#32-type-parameter-constraints](../tech_specs/api_generics.md#32-type-parameter-constraints),
  [api_generics.md#324-documenting-constraints](../tech_specs/api_generics.md#324-documenting-constraints)

## Context Integration

- REQ-GEN-005: Generic methods with context parameters respect cancellation/timeout [type: constraint]. [api_core.md#02-context-integration](../tech_specs/api_core.md#02-context-integration)

## Path Normalization and Validation

### Unicode NFC Normalization

- REQ-GEN-034: Path strings MUST be normalized to NFC before storage [type: constraint]. [api_core.md#214-unicode-normalization](../tech_specs/api_core.md#214-unicode-normalization)
- REQ-GEN-035: NFC normalization ensures consistent lookups across platforms. [api_core.md#214-unicode-normalization](../tech_specs/api_core.md#214-unicode-normalization)
- REQ-GEN-036: NFC normalization resolves macOS (NFD) vs Windows/Linux (NFC) differences. [api_core.md#214-unicode-normalization](../tech_specs/api_core.md#214-unicode-normalization)

### Path Length Limits

- REQ-GEN-037: Path format limit allows paths up to 65,535 bytes [type: constraint]. [api_core.md#215-path-length-limits](../tech_specs/api_core.md#215-path-length-limits)
- REQ-GEN-038: Implementation emits info warning for paths > 260 bytes (Windows default limit). [api_core.md#215-path-length-limits](../tech_specs/api_core.md#215-path-length-limits)
- REQ-GEN-039: Implementation emits warning for paths > 1,024 bytes (macOS limit). [api_core.md#215-path-length-limits](../tech_specs/api_core.md#215-path-length-limits)
- REQ-GEN-040: Implementation emits warning for paths > 4,096 bytes (Linux limit). [api_core.md#215-path-length-limits](../tech_specs/api_core.md#215-path-length-limits)
- REQ-GEN-041: Implementation emits warning for paths > 32,767 bytes (Windows extended limit). [api_core.md#215-path-length-limits](../tech_specs/api_core.md#215-path-length-limits)
- REQ-GEN-042: Path length warnings are non-fatal and operation proceeds [type: constraint]. [api_core.md#215-path-length-limits](../tech_specs/api_core.md#215-path-length-limits)
- REQ-GEN-043: Windows extraction automatically uses extended paths for paths > 260 bytes. [api_core.md#215-path-length-limits](../tech_specs/api_core.md#215-path-length-limits)

### Case Sensitivity

- REQ-GEN-044: Paths are stored case-sensitively to preserve exact names [type: constraint]. [api_core.md#221-case-sensitivity](../tech_specs/api_core.md#221-case-sensitivity)
- REQ-GEN-045: Extraction on case-insensitive filesystems errors for case-conflicting paths. [api_core.md#221-case-sensitivity](../tech_specs/api_core.md#221-case-sensitivity)

### Validation Rules Updates

- REQ-GEN-046: Path validation enforces NFC normalization [type: constraint]. [api_generics.md#134-validation-rules](../tech_specs/api_generics.md#134-validation-rules)
- REQ-GEN-047: Path validation rejects null bytes [type: constraint]. [api_generics.md#134-validation-rules](../tech_specs/api_generics.md#134-validation-rules)
- REQ-GEN-048: Path validation rejects trailing slash for files [type: constraint]. [api_generics.md#134-validation-rules](../tech_specs/api_generics.md#134-validation-rules)

## Best Practices

- REQ-GEN-021: Best practices document recommended generic patterns [type: documentation-only].
  [api_generics.md#3-best-practices](../tech_specs/api_generics.md#3-best-practices)
- REQ-GEN-049: Common samber/lo functions are documented for consistent usage [type: documentation-only] (documentation-only: guidance - DO NOT CREATE FEATURE FILE).
  [api_generics.md#211-common-samberlo-functions](../tech_specs/api_generics.md#211-common-samberlo-functions)
- REQ-GEN-050: Guidance documents when to use `any` in generic constraints [type: documentation-only] (documentation-only: guidance - DO NOT CREATE FEATURE FILE).
  [api_generics.md#321-when-to-use-any](../tech_specs/api_generics.md#321-when-to-use-any)
- REQ-GEN-051: Guidance documents when to use `comparable` in generic constraints [type: documentation-only] (documentation-only: guidance - DO NOT CREATE FEATURE FILE).
  [api_generics.md#322-when-to-use-comparable](../tech_specs/api_generics.md#322-when-to-use-comparable)
- REQ-GEN-052: Guidance recommends using generics consistently across the API surface [type: documentation-only] (documentation-only: guidance - DO NOT CREATE FEATURE FILE).
  [api_generics.md#361-always-use-generics](../tech_specs/api_generics.md#361-always-use-generics)
- REQ-GEN-053: Guidance documents generic impact on NovusPack API design [type: documentation-only] (documentation-only: guidance - DO NOT CREATE FEATURE FILE).
  [api_generics.md#3621-impact-on-novuspack-api](../tech_specs/api_generics.md#3621-impact-on-novuspack-api)
- REQ-GEN-054: Documentation includes supported generic patterns example [type: documentation-only] (documentation-only: examples - DO NOT CREATE FEATURE FILE).
  [api_generics.md#3623-supported-example](../tech_specs/api_generics.md#3623-supported-example)
- REQ-GEN-055: Documentation covers future compatibility considerations for generic APIs [type: documentation-only] (documentation-only: guidance - DO NOT CREATE FEATURE FILE).
  [api_generics.md#3624-future-compatibility](../tech_specs/api_generics.md#3624-future-compatibility)
- REQ-GEN-056: Documentation provides usage guidelines for generic patterns [type: documentation-only] (documentation-only: guidance - DO NOT CREATE FEATURE FILE).
  [api_generics.md#36313-usage-guidelines](../tech_specs/api_generics.md#36313-usage-guidelines)
- REQ-GEN-022: Naming conventions define generic type naming rules [type: documentation-only]. [api_generics.md#31-naming-conventions](../tech_specs/api_generics.md#31-naming-conventions)
- REQ-GEN-023: Documentation defines generic type documentation requirements [type: documentation-only]. [api_generics.md#34-documentation](../tech_specs/api_generics.md#34-documentation)
- REQ-GEN-024: Testing defines generic type testing requirements [type: documentation-only]. [api_generics.md#35-testing](../tech_specs/api_generics.md#35-testing)
- REQ-GEN-033: Generic function usage examples demonstrate practical generic function patterns [type: documentation-only]. [api_generics.md#36-generic-function-usage](../tech_specs/api_generics.md#36-generic-function-usage)
