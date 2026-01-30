@domain:file_mgmt @m2 @REQ-FILEMGMT-318 @spec(api_basic_operations.md#76-session-base-management) @spec(api_file_mgmt_addition.md#overriding-session-base)
Feature: Session Base can be explicitly set, cleared, and queried via Package methods

  @REQ-FILEMGMT-318 @happy
  Scenario: Session Base can be set, cleared, and queried
    Given an open Package
    When Session Base methods are used
    Then Session Base can be explicitly set, cleared, and queried via Package methods
    And the behavior matches the session-base-management specification
    And SetSessionBase, ClearSessionBase, GetSessionBase are available
    And query returns current base or empty when not set
