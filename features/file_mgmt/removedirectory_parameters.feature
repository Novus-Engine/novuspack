@domain:file_mgmt @m2 @REQ-FILEMGMT-332 @spec(api_file_mgmt_removal.md#43-removedirectory-parameters)
Feature: RemoveDirectory parameters include context, directory path, and options

  @REQ-FILEMGMT-332 @happy
  Scenario: RemoveDirectory parameters include context path and options
    Given an open NovusPack package
    When RemoveDirectory is invoked with parameters
    Then parameters include context, directory path, and options as specified
    And the behavior matches the RemoveDirectory parameters specification
    And context is used for cancellation
    And directory path is validated before removal
