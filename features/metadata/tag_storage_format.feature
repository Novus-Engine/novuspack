@domain:metadata @m2 @REQ-META-015 @spec(metadata.md#11-tag-storage-format)
Feature: Tag storage format defines tag encoding structure

  @REQ-META-015 @happy
  Scenario: Tag storage format defines encoding structure
    Given a per-file tags system with tag storage
    When tags are stored or retrieved
    Then the tag storage format defines encoding structure as specified
    And encoding is consistent across operations
    And the behavior matches the tag storage format specification
    And tag data is persisted correctly
