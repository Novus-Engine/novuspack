@domain:file_mgmt @m2 @REQ-FILEMGMT-302 @spec(api_file_mgmt_file_entry.md#5-path-management)
Feature: Path normalization adds leading slash to paths without one

  @REQ-FILEMGMT-302 @happy
  Scenario: Path normalization adds leading slash
    Given input paths that may lack leading slash
    When path normalization is applied
    Then leading slash is added to paths without one
    And the behavior matches the path-management specification
    And normalized paths are stored consistently
    And display conversion strips leading slash for user output
