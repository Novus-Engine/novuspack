@domain:compression @m2 @REQ-COMPR-113 @spec(api_package_compression.md#422-decompresspackage-parameters)
Feature: DecompressPackage Parameter Specification

  @REQ-COMPR-113 @happy
  Scenario: DecompressPackage accepts context parameter
    Given an open NovusPack package
    And package is compressed
    And a valid context
    When DecompressPackage is called with context
    Then context parameter is accepted
    And context supports cancellation
    And context supports timeout handling

  @REQ-COMPR-113 @happy
  Scenario: DecompressPackage uses context for cancellation handling
    Given an open NovusPack package
    And package is compressed
    And a context with cancellation support
    When DecompressPackage is called
    Then context cancellation is checked during operation
    And operation can be cancelled via context
    And cancellation is respected

  @REQ-COMPR-113 @happy
  Scenario: DecompressPackage uses context for timeout handling
    Given an open NovusPack package
    And package is compressed
    And a context with timeout
    When DecompressPackage is called
    Then context timeout is checked during operation
    And operation is terminated when timeout expires
    And timeout is respected

  @REQ-COMPR-113 @happy
  Scenario: DecompressPackage method signature requires only context parameter
    Given an open NovusPack package
    And package is compressed
    When DecompressPackage method signature is examined
    Then method accepts only context parameter
    And no other parameters are required
    And method signature is simple and clear

  @REQ-COMPR-113 @error
  Scenario: DecompressPackage handles context cancellation errors
    Given an open NovusPack package
    And package is compressed
    And a cancelled context
    When DecompressPackage is called
    Then context cancellation error is returned
    And error type is context cancellation
    And error follows structured error format

  @REQ-COMPR-113 @error
  Scenario: DecompressPackage handles context timeout errors
    Given an open NovusPack package
    And package is compressed
    And a context that times out
    When DecompressPackage is called
    Then context timeout error is returned
    And error type is context timeout
    And error follows structured error format
