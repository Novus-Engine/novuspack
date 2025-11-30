@domain:basic_ops @m2 @REQ-API_BASIC-050 @spec(api_basic_operations.md#621-closewithcleanup-behavior)
Feature: CloseWithCleanup behavior

  @REQ-API_BASIC-050 @happy
  Scenario: CloseWithCleanup closes package and performs cleanup
    Given an open NovusPack package
    When CloseWithCleanup is called
    Then package file is closed
    And cleanup operations are performed
    And defragmentation occurs if needed
    And optimization operations occur
    And all resources are released

  @REQ-API_BASIC-050 @happy
  Scenario: CloseWithCleanup takes longer than standard Close
    Given an open NovusPack package
    When CloseWithCleanup is called
    Then cleanup operations take additional time
    And defragmentation adds to execution time
    And optimization adds to execution time
    And operation completes successfully

  @REQ-API_BASIC-050 @happy
  Scenario: CloseWithCleanup performs defragmentation during cleanup
    Given an open NovusPack package with deleted files
    When CloseWithCleanup is called
    Then defragmentation is performed
    And unused space is removed
    And package structure is optimized

  @REQ-API_BASIC-050 @error
  Scenario: CloseWithCleanup handles cleanup errors gracefully
    Given an open NovusPack package
    And cleanup operation fails
    When CloseWithCleanup is called
    Then cleanup error is handled
    And package is still closed
    And error provides details about cleanup failure
