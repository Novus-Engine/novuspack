@domain:compression @m2 @REQ-COMPR-095 @spec(api_package_compression.md#143-structured-error-examples)
Feature: Compression: Structured Error Examples

  @REQ-COMPR-095 @happy
  Scenario: Compression errors are created with rich context
    Given a compression operation that fails
    When structured compression error is created
    Then error includes ErrTypeCompression type
    And error includes algorithm context
    And error includes compression level context
    And error includes input size context
    And error includes operation context

  @REQ-COMPR-095 @happy
  Scenario: Unsupported compression type errors provide details
    Given an unsupported compression type operation
    When structured error is created
    Then error includes ErrTypeUnsupported type
    And error includes compression type context
    And error includes supported types context
    And error includes operation context

  @REQ-COMPR-095 @happy
  Scenario: Memory errors include memory context information
    Given compression operations with insufficient memory
    When structured memory error is created
    Then error includes ErrTypeIO type for memory issues
    And error includes required memory context
    And error includes available memory context
    And error includes algorithm context
    And error includes operation context

  @REQ-COMPR-095 @error
  Scenario: Structured error examples demonstrate error handling patterns
    Given compression operations that may fail
    When structured error examples are applied
    Then error handling patterns are demonstrated
    And error context extraction is shown
    And appropriate error handling responses are illustrated
