@domain:metadata @m2 @REQ-META-024 @spec(metadata.md#127-localization)
Feature: Localization Metadata Tags

  @REQ-META-024 @happy
  Scenario: Localization provides language code support
    Given a NovusPack package
    When localization tag value type is used
    Then Language type (0x0F) supports language codes
    And language code follows ISO 639-1 format
    And value is stored as UTF-8 string

  @REQ-META-024 @happy
  Scenario: Language type stores ISO 639-1 language codes
    Given a NovusPack package
    And a tag with Language value type
    When Language tag is set
    Then value is stored as UTF-8 string
    And language code follows ISO 639-1 standard
    And language codes are two-letter codes

  @REQ-META-024 @happy
  Scenario: Language tags support multilingual packages
    Given a NovusPack package
    And files with different language tags
    When language tags are used
    Then files can be tagged with language codes
    And language tagging enables localization support

  @REQ-META-024 @error
  Scenario: Language type validates ISO 639-1 format
    Given a NovusPack package
    When invalid language code format is provided
    Then format validation detects invalid codes
    And appropriate error is returned
