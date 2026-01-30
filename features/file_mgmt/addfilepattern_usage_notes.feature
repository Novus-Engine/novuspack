@domain:file_mgmt @m2 @REQ-FILEMGMT-071 @spec(api_file_mgmt_addition.md#24-usage-notes)
Feature: AddFilePattern Usage Notes

  @REQ-FILEMGMT-071 @happy
  Scenario: AddFile reads file data from filesystem paths
    Given an open NovusPack package
    And a valid context
    And a filesystem file path
    When AddFile is called
    Then file content is read from filesystem path
    And file is added to package

  @REQ-FILEMGMT-071 @happy
  Scenario: AddFileOptions configures compression, encryption, and metadata
    Given an open NovusPack package
    And a valid context
    When AddFileOptions is used with AddFile
    Then compression can be configured
    And encryption can be configured
    And metadata tags can be configured
    And file processing behavior is customizable

  @REQ-FILEMGMT-071 @happy
  Scenario: AddFile usage notes document best practices
    Given an open NovusPack package
    And a valid context
    When AddFile is used
    Then best practices are documented
    And configuration options are explained
    And stored path determination rules are documented
