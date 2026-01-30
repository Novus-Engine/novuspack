@domain:file_mgmt @m2 @REQ-FILEMGMT-309 @spec(api_file_mgmt_addition.md#2821-conflict-handling-options)
Feature: AddFileOptions conflict handling defines overwrite behavior

  @REQ-FILEMGMT-309 @happy
  Scenario: AddFileOptions conflict handling defines overwrite
    Given AddFileOptions with conflict handling configuration
    When a file is added and path already exists
    Then conflict handling defines overwrite behavior as specified
    And the behavior matches the conflict-handling-options specification
    And overwrite or skip is applied per options
    And error or success is returned consistently
