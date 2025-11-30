@domain:file_mgmt @m2 @REQ-FILEMGMT-132 @spec(api_file_management.md#1331-use-patterns-for-bulk-operations)
Feature: Use Patterns for Bulk Operations

  @REQ-FILEMGMT-132 @happy
  Scenario: Patterns are used for bulk operations instead of individual calls
    Given an open NovusPack package
    And a valid context
    And multiple files to add matching a pattern
    When AddFilePattern is used for bulk operations
    Then bulk operations are performed efficiently
    And performance is better than individual AddFile calls
    And multiple files are added in single operation

  @REQ-FILEMGMT-132 @happy
  Scenario: Patterns improve performance for bulk file operations
    Given an open NovusPack package
    And a valid context
    And many files matching a pattern
    When AddFilePattern is used instead of multiple AddFile calls
    Then operation completes faster
    And resource usage is optimized
    And bulk processing is efficient

  @REQ-FILEMGMT-132 @happy
  Scenario: Patterns support bulk file removal operations
    Given an open NovusPack package
    And a valid context
    And multiple files matching a pattern
    When RemoveFilePattern is used for bulk removal
    Then bulk removal is performed efficiently
    And multiple files are removed in single operation
    And performance is optimized
