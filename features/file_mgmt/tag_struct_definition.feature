@domain:file_mgmt @m2 @REQ-FILEMGMT-423 @spec(api_file_mgmt_file_entry.md#191-tag-struct)
Feature: Tag type definition defines generic Tag structure

  @REQ-FILEMGMT-423 @happy
  Scenario: Tag type definition defines Tag structure
    Given tag type definitions for file entries
    When Tag struct is used for tag representation
    Then Tag type definition defines generic Tag structure as specified
    And the behavior matches the Tag struct specification
    And struct fields are Key and Value
    And type safety is preserved
