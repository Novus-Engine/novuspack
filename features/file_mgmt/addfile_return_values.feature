@domain:file_mgmt @m2 @REQ-FILEMGMT-052 @spec(api_file_management.md#213-returns)
Feature: AddFile Return Values

  @REQ-FILEMGMT-052 @happy
  Scenario: AddFile returns created FileEntry with metadata
    Given an open NovusPack package
    And a valid context
    And a file to be added
    When AddFile is called
    Then created FileEntry is returned
    And FileEntry contains all metadata
    And FileEntry contains compression status
    And FileEntry contains encryption details
    And FileEntry contains checksums

  @REQ-FILEMGMT-052 @happy
  Scenario: AddFile returns error if file addition fails
    Given an open NovusPack package
    And a valid context
    And a file that fails to be added
    When AddFile is called and fails
    Then error is returned
    And error indicates failure reason
    And error follows structured error format
    And no FileEntry is returned on failure

  @REQ-FILEMGMT-052 @error
  Scenario: AddFile handles context cancellation errors
    Given an open NovusPack package
    And a context that is cancelled
    And a file to be added
    When AddFile is called
    And context is cancelled
    Then structured context error is returned
    And error follows structured error format
