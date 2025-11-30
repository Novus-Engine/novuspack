@domain:metadata @m2 @REQ-META-018 @spec(metadata.md#121-basic-types)
Feature: Basic Metadata Types

  @REQ-META-018 @happy
  Scenario: Basic types provide string, integer, float, and boolean support
    Given a NovusPack package
    When basic tag value types are used
    Then string type (0x00) supports UTF-8 encoded strings
    And integer type (0x01) supports 64-bit signed integers
    And float type (0x02) supports 64-bit floating point numbers
    And boolean type (0x03) supports true or false values

  @REQ-META-018 @happy
  Scenario: String type stores UTF-8 encoded values
    Given a NovusPack package
    And a tag with string value type
    When string tag is set
    Then value is stored as UTF-8 encoded string
    And string can contain any valid UTF-8 characters

  @REQ-META-018 @happy
  Scenario: Integer type stores 64-bit signed integers
    Given a NovusPack package
    And a tag with integer value type
    When integer tag is set
    Then value is stored as UTF-8 string representation
    And integer can be any 64-bit signed integer value

  @REQ-META-018 @happy
  Scenario: Float type stores 64-bit floating point numbers
    Given a NovusPack package
    And a tag with float value type
    When float tag is set
    Then value is stored as UTF-8 string representation
    And float can be any 64-bit floating point number

  @REQ-META-018 @happy
  Scenario: Boolean type stores true or false values
    Given a NovusPack package
    And a tag with boolean value type
    When boolean tag is set
    Then value is stored as UTF-8 string "true" or "false"
    And boolean represents logical true or false state

  @REQ-META-018 @error
  Scenario: Basic types validate value format
    Given a NovusPack package
    And a tag with specific value type
    When invalid value format is provided
    Then value validation detects format mismatch
    And appropriate error is returned
