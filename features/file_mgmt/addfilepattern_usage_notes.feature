@domain:file_mgmt @m2 @REQ-FILEMGMT-071 @spec(api_file_management.md#26-usage-notes)
Feature: AddFilePattern Usage Notes

  @REQ-FILEMGMT-071 @happy
  Scenario: AddFile supports various FileSource implementations
    Given an open NovusPack package
    And a valid context
    When AddFile is used with different FileSource types
    Then FilePathSource works with AddFile
    And MemorySource works with AddFile
    And custom FileSource implementations work with AddFile
    And unified interface provides flexibility

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
    And FileSource implementations are documented
