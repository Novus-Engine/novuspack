@domain:file_mgmt @m2 @REQ-FILEMGMT-425 @spec(api_file_mgmt_file_entry.md#193-taggetvalue-method)
Feature: Tag GetValue method returns the tag value

  @REQ-FILEMGMT-425 @happy
  Scenario: Tag GetValue returns tag value
    Given a Tag with a value
    When GetValue is called on the Tag
    Then the tag value is returned as specified
    And the behavior matches the Tag GetValue method specification
    And return type matches Tag value type
    And value is returned without modification
