@domain:file_types @m2 @REQ-FILETYPES-023 @spec(file_type_system.md#31128-system-file-types-10000-10999)
Feature: System File Type Constants

  @REQ-FILETYPES-023 @happy
  Scenario: System file type constants define system file type range
    Given a NovusPack package
    When system file type constants are examined
    Then FileTypeSystemStart is 10000
    And FileTypeSystemEnd is 10999
    And system file types are within range 10000-10999

  @REQ-FILETYPES-023 @happy
  Scenario: Specific system file type constants are defined
    Given a NovusPack package
    When system file type constants are examined
    Then FileTypeRegular is 10000
    And FileTypeDirectory is 10001
    And FileTypeSymlink is 10002

  @REQ-FILETYPES-023 @happy
  Scenario: System file types are recognized by IsSystemFile
    Given a NovusPack package
    When FileTypeRegular is checked with IsSystemFile
    Then IsSystemFile returns true
    When FileTypeDirectory is checked with IsSystemFile
    Then IsSystemFile returns true
    When FileTypeSymlink is checked with IsSystemFile
    Then IsSystemFile returns true

  @REQ-FILETYPES-023 @error
  Scenario: Non-system file types are not recognized by IsSystemFile
    Given a NovusPack package
    When FileTypeText is checked with IsSystemFile
    Then IsSystemFile returns false
    When FileTypeImage is checked with IsSystemFile
    Then IsSystemFile returns false
    When FileTypeBinary is checked with IsSystemFile
    Then IsSystemFile returns false

  @REQ-FILETYPES-023 @happy
  Scenario: System file types support system validation
    Given a NovusPack package
    And a system file with type FileTypeDirectory
    When the system file is processed
    Then system validation is performed
    And path handling is appropriate for system files
