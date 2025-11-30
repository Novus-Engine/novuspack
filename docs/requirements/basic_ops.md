# Basic Operations API Requirements

## Package Structure and Constants

- REQ-API_BASIC-004: Package structure follows format specification [type: architectural]. [api_basic_operations.md#11-package-structure](../tech_specs/api_basic_operations.md#11-package-structure)
- REQ-API_BASIC-005: Package format constants are defined and accessible. [api_basic_operations.md#21-package-format-constants](../tech_specs/api_basic_operations.md#21-package-format-constants)
- REQ-API_BASIC-021: Package structure and loading define package organization [type: architectural]. [api_basic_operations.md#1-package-structure-and-loading](../tech_specs/api_basic_operations.md#1-package-structure-and-loading)
- REQ-API_BASIC-022: Package loading process defines package loading workflow [type: architectural]. [api_basic_operations.md#12-package-loading-process](../tech_specs/api_basic_operations.md#12-package-loading-process)
- REQ-API_BASIC-023: Package constants define package format constants [type: architectural]. [api_basic_operations.md#2-package-constants](../tech_specs/api_basic_operations.md#2-package-constants)
- REQ-API_BASIC-024: Constants provide package format values. [api_basic_operations.md#211-constants](../tech_specs/api_basic_operations.md#211-constants)
- REQ-API_BASIC-025: Package lifecycle operations define package lifecycle management [type: architectural]. [api_basic_operations.md#3-package-lifecycle-operations](../tech_specs/api_basic_operations.md#3-package-lifecycle-operations)

## Package Creation

- REQ-API_BASIC-001: NewPackage creates an empty, valid container in memory. [api_basic_operations.md#41-package-constructor](../tech_specs/api_basic_operations.md#41-package-constructor)
- REQ-API_BASIC-006: Create validates path and directory, configures package in memory. [api_basic_operations.md#42-create-method](../tech_specs/api_basic_operations.md#42-create-method)
- REQ-API_BASIC-007: CreateWithOptions uses Create internally and applies options. [api_basic_operations.md#43-create-with-options](../tech_specs/api_basic_operations.md#43-create-with-options)
- REQ-API_BASIC-008: PackageBuilder pattern supports fluent package creation. [api_basic_operations.md#44-package-builder-pattern](../tech_specs/api_basic_operations.md#44-package-builder-pattern)
- REQ-API_BASIC-020: Create validates target directory exists and is writable, fails if directory missing [type: constraint]. [api_basic_operations.md#42-create-method](../tech_specs/api_basic_operations.md#42-create-method)
- REQ-API_BASIC-026: Package creation defines package creation operations [type: architectural]. [api_basic_operations.md#4-package-creation](../tech_specs/api_basic_operations.md#4-package-creation)
- REQ-API_BASIC-027: NewPackage behavior defines constructor behavior. [api_basic_operations.md#411-newpackage-behavior](../tech_specs/api_basic_operations.md#411-newpackage-behavior)
- REQ-API_BASIC-028: NewPackage creates Package instances in memory that are ready for Create or Open operations. [api_basic_operations.md#412-newpackage-example-usage](../tech_specs/api_basic_operations.md#412-newpackage-example-usage)
- REQ-API_BASIC-029: Create parameters define creation parameters. [api_basic_operations.md#421-create-parameters](../tech_specs/api_basic_operations.md#421-create-parameters)
- REQ-API_BASIC-030: Create behavior defines creation process. [api_basic_operations.md#422-create-behavior](../tech_specs/api_basic_operations.md#422-create-behavior)
- REQ-API_BASIC-031: Create error conditions define creation errors. [api_basic_operations.md#423-create-error-conditions](../tech_specs/api_basic_operations.md#423-create-error-conditions)
- ~~REQ-API_BASIC-032: Create example usage demonstrates creation usage~~ [type: documentation-only] (documentation-only: examples - DO NOT CREATE FEATURE FILE). [api_basic_operations.md#424-create-example-usage](../tech_specs/api_basic_operations.md#424-create-example-usage)
- REQ-API_BASIC-033: CreateWithOptions parameters define option-based creation. [api_basic_operations.md#431-parameters](../tech_specs/api_basic_operations.md#431-parameters)
- REQ-API_BASIC-034: CreateWithOptions behavior defines option-based creation process. [api_basic_operations.md#432-behavior](../tech_specs/api_basic_operations.md#432-behavior)
- REQ-API_BASIC-035: CreateWithOptions error conditions define option-based creation errors. [api_basic_operations.md#433-createwithoptions-error-conditions](../tech_specs/api_basic_operations.md#433-createwithoptions-error-conditions)
- ~~REQ-API_BASIC-036: CreateWithOptions example usage demonstrates option-based creation~~ [type: documentation-only] (documentation-only - has testable scenarios in `createwithoptions_example_usage_creation.feature`). [api_basic_operations.md#434-createwithoptions-example-usage](../tech_specs/api_basic_operations.md#434-createwithoptions-example-usage)
- REQ-API_BASIC-037: PackageBuilder purpose defines builder pattern interface. [api_basic_operations.md#441-purpose](../tech_specs/api_basic_operations.md#441-purpose)
- ~~REQ-API_BASIC-038: PackageBuilder example usage demonstrates builder pattern~~ [type: documentation-only] (documentation-only - has testable scenarios in `packagebuilder_example_usage.feature` and `packagebuilder_purpose_and_usage_creation.feature`). [api_basic_operations.md#442-example-usage](../tech_specs/api_basic_operations.md#442-example-usage)

## Package Opening

- REQ-API_BASIC-002: OpenPackage validates format and returns structured errors. [api_basic_operations.md#51-open-method](../tech_specs/api_basic_operations.md#51-open-method)
- REQ-API_BASIC-009: OpenWithValidation opens package and performs full validation. [api_basic_operations.md#52-open-with-validation](../tech_specs/api_basic_operations.md#52-open-with-validation)
- REQ-API_BASIC-039: Package opening defines package opening operations [type: architectural]. [api_basic_operations.md#5-package-opening](../tech_specs/api_basic_operations.md#5-package-opening)
- REQ-API_BASIC-040: Open parameters define opening parameters. [api_basic_operations.md#511-open-parameters](../tech_specs/api_basic_operations.md#511-open-parameters)
- REQ-API_BASIC-041: Open behavior defines opening process. [api_basic_operations.md#512-open-behavior](../tech_specs/api_basic_operations.md#512-open-behavior)
- REQ-API_BASIC-042: Open error conditions define opening errors. [api_basic_operations.md#513-open-error-conditions](../tech_specs/api_basic_operations.md#513-open-error-conditions)
- ~~REQ-API_BASIC-043: Open example usage demonstrates opening usage~~ [type: documentation-only] (documentation-only - has testable scenarios in `open_example_usage_opening.feature`). [api_basic_operations.md#514-open-example-usage](../tech_specs/api_basic_operations.md#514-open-example-usage)
- REQ-API_BASIC-044: OpenWithValidation behavior defines validation opening process. [api_basic_operations.md#521-openwithvalidation-behavior](../tech_specs/api_basic_operations.md#521-openwithvalidation-behavior)
- REQ-API_BASIC-045: OpenWithValidation error conditions define validation opening errors. [api_basic_operations.md#522-openwithvalidation-error-conditions](../tech_specs/api_basic_operations.md#522-openwithvalidation-error-conditions)

## Package Closing

- REQ-API_BASIC-003: Close flushes, finalizes, and releases. [api_basic_operations.md#61-close-method](../tech_specs/api_basic_operations.md#61-close-method)
- REQ-API_BASIC-010: CloseWithCleanup closes package and performs cleanup operations. [api_basic_operations.md#62-close-with-cleanup](../tech_specs/api_basic_operations.md#62-close-with-cleanup)
- REQ-API_BASIC-046: Package closing defines package closing operations [type: architectural]. [api_basic_operations.md#6-package-closing](../tech_specs/api_basic_operations.md#6-package-closing)
- REQ-API_BASIC-047: Close behavior defines closing process. [api_basic_operations.md#611-close-behavior](../tech_specs/api_basic_operations.md#611-close-behavior)
- REQ-API_BASIC-048: Close error conditions define closing errors. [api_basic_operations.md#612-close-error-conditions](../tech_specs/api_basic_operations.md#612-close-error-conditions)
- ~~REQ-API_BASIC-049: Close example usage demonstrates closing usage~~ [type: documentation-only] (documentation-only - has testable scenarios in `close_example_usage_error.feature`). [api_basic_operations.md#613-close-example-usage](../tech_specs/api_basic_operations.md#613-close-example-usage)
- REQ-API_BASIC-050: CloseWithCleanup behavior defines cleanup closing process. [api_basic_operations.md#621-closewithcleanup-behavior](../tech_specs/api_basic_operations.md#621-closewithcleanup-behavior)

## Package Operations

- REQ-API_BASIC-011: Validate validates package format, structure, and integrity. [api_basic_operations.md#71-package-validation](../tech_specs/api_basic_operations.md#71-package-validation)
- REQ-API_BASIC-012: Defragment optimizes package structure and removes unused space. [api_basic_operations.md#72-package-defragmentation](../tech_specs/api_basic_operations.md#72-package-defragmentation)
- REQ-API_BASIC-013: GetInfo returns comprehensive package information. [api_basic_operations.md#73-package-information](../tech_specs/api_basic_operations.md#73-package-information)
- REQ-API_BASIC-014: ReadHeader reads package header from reader without opening package. [api_basic_operations.md#74-header-inspection](../tech_specs/api_basic_operations.md#74-header-inspection)
- REQ-API_BASIC-015: Package state methods (IsOpen, IsReadOnly, GetPath) report current state. [api_basic_operations.md#75-check-package-state](../tech_specs/api_basic_operations.md#75-check-package-state)
- REQ-API_BASIC-051: Package operations define package utility operations [type: architectural]. [api_basic_operations.md#7-package-operations](../tech_specs/api_basic_operations.md#7-package-operations)
- REQ-API_BASIC-052: Validate behavior defines validation process. [api_basic_operations.md#711-validate-behavior](../tech_specs/api_basic_operations.md#711-validate-behavior)
- REQ-API_BASIC-053: Validate error conditions define validation errors. [api_basic_operations.md#712-validate-error-conditions](../tech_specs/api_basic_operations.md#712-validate-error-conditions)
- ~~REQ-API_BASIC-054: Validate example usage demonstrates validation usage~~ [type: documentation-only] (documentation-only - has testable scenarios in `validate_example_usage_validation.feature`). [api_basic_operations.md#713-validate-example-usage](../tech_specs/api_basic_operations.md#713-validate-example-usage)
- REQ-API_BASIC-055: Defragment behavior defines defragmentation process. [api_basic_operations.md#721-defragment-behavior](../tech_specs/api_basic_operations.md#721-defragment-behavior)
- REQ-API_BASIC-056: Defragment error conditions define defragmentation errors. [api_basic_operations.md#722-defragment-error-conditions](../tech_specs/api_basic_operations.md#722-defragment-error-conditions)
- ~~REQ-API_BASIC-057: Defragment example usage demonstrates defragmentation usage~~ [type: documentation-only] (documentation-only - has testable scenarios in `defragment_example_usage.feature`). [api_basic_operations.md#723-defragment-example-usage](../tech_specs/api_basic_operations.md#723-defragment-example-usage)
- REQ-API_BASIC-058: GetInfo error conditions define information retrieval errors. [api_basic_operations.md#731-getinfo-error-conditions](../tech_specs/api_basic_operations.md#731-getinfo-error-conditions)
- ~~REQ-API_BASIC-059: GetInfo example usage demonstrates information retrieval usage~~ [type: documentation-only] (documentation-only - has testable scenarios in `getinfo_example_usage_information.feature`). [api_basic_operations.md#732-getinfo-example-usage](../tech_specs/api_basic_operations.md#732-getinfo-example-usage)
- REQ-API_BASIC-060: ReadHeader use cases define header reading scenarios. [api_basic_operations.md#741-readheader-use-cases](../tech_specs/api_basic_operations.md#741-readheader-use-cases)
- REQ-API_BASIC-061: ReadHeader parameters define header reading parameters. [api_basic_operations.md#742-readheader-parameters](../tech_specs/api_basic_operations.md#742-readheader-parameters)
- REQ-API_BASIC-062: ReadHeader error conditions define header reading errors. [api_basic_operations.md#743-readheader-error-conditions](../tech_specs/api_basic_operations.md#743-readheader-error-conditions)
- ~~REQ-API_BASIC-063: ReadHeader example usage demonstrates header reading usage~~ [type: documentation-only] (documentation-only - has testable scenarios in `readheader_example_usage_reading.feature`). [api_basic_operations.md#744-readheader-example-usage](../tech_specs/api_basic_operations.md#744-readheader-example-usage)

## Error Handling

- REQ-API_BASIC-016: Error handling returns structured errors for all failure cases. [api_basic_operations.md#8-error-handling](../tech_specs/api_basic_operations.md#8-error-handling)
- REQ-API_BASIC-018: Path parameters validated (non-empty, normalized, valid format) [type: constraint]. [api_basic_operations.md#8-error-handling](../tech_specs/api_basic_operations.md#8-error-handling)
- REQ-API_BASIC-019: Context errors returned when context is cancelled or times out. [api_basic_operations.md#8-error-handling](../tech_specs/api_basic_operations.md#8-error-handling)
- REQ-API_BASIC-064: Structured error system defines error handling system [type: architectural]. [api_basic_operations.md#81-structured-error-system](../tech_specs/api_basic_operations.md#81-structured-error-system)
- REQ-API_BASIC-065: Common error types define standard error classifications. [api_basic_operations.md#82-common-error-types](../tech_specs/api_basic_operations.md#82-common-error-types)
- REQ-API_BASIC-066: Error types used define specific error types. [api_basic_operations.md#821-error-types-used](../tech_specs/api_basic_operations.md#821-error-types-used)
- REQ-API_BASIC-067: Structured error system enables error type checking, message retrieval, context access, and type-safe error handling. [api_basic_operations.md#83-structured-error-examples](../tech_specs/api_basic_operations.md#83-structured-error-examples)
- REQ-API_BASIC-068: Creating structured errors supports error creation. [api_basic_operations.md#831-creating-structured-errors](../tech_specs/api_basic_operations.md#831-creating-structured-errors)
- REQ-API_BASIC-069: Error inspection provides error examination methods. [api_basic_operations.md#832-error-inspection](../tech_specs/api_basic_operations.md#832-error-inspection)
- REQ-API_BASIC-070: Common error scenarios demonstrate typical error cases. [api_basic_operations.md#833-common-error-scenarios](../tech_specs/api_basic_operations.md#833-common-error-scenarios)
- REQ-API_BASIC-071: Error handling best practices define recommended patterns [type: documentation-only]. [api_basic_operations.md#84-error-handling-best-practices](../tech_specs/api_basic_operations.md#84-error-handling-best-practices)
- REQ-API_BASIC-072: Always check for errors defines error checking requirement. [api_basic_operations.md#841-always-check-for-errors](../tech_specs/api_basic_operations.md#841-always-check-for-errors)
- REQ-API_BASIC-073: Use structured errors for better debugging defines error reporting. [api_basic_operations.md#842-use-structured-errors-for-better-debugging](../tech_specs/api_basic_operations.md#842-use-structured-errors-for-better-debugging)
- REQ-API_BASIC-074: Use context for cancellation defines cancellation handling. [api_basic_operations.md#843-use-context-for-cancellation](../tech_specs/api_basic_operations.md#843-use-context-for-cancellation)
- REQ-API_BASIC-075: Handle different error types appropriately defines error handling. [api_basic_operations.md#844-handle-different-error-types-appropriately](../tech_specs/api_basic_operations.md#844-handle-different-error-types-appropriately)
- REQ-API_BASIC-076: Clean up resources defines resource management requirement. [api_basic_operations.md#845-clean-up-resources](../tech_specs/api_basic_operations.md#845-clean-up-resources)
- REQ-API_BASIC-082: Error handling best practices define error patterns [type: documentation-only]. [api_basic_operations.md#92-error-handling](../tech_specs/api_basic_operations.md#92-error-handling)
- REQ-API_BASIC-083: Wrap errors with context defines error wrapping pattern. [api_basic_operations.md#921-wrap-errors-with-context](../tech_specs/api_basic_operations.md#921-wrap-errors-with-context)
- REQ-API_BASIC-084: Handle specific error types defines targeted error handling. [api_basic_operations.md#922-handle-specific-error-types](../tech_specs/api_basic_operations.md#922-handle-specific-error-types)

## Context Integration

- REQ-API_BASIC-017: All methods accept context.Context and respect cancellation/timeout [type: constraint]. [api_basic_operations.md#02-context-integration](../tech_specs/api_basic_operations.md#02-context-integration)

## Best Practices

- REQ-API_BASIC-077: Best practices document recommended usage patterns [type: documentation-only]. [api_basic_operations.md#9-best-practices](../tech_specs/api_basic_operations.md#9-best-practices)
- REQ-API_BASIC-078: Package lifecycle management defines lifecycle patterns [type: documentation-only]. [api_basic_operations.md#91-package-lifecycle-management](../tech_specs/api_basic_operations.md#91-package-lifecycle-management)
- REQ-API_BASIC-079: Always use defer for cleanup defines cleanup pattern [type: documentation-only]. [api_basic_operations.md#911-always-use-defer-for-cleanup](../tech_specs/api_basic_operations.md#911-always-use-defer-for-cleanup)
- REQ-API_BASIC-080: Check package state before operations defines state checking requirement. [api_basic_operations.md#912-check-package-state-before-operations](../tech_specs/api_basic_operations.md#912-check-package-state-before-operations)
- REQ-API_BASIC-081: Use appropriate context timeouts defines timeout configuration [type: documentation-only]. [api_basic_operations.md#913-use-appropriate-context-timeouts](../tech_specs/api_basic_operations.md#913-use-appropriate-context-timeouts)
- REQ-API_BASIC-085: Resource management defines resource handling patterns [type: documentation-only]. [api_basic_operations.md#93-resource-management](../tech_specs/api_basic_operations.md#93-resource-management)
- REQ-API_BASIC-086: Use context for resource management defines context-based management [type: documentation-only]. [api_basic_operations.md#931-use-context-for-resource-management](../tech_specs/api_basic_operations.md#931-use-context-for-resource-management)
- REQ-API_BASIC-087: Handle cleanup errors gracefully defines cleanup error handling [type: documentation-only]. [api_basic_operations.md#932-handle-cleanup-errors-gracefully](../tech_specs/api_basic_operations.md#932-handle-cleanup-errors-gracefully)
