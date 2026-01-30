@domain:file_mgmt @m2 @REQ-FILEMGMT-275 @spec(api_file_mgmt_queries.md#2-single-entry-lookups)
Feature: Single-entry lookups define methods to retrieve one FileEntry by identifier

  @REQ-FILEMGMT-275 @happy
  Scenario: Single-entry lookups retrieve one FileEntry
    Given an open NovusPack package with file entries
    When single-entry lookup methods are used
    Then methods to retrieve one FileEntry by identifier are defined
    And the behavior matches the single-entry-lookups specification
    And GetFileByPath, GetFileByOffset, GetFileByFileID are available
    And nil or error is returned when not found
