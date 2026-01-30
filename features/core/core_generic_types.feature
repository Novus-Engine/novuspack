@domain:core @m2 @REQ-CORE-063 @spec(api_core.md#111-core-generic-types) @spec(api_generics.md#1-core-generic-types)
Feature: Core Generic Types
    As a package API developer
    I want foundational generic type abstractions
    So that type-safe patterns are used consistently throughout the API

    Background:
        Given core generic types provide foundational abstractions
        And types include Option, Result, PathEntry, and others

    @REQ-CORE-063 @happy
    Scenario: Option type provides safe optional value handling
        When an optional value is represented
        Then the Option[T] type MUST be used
        And Option provides Some and None variants with type safety

    @REQ-CORE-063 @happy
    Scenario: Result type provides type-safe error handling
        When an operation may fail
        Then the Result[T] type MAY be used for functional error handling
        And Result provides Ok and Err variants

    @REQ-CORE-063 @happy
    Scenario: PathEntry type provides consistent path representation
        When file or directory paths are stored
        Then the PathEntry type MUST be used
        And PathEntry enforces leading slash and path validation rules

    @REQ-CORE-063 @happy
    Scenario: Generic types support compile-time type safety
        When generic types are used in APIs
        Then type constraints MUST be enforced at compile time
        And type mismatches MUST be caught before runtime

    @REQ-CORE-063 @happy
    Scenario: Generic types integrate with error system
        When generic operations fail
        Then failures MUST use the structured error system
        And error types MUST align with PackageError categories

    @REQ-CORE-063 @happy
    Scenario: Core generic types are reused across all APIs
        When new API methods are added
        Then they SHOULD reuse core generic types
        And new generic types SHOULD only be added when necessary
