@domain:basic_ops @m1 @REQ-API_BASIC-013 @spec(api_basic_operations.md#74-header-inspection)
Feature: Header inspection without opening package

  @happy
  Scenario: ReadHeader reads package header from reader
    Given a NovusPack package file
    When ReadHeader is called with a reader
    Then header is read successfully
    And Header structure is returned
    And header fields are populated correctly

  @happy
  Scenario: ReadHeader validates header format without loading package
    Given a NovusPack package file
    When ReadHeader is called
    Then header magic number is validated
    And header format version is validated
    And header structure is validated
    And package data is not loaded

  @happy
  Scenario: ReadHeader allows metadata inspection before opening
    Given a NovusPack package file
    When ReadHeader is called
    Then header metadata is accessible
    And VendorID and AppID can be inspected
    And package comment size can be checked
    And signature offset can be inspected
    And package version information is available

  @happy
  Scenario: ReadHeader works with stream processing
    Given a stream containing a NovusPack package header
    When ReadHeader is called with stream reader
    Then header is read from stream
    And only header bytes are consumed
    And remaining stream is available for further processing

  @error
  Scenario: ReadHeader fails if header format is invalid
    Given a file with invalid package header
    When ReadHeader is called
    Then a structured validation error is returned
    And error indicates invalid header format

  @error
  Scenario: ReadHeader fails if package version is unsupported
    Given a NovusPack package file with unsupported version
    When ReadHeader is called
    Then a structured unsupported error is returned
    And error indicates unsupported version

  @error
  Scenario: ReadHeader fails with insufficient data
    Given a file with incomplete header data
    When ReadHeader is called
    Then a structured I/O error is returned
    And error indicates unexpected end of file

  @error
  Scenario: ReadHeader respects context cancellation
    Given a NovusPack package file
    And a cancelled context
    When ReadHeader is called
    Then a structured context error is returned
    And header reading is cancelled
