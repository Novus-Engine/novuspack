@domain:file_mgmt @m2 @REQ-FILEMGMT-415 @spec(api_file_mgmt_file_entry.md#11-hashentry-struct)
Feature: HashEntry type represents a hash with type and purpose

  @REQ-FILEMGMT-415 @happy
  Scenario: HashEntry type represents hash with type and purpose
    Given a FileEntry or content with hash information
    When HashEntry is used for hash representation
    Then HashEntry type represents a hash with type and purpose as specified
    And the behavior matches the HashEntry struct specification
    And hash value, type, and purpose are available
    And type safety is preserved
