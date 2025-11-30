@domain:basic_ops @m2 @REQ-API_BASIC-060 @spec(api_basic_operations.md#741-readheader-use-cases)
Feature: ReadHeader Use Cases

  @REQ-API_BASIC-060 @happy
  Scenario: ReadHeader reads package header without opening package
    Given a NovusPack package file on disk
    And a reader for the package file
    When ReadHeader is called with the reader
    Then package header is read successfully
    And header contains magic number
    And header contains format version
    And header contains package metadata
    And package is not opened

  @REQ-API_BASIC-060 @happy
  Scenario: ReadHeader allows inspection before opening
    Given a NovusPack package file
    When ReadHeader is used to inspect header
    Then package format can be validated
    And package version can be checked
    And package metadata can be examined
    And package file is not locked
    And package can still be opened normally

  @REQ-API_BASIC-061 @happy
  Scenario: ReadHeader accepts reader and context parameters
    Given a NovusPack package file
    And a context for cancellation
    And a reader for the package file
    When ReadHeader is called with context and reader
    Then header is read from reader
    And context cancellation is respected
    And reader position is managed correctly

  @REQ-API_BASIC-062 @error
  Scenario: ReadHeader returns error for invalid reader
    Given an invalid or nil reader
    When ReadHeader is called
    Then validation error is returned
    And error indicates invalid reader

  @REQ-API_BASIC-062 @error
  Scenario: ReadHeader returns error for corrupted header
    Given a package file with corrupted header
    When ReadHeader is called
    Then validation error is returned
    And error indicates header corruption
    And error provides details about corruption

  @REQ-API_BASIC-062 @error
  Scenario: ReadHeader respects context cancellation
    Given a package file reader
    And a cancelled context
    When ReadHeader is called with cancelled context
    Then context error is returned
    And error type is context cancellation
