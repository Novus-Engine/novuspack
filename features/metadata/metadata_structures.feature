@skip @domain:metadata @m2 @spec(metadata.md#122-structured-data)
Feature: Metadata Structures

# This feature captures metadata structure and value-type encoding expectations from the metadata spec.
# Detailed runnable scenarios live in the dedicated metadata feature files.

  @format
  Scenario: Structured metadata values support JSON and YAML encodings
    Given a metadata tag value that represents structured data
    When the value type is JSON
    Then the stored value is a UTF-8 encoded JSON object or array
    And parsers treat invalid JSON as a validation error

  @format
  Scenario: Structured metadata values support string lists
    Given a metadata tag value that represents a list of strings
    When the value type is StringList
    Then the stored value is a comma-separated list of strings
    And consumers split the list without losing ordering information
