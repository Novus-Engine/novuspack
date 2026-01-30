@domain:generics @m2 @REQ-GEN-032 @spec(api_generics.md#15-data-structure-operations)
Feature: Data Structure Operations
    As a package developer
    I want generic data structure utilities
    So that common data manipulation patterns are reusable across the API

    Background:
        Given data structure operations are provided in the generics package
        And operations support map, set, and aggregation utilities

    @REQ-GEN-032 @happy
    Scenario: Map operations support key-value transformations
        When map operations are used on data structures
        Then map transformations MUST preserve key-value relationships
        And operations MUST support type-safe transformations

    @REQ-GEN-032 @happy
    Scenario: Set operations support unique value collections
        When set operations are used on data structures
        Then set operations MUST enforce uniqueness
        And operations MUST support union, intersection, and difference

    @REQ-GEN-032 @happy
    Scenario: Aggregation operations support collection reduction
        When aggregation operations are used on collections
        Then operations MUST support reduction patterns (sum, count, min, max)
        And aggregations MUST handle empty collections safely

    @REQ-GEN-032 @happy
    Scenario: Data structure operations preserve type safety
        When generic data structure operations are applied
        Then all operations MUST maintain compile-time type safety
        And runtime type validation MUST prevent type mismatches

    @REQ-GEN-032 @happy
    Scenario: Data structure operations handle nil values
        When data structure operations receive nil or empty inputs
        Then operations MUST handle nil values gracefully
        And operations MUST NOT panic on nil inputs
