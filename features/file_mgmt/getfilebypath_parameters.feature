@domain:file_mgmt @m2 @REQ-FILEMGMT-278 @spec(api_file_mgmt_queries.md#213-getfilebypath-parameters)
Feature: GetFileByPath parameters define virtual path input

  @REQ-FILEMGMT-278 @happy
  Scenario: GetFileByPath parameters define path input
    Given an open NovusPack package
    When GetFileByPath is invoked with parameters
    Then parameters define virtual path input as specified
    And the behavior matches the parameters specification
    And path is normalized with leading slash
    And context may be used for cancellation
