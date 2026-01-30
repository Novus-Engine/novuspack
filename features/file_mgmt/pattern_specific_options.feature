@domain:file_mgmt @m2 @REQ-FILEMGMT-070 @spec(api_file_mgmt_addition.md#2523-pattern-specific-options)
Feature: Pattern-Specific Options

  @REQ-FILEMGMT-070 @happy
  Scenario: Pattern-specific options control exclude patterns
    Given an open NovusPack package
    And a file pattern to match
    And a valid context
    When AddFileOptions with ExcludePatterns is set
    Then exclude patterns are configured
    And ExcludePatterns option contains patterns to exclude
    And matching files are excluded from processing

  @REQ-FILEMGMT-070 @happy
  Scenario: Pattern-specific options control maximum file size
    Given an open NovusPack package
    And a file pattern to match
    And a valid context
    When AddFileOptions with MaxFileSize is set
    Then maximum file size is configured
    And MaxFileSize option limits file size inclusion
    And files exceeding MaxFileSize are excluded
    And zero value means no limit

  @REQ-FILEMGMT-070 @happy
  Scenario: Pattern-specific options control path preservation
    Given an open NovusPack package
    And a file pattern to match
    And a valid context
    When AddFileOptions with PreservePaths is set
    Then path preservation is configured
    And PreservePaths option preserves directory structure
    And directory structure is maintained in package

  @REQ-FILEMGMT-070 @happy
  Scenario: Pattern-specific options have default values
    Given an open NovusPack package
    And a file pattern to match
    And a valid context
    When AddFileOptions with nil or default values is used
    Then exclude patterns default to nil
    And maximum file size defaults to 0 (no limit)
    And path preservation defaults to false

  @REQ-FILEMGMT-070 @happy
  Scenario: Pattern-specific options only apply to pattern operations
    Given an open NovusPack package
    And a single file to add
    And a valid context
    When AddFileOptions with pattern-specific options is used with AddFile
    Then pattern-specific options are ignored
    And only file processing and encryption options apply
    And single file addition proceeds normally

  @REQ-FILEMGMT-070 @error
  Scenario: Pattern-specific options validate exclude patterns
    Given an open NovusPack package
    And a file pattern to match
    And invalid exclude patterns
    And a valid context
    When AddFileOptions with invalid exclude patterns is used
    Then a structured error is returned
    And error indicates invalid exclude patterns
    And error follows structured error format
