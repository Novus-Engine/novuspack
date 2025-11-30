@domain:compression @m2 @REQ-COMPR-094 @spec(api_package_compression.md#1421-compression-error-types)
Feature: Compression Error Types

  @REQ-COMPR-094 @error
  Scenario: ErrTypeCompression indicates compression and decompression operation failures
    Given a compression operation
    And compression or decompression operation fails
    When error is returned
    Then ErrTypeCompression error type is used
    And error indicates compression operation failure
    And error is compression-specific

  @REQ-COMPR-094 @error
  Scenario: ErrTypeValidation indicates invalid compression parameters
    Given a compression operation
    And invalid compression parameters are provided
    When error is returned
    Then ErrTypeValidation error type is used
    And error indicates invalid parameters
    And error is validation-specific

  @REQ-COMPR-094 @error
  Scenario: ErrTypeIO indicates I/O errors during compression operations
    Given a compression operation
    And I/O error occurs during compression
    When error is returned
    Then ErrTypeIO error type is used
    And error indicates I/O operation failure
    And error is I/O-specific

  @REQ-COMPR-094 @error
  Scenario: ErrTypeContext indicates context cancellation and timeout errors
    Given a compression operation
    And context is cancelled or timeout occurs
    When error is returned
    Then ErrTypeContext error type is used
    And error indicates context issue
    And error is context-specific

  @REQ-COMPR-094 @error
  Scenario: ErrTypeCorruption indicates corrupted compressed data errors
    Given a compression operation
    And compressed data is corrupted
    When error is returned
    Then ErrTypeCorruption error type is used
    And error indicates data corruption
    And error is corruption-specific

  @REQ-COMPR-094 @error
  Scenario: ErrTypeUnsupported indicates unsupported compression algorithms
    Given a compression operation
    And unsupported compression algorithm is used
    When error is returned
    Then ErrTypeUnsupported error type is used
    And error indicates unsupported feature
    And error is unsupported-specific

  @REQ-COMPR-094 @happy
  Scenario: Compression error types enable appropriate error handling
    Given a compression operation
    When error occurs with specific error type
    Then error type determines handling strategy
    And different error types are handled appropriately
    And error handling is type-aware
