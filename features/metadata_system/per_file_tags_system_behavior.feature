@domain:metadata_system @m2 @REQ-METASYS-001 @spec(metadata.md#1-per-file-tags-system-specification)
Feature: Per-file tags system behavior

  @happy
  Scenario: Tags schema and constraints validated
    Given a file entry
    When I set tags according to schema
    Then tags should be persisted and validated

  @happy
  Scenario: Tags support multiple value types
    Given a file entry
    When tags with different value types are set
    Then string tags are supported
    And integer tags are supported
    And float tags are supported
    And boolean tags are supported
    And structured data tags are supported

  @happy
  Scenario: Tags are stored in OptionalData
    Given a file entry with tags
    When file entry is examined
    Then tags are stored in OptionalData.Tags
    And Tags field is accessible
    And tags are persisted correctly

  @error
  Scenario: Invalid tag schemas are rejected
    Given a file entry
    When invalid tag schema is used
    Then structured validation error is returned
    And error indicates schema violation
