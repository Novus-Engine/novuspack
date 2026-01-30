@domain:file_mgmt @m2 @REQ-FILEMGMT-148 @spec(api_file_mgmt_updates.md#121-updatefilepattern-purpose)
Feature: UpdateFilePattern purpose is to update multiple files via pattern

  @REQ-FILEMGMT-148 @happy
  Scenario: UpdateFilePattern updates multiple files via pattern
    Given a package and a pattern for file updates
    When UpdateFilePattern is called
    Then the purpose is to update multiple files via pattern
    And the behavior matches the UpdateFilePattern purpose specification
    And pattern matching and bulk updates are performed
