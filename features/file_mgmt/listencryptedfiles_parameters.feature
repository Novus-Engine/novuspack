@domain:file_mgmt @m2 @REQ-FILEMGMT-292 @spec(api_file_mgmt_queries.md#433-listencryptedfiles-parameters)
Feature: ListEncryptedFiles parameters define no input parameters

  @REQ-FILEMGMT-292 @happy
  Scenario: ListEncryptedFiles has no input parameters
    Given an open NovusPack package
    When ListEncryptedFiles is invoked
    Then parameters define no input parameters as specified
    And the behavior matches the parameters specification
    And method is pure in-memory over package state
    And context may be used for cancellation
