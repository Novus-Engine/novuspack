@domain:file_mgmt @m2 @REQ-FILEMGMT-424 @spec(api_file_mgmt_file_entry.md#192-newtag-function)
Feature: NewTag function creates a new type-safe tag

  @REQ-FILEMGMT-424 @happy
  Scenario: NewTag creates type-safe tag
    Given tag key and value for a file entry
    When NewTag is called with key and value
    Then a new type-safe tag is created as specified
    And the behavior matches the NewTag function specification
    And returned Tag has correct key and value
    And type safety is preserved
