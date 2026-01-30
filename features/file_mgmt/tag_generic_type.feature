@domain:file_mgmt @m2 @REQ-FILEMGMT-422 @spec(api_file_mgmt_file_entry.md#19-tag-generic-type)
Feature: Tag generic type represents a type-safe tag with a specific value type

  @REQ-FILEMGMT-422 @happy
  Scenario: Tag generic type represents type-safe tag
    Given tag metadata for file entries
    When Tag generic type is used for tag representation
    Then Tag generic type represents a type-safe tag with a specific value type as specified
    And the behavior matches the Tag generic type specification
    And Tag[T] enforces value type at compile time
    And GetValue and SetValue are type-safe
