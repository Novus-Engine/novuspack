@domain:signatures @m2 @v2 @REQ-SIG-059 @spec(api_signatures.md#5331-core-error-handling-functions)
Feature: Core Signature Error Handling Functions

  @REQ-SIG-059 @happy
  Scenario: Core error handling functions provide signature validation utilities
    Given a NovusPack package
    When signature validation is performed
    Then validateSignatureInternal function is available
    And function uses structured error system
    And function provides error handling utilities
    And function follows structured error patterns

  @REQ-SIG-059 @happy
  Scenario: Core error handling functions support generic error helpers
    Given a NovusPack package
    When signature operations use error handling
    Then GetTypedContext helper function is available
    And NewTypedPackageError helper function is available
    And WithTypedContext helper function is available
    And WrapWithContext helper function is available
    And MapError helper function is available

  @REQ-SIG-059 @happy
  Scenario: Core error handling functions follow structured error system
    Given a NovusPack package
    When signature error handling functions are used
    Then functions use structured error system from api_core.md
    And functions provide type-safe error context access
    And functions enable modern error handling patterns

  @REQ-SIG-059 @error
  Scenario: Core error handling functions handle errors correctly
    Given a NovusPack package
    When error handling functions encounter errors
    Then functions return structured errors
    And functions provide error context information
    And functions follow error handling patterns
