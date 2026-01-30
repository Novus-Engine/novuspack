@domain:file_mgmt @m2 @REQ-FILEMGMT-272 @spec(api_file_mgmt_queries.md#12-purpose)
Feature: Purpose defines basic file existence checks and listing operations

  @REQ-FILEMGMT-272 @happy
  Scenario: Purpose defines file existence and listing
    Given an open NovusPack package
    When file existence or listing purpose is applied
    Then basic file existence checks and listing operations are defined
    And the behavior matches the purpose specification
    And operations are pure in-memory over package state
    And results are consistent with package contents
