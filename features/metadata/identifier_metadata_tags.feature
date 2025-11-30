@domain:metadata @m2 @REQ-META-020 @spec(metadata.md#123-identifiers)
Feature: Identifier Metadata Tags

  @REQ-META-020 @happy
  Scenario: Identifiers provide UUID, hash, and version support
    Given a NovusPack package
    When identifier tag value types are used
    Then UUID type (0x07) supports UUID strings
    And Hash type (0x08) supports hash/checksum strings
    And Version type (0x09) supports semantic version strings

  @REQ-META-020 @happy
  Scenario: UUID type stores UUID strings
    Given a NovusPack package
    And a tag with UUID value type
    When UUID tag is set
    Then value is stored as UTF-8 string
    And UUID can be any valid UUID format

  @REQ-META-020 @happy
  Scenario: Hash type stores hash/checksum strings
    Given a NovusPack package
    And a tag with Hash value type
    When Hash tag is set
    Then value is stored as UTF-8 string
    And hash can be any hash or checksum format

  @REQ-META-020 @happy
  Scenario: Version type stores semantic version strings
    Given a NovusPack package
    And a tag with Version value type
    When Version tag is set
    Then value is stored as UTF-8 string
    And version follows semantic version format

  @REQ-META-020 @error
  Scenario: Identifier types validate format
    Given a NovusPack package
    When invalid identifier format is provided
    Then format validation detects invalid values
    And appropriate error is returned
