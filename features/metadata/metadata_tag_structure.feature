@domain:metadata @m2 @REQ-META-016 @spec(metadata.md#111-tag-structure)
Feature: Metadata Tag Structure

  @REQ-META-016 @happy
  Scenario: Tag structure defines binary tag storage format
    Given a NovusPack package
    When tag structure is examined
    Then TagCount field is 2 bytes indicating number of tags
    And Tags field is variable-length array of tag entries
    And each tag entry follows defined structure

  @REQ-META-016 @happy
  Scenario: Tag entry structure contains key and value components
    Given a NovusPack package
    And a tag entry
    When tag entry structure is examined
    Then KeyLength is 1 byte indicating key string length
    And Key is UTF-8 encoded string with no null termination
    And ValueType is 1 byte indicating value type
    And ValueLength is 2 bytes indicating value length
    And Value is variable-length UTF-8 encoded based on ValueType

  @REQ-META-016 @happy
  Scenario: Tag entry key is UTF-8 encoded string
    Given a NovusPack package
    And a tag entry
    When tag entry key is examined
    Then key is UTF-8 encoded
    And key has no null termination
    And key length is stored in KeyLength field

  @REQ-META-016 @happy
  Scenario: Tag entry value type specifies value format
    Given a NovusPack package
    And a tag entry
    When ValueType is examined
    Then ValueType indicates value format (string, integer, float, boolean, JSON, YAML, etc.)
    And Value is UTF-8 encoded based on ValueType
    And ValueLength indicates number of bytes in value

  @REQ-META-016 @error
  Scenario: Tag structure validates format constraints
    Given a NovusPack package
    When invalid tag structure is detected
    Then format validation detects violations
    And appropriate errors are returned
