@domain:file_mgmt @m2 @REQ-FILEMGMT-264 @spec(api_file_mgmt_queries.md#423-listcompressedfiles-parameters)
Feature: ListCompressedFiles parameters define no input parameters

  @REQ-FILEMGMT-264 @happy
  Scenario: ListCompressedFiles has no input parameters
    Given an open NovusPack package
    When ListCompressedFiles is invoked
    Then parameters define no input parameters as specified
    And the behavior matches the parameters specification
    And method is pure in-memory over package state
    And context may be used for cancellation
