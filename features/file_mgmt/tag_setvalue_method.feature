@domain:file_mgmt @m2 @REQ-FILEMGMT-426 @spec(api_file_mgmt_file_entry.md#194-tagsetvalue-method)
Feature: Tag SetValue method sets the tag value

  @REQ-FILEMGMT-426 @happy
  Scenario: Tag SetValue sets tag value
    Given a Tag and a new value
    When SetValue is called on the Tag with the value
    Then the tag value is set as specified
    And the behavior matches the Tag SetValue method specification
    And value type is validated
    And GetValue returns the set value after SetValue
