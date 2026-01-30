@domain:file_mgmt @m2 @REQ-FILEMGMT-268 @spec(api_file_mgmt_queries.md#1-file-existence-and-listing)
Feature: File existence and listing defines basic file existence checks and listing operations

  @REQ-FILEMGMT-268 @happy
  Scenario: File existence and listing defines operations
    Given an open NovusPack package
    When file existence or listing operations are used
    Then basic file existence checks and listing are defined as specified
    And the behavior matches the file-existence-and-listing specification
    And FileExists and ListFiles are available
    And results are consistent with package state
