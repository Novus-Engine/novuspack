@skip @domain:metadata @m2 @spec(metadata.md#11-tag-storage-format)
Feature: Metadata Definitions

# This feature captures metadata encoding expectations from the metadata spec.
# Detailed runnable scenarios live in the dedicated metadata feature files.

  @format
  Scenario: Tag storage encodes key, value type, length, and value bytes
    Given a file entry with one or more metadata tags
    When the tags are serialized
    Then the serialized format includes a tag count
    And each tag includes a key length, key bytes, value type, value length, and value bytes

  @format
  Scenario: Tag keys and values are stored as UTF-8 strings
    Given a tag key and value to be stored
    When the tag is serialized
    Then the key is stored as a UTF-8 string without null termination
    And the value is stored as a UTF-8 string representation according to its value type
