@domain:core @m1 @REQ-CORE-004 @spec(api_core.md#11-package-reader-interface)
Feature: PackageReader interface methods

  @happy
  Scenario: ReadFile reads file content from package
    Given an open NovusPack package with files
    When ReadFile is called with a file path
    Then file content is returned as byte slice
    And file content matches stored data
    And decryption is applied if file is encrypted
    And decompression is applied if file is compressed

  @happy
  Scenario: ListFiles returns all file information
    Given an open NovusPack package with multiple files
    When ListFiles is called
    Then a slice of FileInfo is returned
    And all files in package are included
    And file information includes paths and metadata

  @happy
  Scenario: GetMetadata returns package metadata
    Given an open NovusPack package
    When GetMetadata is called
    Then PackageInfo structure is returned
    And package metadata is complete
    And metadata includes comment, VendorID, AppID

  @happy
  Scenario: Validate performs package validation
    Given an open NovusPack package
    When Validate is called via PackageReader interface
    Then comprehensive package validation is performed
    And validation results are returned

  @happy
  Scenario: GetInfo returns package information
    Given an open NovusPack package
    When GetInfo is called via PackageReader interface
    Then PackageInfo is returned
    And package information is complete

  @error
  Scenario: ReadFile fails for non-existent file
    Given an open NovusPack package
    When ReadFile is called with non-existent path
    Then a structured validation error is returned
    And error type is ErrTypeValidation

  @error
  Scenario: ReadFile respects context cancellation
    Given an open NovusPack package
    And a cancelled context
    When ReadFile is called
    Then a structured context error is returned
