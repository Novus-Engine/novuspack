@domain:file_mgmt @m2 @REQ-FILEMGMT-209 @spec(api_file_management.md#9222-parameters)
Feature: GetFileByHash Parameter Specification

  @REQ-FILEMGMT-209 @happy
  Scenario: GetFileByHash parameters include context, hashType, and hashData
    Given an open NovusPack package
    And a valid context
    And a hash type and hash data
    When GetFileByHash is called
    Then context parameter is accepted for cancellation and timeout handling
    And hashType parameter specifies hash algorithm type
    And hashData parameter contains hash value bytes
    And parameters are validated

  @REQ-FILEMGMT-209 @happy
  Scenario: GetFileByHash supports multiple hash types
    Given an open NovusPack package
    And a valid context
    When GetFileByHash is called with different hash types
    Then SHA-256 hash type is supported
    And SHA-512 hash type is supported
    And BLAKE3 hash type is supported
    And XXH3 hash type is supported

  @REQ-FILEMGMT-209 @happy
  Scenario: GetFileByHash context supports cancellation
    Given an open NovusPack package
    And a context that can be cancelled
    And a hash type and hash data
    When GetFileByHash is called
    And context is cancelled
    Then operation respects context cancellation
    And structured context error is returned

  @REQ-FILEMGMT-209 @happy
  Scenario: GetFileByHash context supports timeout handling
    Given an open NovusPack package
    And a context with timeout
    And a hash type and hash data
    When GetFileByHash is called
    And timeout is exceeded
    Then operation respects context timeout
    And structured context timeout error is returned

  @REQ-FILEMGMT-209 @error
  Scenario: GetFileByHash handles package not open errors
    Given a closed NovusPack package
    And a valid context
    And a hash type and hash data
    When GetFileByHash is called
    Then a structured error is returned
    And error indicates package is not open
    And error follows structured error format

  @REQ-FILEMGMT-209 @error
  Scenario: GetFileByHash handles invalid hash type
    Given an open NovusPack package
    And a valid context
    And an invalid hash type
    When GetFileByHash is called with invalid hash type
    Then a structured error is returned
    And error indicates unsupported hash type
    And error follows structured error format
