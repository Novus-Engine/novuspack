@domain:metadata @m2 @REQ-META-048 @spec(metadata.md#custom-metadata-example)
Feature: Custom Metadata Example

  @REQ-META-048 @happy
  Scenario: Custom metadata example stores build number
    Given an open NovusPack package
    And custom metadata
    When custom metadata example is examined
    Then build_number field can be stored as integer
    And build number value is accessible
    And example demonstrates integer custom field

  @REQ-META-048 @happy
  Scenario: Custom metadata example stores beta version flag
    Given an open NovusPack package
    And custom metadata
    When custom metadata example is examined
    Then beta_version field can be stored as boolean
    And beta version flag is accessible
    And example demonstrates boolean custom field

  @REQ-META-048 @happy
  Scenario: Custom metadata example stores DLC ready flag
    Given an open NovusPack package
    And custom metadata
    When custom metadata example is examined
    Then dlc_ready field can be stored as boolean
    And DLC ready flag is accessible
    And example demonstrates boolean custom field

  @REQ-META-048 @happy
  Scenario: Custom metadata example stores achievements count
    Given an open NovusPack package
    And custom metadata
    When custom metadata example is examined
    Then achievements field can be stored as integer
    And achievements count is accessible
    And example demonstrates integer custom field

  @REQ-META-048 @happy
  Scenario: Custom metadata example demonstrates extensible key-value pairs
    Given an open NovusPack package
    And custom metadata with example fields
    When custom metadata is retrieved
    Then custom object provides extensible key-value pairs
    And additional metadata can be stored
    And custom fields provide flexibility

  @REQ-META-048 @happy
  Scenario: Custom metadata example supports various value types
    Given an open NovusPack package
    And custom metadata with various types
    When custom metadata is examined
    Then string values are supported
    And integer values are supported
    And boolean values are supported
    And all valid tag value types are supported

  @REQ-META-011 @error
  Scenario: Custom metadata example validation fails with invalid keys
    Given an open NovusPack package
    And custom metadata with invalid keys
    When custom metadata is validated
    Then structured validation error is returned
    And error indicates invalid custom metadata keys

  @REQ-META-011 @error
  Scenario: Custom metadata example validation fails with invalid values
    Given an open NovusPack package
    And custom metadata with invalid values
    When custom metadata is validated
    Then structured validation error is returned
    And error indicates invalid custom metadata values
