@domain:metadata @m2 @REQ-META-100 @spec(api_metadata.md#842-pathmetadataentry-filesystem-properties)
Feature: PathMetadataEntry Filesystem Properties

  @REQ-META-100 @happy
  Scenario: PathMetadataEntry filesystem properties provide path filesystem information
    Given a NovusPack package
    And a PathMetadataEntry
    When filesystem properties are examined
    Then Mode provides Unix/Linux path permissions
    And UID and GID provide user and group IDs
    And ACL provides Access Control List entries
    And WindowsAttrs provides Windows path attributes
    And ExtendedAttrs provides extended attributes map
    And Flags provides filesystem-specific flags

  @REQ-META-100 @happy
  Scenario: Unix/Linux filesystem properties are supported
    Given a NovusPack package
    And a PathMetadataEntry
    When Unix filesystem properties are set
    Then Mode stores path permissions as octal
    And UID stores user ID
    And GID stores group ID
    And ACL stores access control list entries

  @REQ-META-100 @happy
  Scenario: Windows filesystem properties are supported
    Given a NovusPack package
    And a PathMetadataEntry
    When Windows filesystem properties are set
    Then WindowsAttrs stores Windows attributes
    And Windows-specific properties are preserved

  @REQ-META-100 @happy
  Scenario: Extended attributes are supported
    Given a NovusPack package
    And a PathMetadataEntry
    When extended attributes are set
    Then ExtendedAttrs stores extended attributes map
    And extended attributes can be any key-value pairs
    And filesystem-specific flags can be stored

  @REQ-META-100 @error
  Scenario: Filesystem properties validate values
    Given a NovusPack package
    When invalid filesystem properties are provided
    Then property validation detects invalid values
    And appropriate errors are returned
