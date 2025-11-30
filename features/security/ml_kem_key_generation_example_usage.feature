@domain:security @m2 @REQ-SEC-100 @spec(api_security.md#525-example-usage)
Feature: ML-KEM Key Generation Example Usage

  @REQ-SEC-100 @happy
  Scenario: ML-KEM key generation example demonstrates key creation
    Given an open NovusPack package
    And a valid context
    And security level 3
    When GenerateMLKEMKey is called with level 3
    Then ML-KEM key pair is generated
    And key follows example usage pattern
    And key.Clear is called to clear sensitive data
    And error handling demonstrates proper error checking

  @REQ-SEC-100 @happy
  Scenario: ML-KEM key generation example demonstrates error handling
    Given an open NovusPack package
    And a valid context
    And invalid security level
    When GenerateMLKEMKey example is followed
    Then error is checked properly
    And error message is formatted correctly
    And example demonstrates structured error handling

  @REQ-SEC-100 @happy
  Scenario: ML-KEM key generation example demonstrates context usage
    Given an open NovusPack package
    And a valid context
    And ML-KEM key generation requirements
    When example usage is followed
    Then context.Context is used correctly
    And context cancellation is handled
    And context timeout is handled
