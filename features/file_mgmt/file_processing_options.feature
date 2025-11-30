@domain:file_mgmt @m2 @REQ-FILEMGMT-068 @spec(api_file_management.md#2521-file-processing-options)
Feature: File Processing Options

  @REQ-FILEMGMT-068 @happy
  Scenario: File processing options control compression behavior
    Given an open NovusPack package
    And a file to be added
    And a valid context
    When AddFileOptions with Compress option is set
    Then compression behavior is controlled
    And Compress option enables or disables compression
    And CompressionType option specifies algorithm
    And CompressionLevel option controls compression level

  @REQ-FILEMGMT-068 @happy
  Scenario: File processing options control file type
    Given an open NovusPack package
    And a file to be added
    And a valid context
    When AddFileOptions with FileType option is set
    Then file type identifier is configured
    And FileType option specifies file type
    And default file type is regular file

  @REQ-FILEMGMT-068 @happy
  Scenario: File processing options support per-file tags
    Given an open NovusPack package
    And a file to be added
    And a valid context
    When AddFileOptions with Tags option is set
    Then per-file tags are configured
    And Tags option contains key-value pairs
    And tags are applied to the file entry

  @REQ-FILEMGMT-068 @happy
  Scenario: File processing options have default values
    Given an open NovusPack package
    And a file to be added
    And a valid context
    When AddFileOptions with nil or default values is used
    Then compression defaults to false
    And compression type defaults to none
    And file type defaults to regular file
    And tags default to nil

  @REQ-FILEMGMT-068 @error
  Scenario: File processing options validate compression settings
    Given an open NovusPack package
    And a file to be added
    And an invalid compression type
    And a valid context
    When AddFileOptions with invalid compression type is used
    Then a structured error is returned
    And error indicates unsupported compression type
    And error follows structured error format
