@domain:file_mgmt @m2 @REQ-FILEMGMT-026 @spec(api_file_mgmt_addition.md#21-addfile-package-method)
Feature: Filesystem path reading and streaming

  @REQ-FILEMGMT-026 @happy
  Scenario: AddFile supports streaming for large files
    Given an open NovusPack package
    And a valid context
    And a filesystem file path
    And a large file to process
    When AddFile is called
    Then streaming is used for large files when needed
    And memory usage is controlled

  @REQ-FILEMGMT-026 @error
  Scenario: AddFile returns structured I/O error for inaccessible path
    Given an open NovusPack package
    And a valid context
    And an inaccessible file path
    When AddFile is called
    Then structured I/O error is returned
    And error follows structured error format
