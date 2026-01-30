@domain:file_mgmt @m2 @REQ-FILEMGMT-248 @spec(api_file_mgmt_file_entry.md#32-tag-type)
Feature: Tag type provides type-safe tag operations with Tag[T] structure

  @REQ-FILEMGMT-248 @happy
  Scenario: Tag type provides type-safe tag operations
    Given a FileEntry or path metadata with tags
    When tag operations are performed using Tag structure
    Then type-safe tag operations are provided as specified
    And Tag structure enforces type safety
    And the behavior matches the tag-type specification
    And tag values are get and set with type safety
