@domain:metadata @m2 @REQ-META-022 @spec(metadata.md#125-networkcommunication)
Feature: Network and Communication Metadata Tags

  @REQ-META-022 @happy
  Scenario: Network/communication provides URL and email support
    Given a NovusPack package
    When network/communication tag value types are used
    Then URL type (0x0B) supports URL strings
    And Email type (0x0C) supports email addresses
    And values are stored as UTF-8 strings

  @REQ-META-022 @happy
  Scenario: URL type stores URL strings
    Given a NovusPack package
    And a tag with URL value type
    When URL tag is set
    Then value is stored as UTF-8 string
    And URL can be any valid URL format

  @REQ-META-022 @happy
  Scenario: Email type stores email addresses
    Given a NovusPack package
    And a tag with Email value type
    When Email tag is set
    Then value is stored as UTF-8 string
    And email can be any valid email format

  @REQ-META-022 @error
  Scenario: Network/communication types validate format
    Given a NovusPack package
    When invalid URL or email format is provided
    Then format validation detects invalid values
    And appropriate error is returned
