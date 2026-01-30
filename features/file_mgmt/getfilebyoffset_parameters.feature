@domain:file_mgmt @m2 @REQ-FILEMGMT-283 @spec(api_file_mgmt_queries.md#223-getfilebyoffset-parameters)
Feature: GetFileByOffset parameters define offset input

  @REQ-FILEMGMT-283 @happy
  Scenario: GetFileByOffset parameters define offset input
    Given an open NovusPack package
    When GetFileByOffset is invoked with parameters
    Then parameters define offset input as specified
    And the behavior matches the parameters specification
    And offset is validated before lookup
    And context may be used for cancellation
