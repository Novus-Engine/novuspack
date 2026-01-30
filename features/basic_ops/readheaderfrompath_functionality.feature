@domain:basic_ops @m1 @REQ-API_BASIC-104 @spec(api_basic_operations.md#743-readheaderfrompath)
Feature: ReadHeaderFromPath functionality

  @REQ-API_BASIC-104 @happy
  Scenario: ReadHeaderFromPath reads header from file path
    Given a NovusPack package file at "/path/to/package.nvpk"
    When ReadHeaderFromPath is called with the file path
    Then header is read successfully
    And PackageHeader structure is returned
    And header fields are populated correctly
    And file is closed automatically

  @REQ-API_BASIC-104 @happy
  Scenario: ReadHeaderFromPath provides automatic file management
    Given a NovusPack package file at "/path/to/package.nvpk"
    When ReadHeaderFromPath is called
    Then file is opened automatically
    And header is read from file
    And file is closed automatically
    And no manual file handle management is required

  @REQ-API_BASIC-104 @happy
  Scenario: ReadHeaderFromPath validates package format
    Given a NovusPack package file
    When ReadHeaderFromPath is called
    Then header magic number is validated
    And header format version is validated
    And header structure is validated

  @REQ-API_BASIC-104 @happy
  Scenario: ReadHeaderFromPath exposes package metadata
    Given a NovusPack package file with metadata
    When ReadHeaderFromPath is called
    Then header metadata is accessible
    And format version is available
    And magic number is available
    And VendorID and AppID can be inspected
    And signature information is accessible

  @REQ-API_BASIC-105 @error
  Scenario: ReadHeaderFromPath fails with invalid path
    Given an invalid file path
    When ReadHeaderFromPath is called with invalid path
    Then a structured validation error is returned
    And error indicates invalid file path

  @REQ-API_BASIC-105 @error
  Scenario: ReadHeaderFromPath fails with nonexistent file
    Given a file path that does not exist
    When ReadHeaderFromPath is called
    Then a structured I/O error is returned
    And error indicates file not found

  @REQ-API_BASIC-105 @error
  Scenario: ReadHeaderFromPath fails with permission denied
    Given a NovusPack package file without read permissions
    When ReadHeaderFromPath is called
    Then a structured security error is returned
    And error indicates permission denied

  @REQ-API_BASIC-105 @error
  Scenario: ReadHeaderFromPath fails with corrupted header
    Given a package file with corrupted header
    When ReadHeaderFromPath is called
    Then a structured validation error is returned
    And error indicates invalid header format

  @REQ-API_BASIC-105 @error
  Scenario: ReadHeaderFromPath fails with unsupported version
    Given a NovusPack package file with unsupported version
    When ReadHeaderFromPath is called
    Then a structured unsupported error is returned
    And error indicates unsupported package version

  @REQ-API_BASIC-105 @error
  Scenario: ReadHeaderFromPath respects context cancellation
    Given a NovusPack package file
    And a cancelled context
    When ReadHeaderFromPath is called with cancelled context
    Then a structured context error is returned
    And operation is cancelled before completion
