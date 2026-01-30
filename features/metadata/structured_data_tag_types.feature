@domain:metadata @m2 @REQ-META-019 @spec(metadata.md#122-structured-data)
Feature: Structured data provides JSON, YAML, and string list support

  @REQ-META-019 @happy
  Scenario: Structured data tag types provide JSON YAML string list support
    Given tag value types for path metadata or file entries
    When structured data is used as a tag value
    Then JSON, YAML, and string list support are provided as specified
    And the behavior matches the structured data specification
    And encoding and decoding are consistent
    And type safety is preserved for structured data
