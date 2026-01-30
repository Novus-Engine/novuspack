@domain:file_mgmt @m2 @REQ-FILEMGMT-417 @spec(api_file_mgmt_file_entry.md#13-hashpurpose-type)
Feature: HashPurpose type represents hash purposes

  @REQ-FILEMGMT-417 @happy
  Scenario: HashPurpose type represents hash purposes
    Given hash information for file content
    When HashPurpose is used for hash purpose identification
    Then HashPurpose type represents hash purposes as specified
    And the behavior matches the HashPurpose type specification
    And purpose values are defined and queryable
    And type safety is preserved
