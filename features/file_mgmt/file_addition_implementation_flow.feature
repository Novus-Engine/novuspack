@domain:file_mgmt @m2 @REQ-FILEMGMT-010 @spec(api_file_management.md#31-processing-order-requirements)
Feature: File Addition Processing Order

  @happy
  Scenario: File addition follows required processing sequence
    Given a file to be added
    When AddFile is called
    Then file validation occurs first
    Then file data is loaded
    Then deduplication check occurs
    Then compression is applied if configured
    Then encryption is applied if configured
    Then file entry is created
    Then file entry is written to package
    Then file data is written to package
    Then file index is updated
    And processing order is correct

  @happy
  Scenario: Processing order handles dependencies correctly
    Given a file with compression and encryption
    When file is added
    Then compression occurs before encryption
    Then encryption occurs before writing
    Then writing occurs before index update
    And dependencies are respected

  @happy
  Scenario: Error handling preserves package state on failure
    Given a file addition operation
    When error occurs during processing
    Then partial changes are rolled back
    And package state remains consistent
    And error is returned with context

  @happy
  Scenario: Performance requirements are met
    Given large files or many files
    When files are added
    Then operations complete within reasonable time
    And memory usage is controlled
    And I/O operations are efficient

  @error
  Scenario: Processing failures are handled gracefully
    Given a file addition operation
    When compression fails
    Then error is returned
    And package state is not corrupted
    When encryption fails
    Then error is returned
    And package state is not corrupted

  @error
  Scenario: Processing respects context cancellation
    Given a long-running file addition
    And a cancelled context
    When processing is in progress
    Then operation is cancelled
    And partial state is cleaned up
    And context error is returned
