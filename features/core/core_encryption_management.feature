@domain:core @m2 @REQ-CORE-033 @spec(api_core.md#5-encryption-management)
Feature: Core Encryption Management

  @REQ-CORE-033 @happy
  Scenario: Encryption management defines encryption capabilities
    Given an open NovusPack package
    And a valid context
    When encryption management is used
    Then encryption capabilities are available
    And file encryption is supported
    And encryption operations are accessible

  @REQ-CORE-033 @happy
  Scenario: Encryption management integrates with core package interface
    Given an open NovusPack package
    And encryption operations are needed
    When encryption management is used
    Then encryption integrates with core interface
    And encryption follows core patterns
    And encryption uses structured error system

  @REQ-CORE-033 @happy
  Scenario: Encryption management provides ML-KEM key management
    Given an open NovusPack package
    And a valid context
    When ML-KEM key management is accessed
    Then ML-KEM key generation is available
    And post-quantum encryption keys can be managed
    And key operations support security levels 1-5

  @REQ-CORE-033 @happy
  Scenario: Encryption management supports key operations
    Given an open NovusPack package
    And a valid context
    And an ML-KEM key is available
    When key operations are performed
    Then encrypt operation is supported
    And decrypt operation is supported
    And key lifecycle management is available

  @REQ-CORE-033 @happy
  Scenario: Encryption management supports multiple security levels
    Given an open NovusPack package
    And a valid context
    When encryption keys are generated
    Then security levels 1-5 are supported
    And each level provides appropriate security guarantees
    And higher levels provide increased security

  @REQ-CORE-033 @error
  Scenario: Encryption management handles invalid security levels
    Given an open NovusPack package
    And a valid context
    When an invalid security level is requested
    Then a structured error is returned
    And error indicates invalid security level
    And error follows structured error format

  @REQ-CORE-033 @happy
  Scenario: Encryption management respects context cancellation
    Given an open NovusPack package
    And a context that can be cancelled
    When encryption operation is performed
    And context is cancelled
    Then operation is cancelled gracefully
    And structured context error is returned
