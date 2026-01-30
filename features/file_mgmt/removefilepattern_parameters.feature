@domain:file_mgmt @m2 @REQ-FILEMGMT-326 @spec(api_file_mgmt_removal.md#33-removefilepattern-parameters)
Feature: RemoveFilePattern parameters include context and pattern

  @REQ-FILEMGMT-326 @happy
  Scenario: RemoveFilePattern parameters include context and pattern
    Given an open NovusPack package with file entries
    When RemoveFilePattern is invoked with parameters
    Then parameters include context and pattern as specified
    And the behavior matches the RemoveFilePattern parameters specification
    And context is used for cancellation
    And pattern is validated before removal
