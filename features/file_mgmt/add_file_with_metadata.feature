@domain:file_mgmt @m2 @spec(api_file_management.md#21-add-file)
Feature: Add file with metadata

  @REQ-FILEMGMT-001 @happy
  Scenario: Adding a file updates index and metadata
    Given an open package
    When I add a file "hello.txt" with size 5 and metadata
    Then the index should include "hello.txt" with size 5 and metadata present

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-038 @error
  Scenario: AddFile validates file path parameter
    Given an open writable package
    When AddFile is called with empty path
    Then structured validation error is returned
    And error indicates invalid path

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-038 @error
  Scenario: AddFile rejects whitespace-only paths
    Given an open writable package
    When AddFile is called with whitespace-only path
    Then structured validation error is returned
    And error indicates invalid path

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-039 @error
  Scenario: AddFile validates file data parameter
    Given an open writable package
    When AddFile is called with nil data
    Then structured validation error is returned
    And error indicates invalid data

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-041 @error
  Scenario: AddFile respects context cancellation
    Given an open writable package
    And a cancelled context
    When AddFile is called
    Then structured context error is returned
    And error type is context cancellation

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-040 @error
  Scenario: AddFile validates encryption type parameter
    Given an open writable package
    When AddFileWithEncryption is called with invalid encryption type
    Then structured validation error is returned
    And error indicates unsupported encryption type
