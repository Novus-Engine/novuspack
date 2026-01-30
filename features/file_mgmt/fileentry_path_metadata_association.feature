@domain:file_mgmt @m2 @REQ-FILEMGMT-048 @spec(api_file_mgmt_file_entry.md#5-path-management)
Feature: FileEntry path metadata association methods manage path relationships

  @REQ-FILEMGMT-048 @happy
  Scenario: FileEntry path metadata association methods manage path relationships
    Given a FileEntry with one or more paths
    When path metadata association methods are used
    Then path relationships are managed as specified
    And the behavior matches the path management specification
    And path metadata is consistent with the path management contract
