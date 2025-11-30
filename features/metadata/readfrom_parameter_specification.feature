@domain:metadata @m2 @REQ-META-057 @spec(api_metadata.md#112-readfrom-parameters)
Feature: ReadFrom Parameter Specification

  @REQ-META-057 @happy
  Scenario: ReadFrom reads package comment from io.Reader
    Given a NovusPack package
    And a PackageComment
    And an io.Reader with comment data
    When ReadFrom is called with reader parameter
    Then comment data is read from reader
    And number of bytes read is returned
    And comment is populated with reader data

  @REQ-META-057 @happy
  Scenario: ReadFrom parameter r is io.Reader
    Given a NovusPack package
    And a PackageComment
    When ReadFrom is called
    Then r parameter accepts io.Reader interface
    And reader can be any io.Reader implementation
    And data is read from reader source

  @REQ-META-057 @happy
  Scenario: ReadFrom returns number of bytes read
    Given a NovusPack package
    And a PackageComment
    And an io.Reader with comment data
    When ReadFrom is called
    Then function returns int64 indicating bytes read
    And bytes read count matches actual data read
    And return value can be used for verification

  @REQ-META-057 @error
  Scenario: ReadFrom handles I/O errors
    Given a NovusPack package
    And a PackageComment
    And an io.Reader that fails
    When ReadFrom is called
    Then I/O errors are returned
    And ErrIOError is returned for I/O problems
    And error follows structured error format
