@domain:metadata @m2 @REQ-META-047 @spec(metadata.md#custom-metadata)
Feature: Custom Metadata Tags

  @REQ-META-047 @happy
  Scenario: Custom metadata provides custom tag support
    Given a NovusPack package
    When custom metadata is used
    Then custom object provides extensible key-value pairs
    And additional metadata can be stored
    And custom tags support any valid tag value types

  @REQ-META-047 @happy
  Scenario: Custom metadata stores extensible key-value pairs
    Given a NovusPack package
    And custom metadata fields
    When custom metadata is set
    Then custom object stores key-value pairs
    And values can be any supported tag value type
    And custom fields extend package metadata

  @REQ-META-047 @happy
  Scenario: Custom metadata examples demonstrate usage
    Given a NovusPack package
    When custom metadata examples are examined
    Then build_number can be stored as integer
    And beta_version can be stored as boolean
    And dlc_ready can be stored as boolean
    And achievements can be stored as integer
    And custom fields provide flexibility

  @REQ-META-047 @error
  Scenario: Custom metadata validates key-value pairs
    Given a NovusPack package
    When invalid custom metadata is provided
    Then key validation detects invalid keys
    And value validation detects invalid values
    And appropriate errors are returned
