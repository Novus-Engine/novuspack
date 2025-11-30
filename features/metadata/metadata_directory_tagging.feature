@domain:metadata @m2 @REQ-META-049 @spec(metadata.md#directory-tagging)
Feature: Metadata Directory Tagging

  @REQ-META-049 @happy
  Scenario: Directory tagging sets tags on directories
    Given a NovusPack package
    And a directory
    When directory tagging is used
    Then tags are set on directories
    And tags are inherited by child files
    And directory tags support inheritance

  @REQ-META-049 @happy
  Scenario: Directory tags are inherited by child files
    Given a NovusPack package
    And a directory with tags
    And child files in the directory
    When directory tagging is used
    Then tags set on directories are inherited by child files
    And inheritance follows tag inheritance rules
    And child files receive directory tags

  @REQ-META-049 @happy
  Scenario: Directory tagging example demonstrates usage
    Given a NovusPack package
    And a textures directory
    When directory is tagged
    Then category tag is set to "texture"
    And compression tag is set to "lossless"
    And mipmaps tag is set to true
    And all files in directory inherit these tags

  @REQ-META-049 @error
  Scenario: Directory tagging validates tag values
    Given a NovusPack package
    When invalid directory tag values are provided
    Then tag validation detects invalid values
    And appropriate errors are returned
