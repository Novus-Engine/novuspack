@domain:metadata @m2 @REQ-META-021 @spec(metadata.md#124-time)
Feature: Time Metadata Tags

  @REQ-META-021 @happy
  Scenario: Time provides timestamp support
    Given a NovusPack package
    When time tag value type is used
    Then Timestamp type (0x0A) supports ISO8601 timestamps
    And value is stored as UTF-8 string
    And timestamp follows ISO8601 format

  @REQ-META-021 @happy
  Scenario: Timestamp type stores ISO8601 timestamps
    Given a NovusPack package
    And a tag with Timestamp value type
    When Timestamp tag is set
    Then value is stored as UTF-8 string
    And timestamp follows ISO8601 standard
    And timestamp includes date and time information

  @REQ-META-021 @happy
  Scenario: Timestamps support package metadata timestamps
    Given a NovusPack package
    And package metadata with timestamps
    When timestamps are used
    Then created timestamp can be stored
    And modified timestamp can be stored
    And timestamps enable temporal tracking

  @REQ-META-021 @error
  Scenario: Timestamp type validates ISO8601 format
    Given a NovusPack package
    When invalid timestamp format is provided
    Then format validation detects invalid timestamps
    And appropriate error is returned
