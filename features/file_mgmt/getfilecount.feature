@domain:file_mgmt @m2 @REQ-FILEMGMT-227 @REQ-FILEMGMT-230 @REQ-FILEMGMT-231 @spec(api_file_mgmt_queries.md#41-getfilecount)
Feature: GetFileCount

  @REQ-FILEMGMT-227 @happy
  Scenario: GetFileCount returns total number of files in package
    Given an open NovusPack package
    And a valid context
    And files exist in the package
    When GetFileCount is called
    Then total number of files is returned
    And integer count is accurate
    And error is nil on success

  @REQ-FILEMGMT-230 @happy
  Scenario: GetFileCount returns integer count and error
    Given an open NovusPack package
    And a valid context
    When GetFileCount is called
    Then integer count is returned
    And error is nil on success
    And count represents total file count

  @REQ-FILEMGMT-231 @happy
  Scenario: GetFileCount supports package statistics
    Given an open NovusPack package
    And a valid context
    When GetFileCount is used for statistics
    Then package statistics can be generated
    And file count information is available
    And package metrics can be calculated

  @REQ-FILEMGMT-231 @happy
  Scenario: GetFileCount supports progress tracking
    Given an open NovusPack package
    And a valid context
    And file operations are being performed
    When GetFileCount is used for progress tracking
    Then file count can be tracked during operations
    And progress can be calculated
    And operation status can be determined

  @REQ-FILEMGMT-231 @happy
  Scenario: GetFileCount supports validation and bounds checking
    Given an open NovusPack package
    And a valid context
    When GetFileCount is used for validation
    Then file count can be validated
    And bounds checking can be performed
    And package integrity validation is supported

  @REQ-FILEMGMT-227 @error
  Scenario: GetFileCount handles package not open errors
    Given a closed NovusPack package
    And a valid context
    When GetFileCount is called
    Then a structured error is returned
    And error indicates package is not open
    And error follows structured error format

  @REQ-FILEMGMT-227 @happy
  Scenario: GetFileCount respects context cancellation
    Given an open NovusPack package
    And a context that can be cancelled
    When GetFileCount is called
    And context is cancelled
    Then operation respects context cancellation
    And structured context error is returned
