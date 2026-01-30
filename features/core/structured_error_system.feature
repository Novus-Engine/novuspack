@domain:core @m2 @REQ-CORE-062 @spec(api_core.md#105-structured-error-system)
Feature: Structured Error System
    As a package API consumer
    I want consistent error handling across all package operations
    So that errors are predictable, inspectable, and actionable

    Background:
        Given the structured error system is used for all package errors
        And errors use the PackageError structure with typed contexts

    @REQ-CORE-062 @happy
    Scenario: All package operations return structured errors
        When any package operation fails
        Then it MUST return a PackageError or nil
        And errors MUST include error type, message, and optional context

    @REQ-CORE-062 @happy
    Scenario: Error types categorize failure modes
        When a structured error is created
        Then it MUST have an ErrorType from the defined error type categories
        And error types include Validation, IO, Security, Corruption, Unsupported, Context, Encryption, Compression, Signature, State, Format, Version, Resource, Metadata

    @REQ-CORE-062 @happy
    Scenario: Errors include typed context for debugging
        When a structured error is created
        Then it MAY include a typed context with operation-specific details
        And context types include FileErrorContext, EncryptionErrorContext, etc.

    @REQ-CORE-062 @happy
    Scenario: Errors support cause chain wrapping
        When an error wraps an underlying error
        Then the cause chain MUST be preserved
        And errors MUST be unwrappable using standard Go error patterns

    @REQ-CORE-062 @happy
    Scenario: Error inspection uses helper functions
        When checking for specific error conditions
        Then helper functions MUST enable type-safe error inspection
        And functions include AsPackageError, GetErrorContext, HasErrorType

    @REQ-CORE-062 @happy
    Scenario: Error system is consistent across all APIs
        When errors are returned from different package operations
        Then error structure and inspection patterns MUST be identical
        And consumers MUST be able to handle errors uniformly

    @REQ-CORE-062 @happy
    Scenario: PackageError structure contains required fields
        When a PackageError is created
        Then it MUST have a Type field with ErrorType value
        And it MUST have a Message field with human-readable description
        And it MAY have a Cause field wrapping an underlying error
        And it MAY have a Context map for additional debugging information

    @REQ-CORE-062 @happy
    Scenario: Error method formats error messages correctly
        Given a PackageError with message and cause
        When Error method is called
        Then it MUST return a formatted error string
        And format MUST be "{Message}: {Cause}" if cause exists
        And format MUST be "{Message}" if no cause exists
        And the method MUST implement the error interface

    @REQ-CORE-062 @happy
    Scenario: Is method enables error matching and comparison
        Given a PackageError instance
        When Is method is called with a target error
        Then it MUST return true if errors match
        And it MUST support Go standard error comparison patterns
        And it MUST enable errors.Is compatibility
        And error matching MUST work with wrapped errors
