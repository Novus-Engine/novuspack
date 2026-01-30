@domain:file_mgmt @m2 @REQ-FILEMGMT-418 @spec(api_file_mgmt_file_entry.md#14-tagvaluetype-type)
Feature: TagValueType type represents the type of a tag value

  @REQ-FILEMGMT-418 @happy
  Scenario: TagValueType type represents tag value type
    Given tag metadata for file entries
    When TagValueType is used for tag value type identification
    Then TagValueType type represents the type of a tag value as specified
    And the behavior matches the TagValueType type specification
    And value types are defined and queryable
    And type safety is preserved
