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
    And file information includes PrimaryPath, Paths, FileID, and metadata fields
    And results are sorted by PrimaryPath alphabetically
    And results are stable across calls when package state has not changed

  @happy
  Scenario: GetInfo returns lightweight package information
    Given an open NovusPack package
    When GetInfo is called
    Then PackageInfo structure is returned
    And package information includes header-derived fields
    And package information includes computed package-level stats
    And package information includes signature summary
    And package information does not include individual FileEntry metadata
    And package information does not include special metadata file contents

  @happy
  Scenario: GetMetadata returns comprehensive package metadata
    Given an open NovusPack package
    When GetMetadata is called
    Then PackageMetadata structure is returned
    And metadata includes all fields from GetInfo
    And metadata includes FileEntry metadata for all files
    And metadata includes special metadata files contents
    And metadata includes path metadata entries with inheritance chains

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

  @happy
  Scenario: ListFiles does not require context parameter
    Given an open NovusPack package
    When ListFiles is called without context
    Then file information is returned successfully

  @happy
  Scenario: GetInfo does not require context parameter
    Given an open NovusPack package
    When GetInfo is called without context
    Then package information is returned successfully

  @happy
  Scenario: GetMetadata does not require context parameter
    Given an open NovusPack package
    When GetMetadata is called without context
    Then comprehensive metadata is returned successfully

  @error
  Scenario: ReadFile fails for invalid package path
    Given an open NovusPack package
    When ReadFile is called with empty path
    Then a structured validation error is returned
    And error type is ErrTypeValidation

  @error
  Scenario: ReadFile fails for path that escapes package root
    Given an open NovusPack package
    When ReadFile is called with path "../../etc/passwd"
    Then a structured validation error is returned
    And error type is ErrTypeValidation
