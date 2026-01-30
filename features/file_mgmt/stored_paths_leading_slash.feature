@domain:file_mgmt @m2 @REQ-FILEMGMT-301 @spec(api_file_mgmt_file_entry.md#5-path-management)
Feature: All stored paths MUST have leading slash for full path references

  @REQ-FILEMGMT-301 @happy
  Scenario: Stored paths have leading slash
    Given file entries with stored paths
    When paths are stored or retrieved
    Then all stored paths have leading slash for full path references
    And the behavior matches the path-management specification
    And path normalization enforces leading slash
    And display paths strip leading slash for user output
