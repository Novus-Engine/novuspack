@domain:file_mgmt @m2 @REQ-FILEMGMT-416 @spec(api_file_mgmt_file_entry.md#12-hashtype-type)
Feature: HashType type represents hash algorithm types

  @REQ-FILEMGMT-416 @happy
  Scenario: HashType type represents hash algorithm types
    Given hash information for file content
    When HashType is used for hash algorithm identification
    Then HashType type represents hash algorithm types as specified
    And the behavior matches the HashType type specification
    And algorithm types are defined and queryable
    And type safety is preserved
