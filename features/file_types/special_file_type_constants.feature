@domain:file_types @m2 @REQ-FILETYPES-024 @spec(file_type_system.md#31129-special-file-types-65000-65535)
Feature: Special File Type Constants

  @REQ-FILETYPES-024 @happy
  Scenario: Special file type constants define special file type range
    Given a NovusPack package
    When special file type constants are examined
    Then FileTypeSpecialStart is 65000
    And FileTypeSpecialEnd is 65535
    And special file types are within range 65000-65535

  @REQ-FILETYPES-024 @happy
  Scenario: Specific special file type constants are defined
    Given a NovusPack package
    When special file type constants are examined
    Then FileTypeMetadata is 65000
    And FileTypeManifest is 65001
    And FileTypeIndex is 65002
    And FileTypeSignature is 65003

  @REQ-FILETYPES-024 @happy
  Scenario: Special file types are recognized by IsSpecialFile
    Given a NovusPack package
    When FileTypeMetadata is checked with IsSpecialFile
    Then IsSpecialFile returns true
    When FileTypeManifest is checked with IsSpecialFile
    Then IsSpecialFile returns true
    When FileTypeIndex is checked with IsSpecialFile
    Then IsSpecialFile returns true
    When FileTypeSignature is checked with IsSpecialFile
    Then IsSpecialFile returns true

  @REQ-FILETYPES-024 @error
  Scenario: Non-special file types are not recognized by IsSpecialFile
    Given a NovusPack package
    When FileTypeText is checked with IsSpecialFile
    Then IsSpecialFile returns false
    When FileTypeImage is checked with IsSpecialFile
    Then IsSpecialFile returns false
    When FileTypeBinary is checked with IsSpecialFile
    Then IsSpecialFile returns false

  @REQ-FILETYPES-024 @happy
  Scenario: Special file types support special processing
    Given a NovusPack package
    And a special file with type FileTypeMetadata
    When the special file is processed
    Then special processing is performed
    And reserved handling is appropriate for special files
