@domain:signatures @m2 @REQ-SIG-057 @spec(api_signatures.md#532-error-handling-patterns)
Feature: Signature Error Handling Patterns

  @REQ-SIG-057 @happy
  Scenario: Error inspection pattern uses IsPackageError
    Given a NovusPack package
    When signature operation returns error
    Then IsPackageError is used to check for structured errors
    And error type is checked using pkgErr.Type
    And error handling switches on error category
    And error handling is type-safe and programmatic

  @REQ-SIG-057 @happy
  Scenario: Type-safe context access uses GetTypedContext
    Given a NovusPack package
    When signature error with context is handled
    Then GetTypedContext is used to access typed error context
    And context is accessed with type-safe SignatureErrorContext
    And context provides debugging information
    And context enables detailed error analysis

  @REQ-SIG-057 @happy
  Scenario: Error creation pattern uses NewTypedPackageError
    Given a NovusPack package
    When signature error is created
    Then NewTypedPackageError is used with error type
    And appropriate context structure is provided
    And error follows recommended error creation pattern
    And error provides comprehensive error information

  @REQ-SIG-057 @happy
  Scenario: Error handling patterns support modern error handling
    Given a NovusPack package
    When signature operations use error handling patterns
    Then errors are handled using structured error system
    And errors are inspected using IsPackageError
    And errors are analyzed using GetTypedContext
    And errors are created using NewTypedPackageError

  @REQ-SIG-057 @error
  Scenario: Error handling patterns handle error conditions correctly
    Given a NovusPack package
    When error handling patterns encounter errors
    Then error inspection correctly identifies structured errors
    And context access provides debugging information
    And error creation produces valid structured errors
