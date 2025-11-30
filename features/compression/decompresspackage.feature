@domain:compression @m2 @REQ-COMPR-111 @spec(api_package_compression.md#42-decompresspackage)
Feature: DecompressPackage

  @REQ-COMPR-111 @happy
  Scenario: DecompressPackage decompresses package content in memory
    Given an open NovusPack package
    And package is compressed
    And a valid context
    When DecompressPackage is called
    Then all compressed content is decompressed
    And package compression state is updated in memory
    And package header compression flags are cleared
    And all other package data is preserved

  @REQ-COMPR-111 @happy
  Scenario: DecompressPackage accepts context parameter
    Given an open NovusPack package
    And package is compressed
    And a valid context
    When DecompressPackage is called with context
    Then context is accepted as parameter
    And context supports cancellation
    And context supports timeout handling

  @REQ-COMPR-111 @happy
  Scenario: DecompressPackage operates on in-memory package
    Given an open NovusPack package in memory
    And package is compressed
    And a valid context
    When DecompressPackage is called
    Then decompression occurs in memory
    And no file I/O operations are performed
    And package remains in memory after decompression

  @REQ-COMPR-111 @error
  Scenario: DecompressPackage returns error when package is not compressed
    Given an open NovusPack package
    And package is not compressed
    And a valid context
    When DecompressPackage is called
    Then validation error is returned
    And error indicates package is not compressed
    And error follows structured error format

  @REQ-COMPR-111 @error
  Scenario: DecompressPackage returns error on decompression failure
    Given an open NovusPack package
    And package has corrupted compressed data
    And a valid context
    When DecompressPackage is called
    Then compression error is returned
    And error indicates decompression failure
    And error follows structured error format

  @REQ-COMPR-111 @error
  Scenario: DecompressPackage handles context cancellation
    Given an open NovusPack package
    And package is compressed
    And a cancelled context
    When DecompressPackage is called
    Then context cancellation error is returned
    And error type is context cancellation
