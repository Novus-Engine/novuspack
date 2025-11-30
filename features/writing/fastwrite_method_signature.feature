@domain:writing @m2 @REQ-WRITE-017 @spec(api_writing.md#21-fastwrite-method-signature)
Feature: FastWrite Method Signature

  @REQ-WRITE-017 @happy
  Scenario: FastWrite method signature accepts context, path, and compressionType
    Given an open NovusPack package
    When FastWrite is called with context.Context, path string, and compressionType uint8
    Then method signature matches specification
    And context supports cancellation and timeout
    And path parameter specifies target file location
    And compressionType parameter specifies compression type (0-3)

  @REQ-WRITE-017 @happy
  Scenario: FastWrite method signature accepts context.Context parameter
    Given an open NovusPack package
    When FastWrite method is examined
    Then method accepts context.Context as first parameter
    And context supports cancellation during operation
    And context supports timeout handling
    And context integration is implemented

  @REQ-WRITE-017 @happy
  Scenario: FastWrite method signature accepts path string parameter
    Given an open NovusPack package
    When FastWrite method is examined
    Then method accepts path string as second parameter
    And path specifies target package file location
    And path parameter is validated (non-empty, normalized, writable)

  @REQ-WRITE-017 @happy
  Scenario: FastWrite method signature accepts compressionType uint8 parameter
    Given an open NovusPack package
    When FastWrite method is examined
    Then method accepts compressionType uint8 as third parameter
    And compressionType 0 indicates no compression
    And compressionType 1-3 indicates specific compression types
    And compressionType parameter is validated

  @REQ-WRITE-017 @happy
  Scenario: FastWrite method signature returns error type
    Given an open NovusPack package
    When FastWrite method is examined
    Then method returns error type
    And error indicates operation result
    And nil indicates successful operation
    And non-nil error indicates failure

  @REQ-WRITE-017 @error
  Scenario: FastWrite method signature validates path parameter
    Given an open NovusPack package
    And path parameter is empty
    When FastWrite is called
    Then validation error is returned
    And error indicates path must be non-empty
    And error follows structured error format

  @REQ-WRITE-017 @error
  Scenario: FastWrite method signature validates compressionType parameter
    Given an open NovusPack package
    And compressionType parameter is unsupported (>3)
    When FastWrite is called
    Then validation error is returned
    And error indicates unsupported compression type
    And error follows structured error format
