@domain:file_types @m2 @REQ-FILETYPES-013 @spec(file_type_system.md#311-filetype-definition)
Feature: FileType Definition

  @REQ-FILETYPES-013 @happy
  Scenario: FileType definition provides file type type definition
    Given a NovusPack package
    When FileType type is examined
    Then FileType is defined as uint16
    And FileType represents a file type identifier
    And FileType can hold values from 0 to 65535

  @REQ-FILETYPES-013 @happy
  Scenario: FileType is the authoritative definition for file types
    Given a NovusPack package
    When file type system is used
    Then FileType definition is the authoritative source
    And all other references to file types link to FileType definition
    And FileType definition provides consistent file type representation

  @REQ-FILETYPES-013 @happy
  Scenario: FileType values are used throughout the system
    Given a NovusPack package
    When file type operations are performed
    Then FileType values are used for type identification
    And FileType values are used for category checking
    And FileType values are used for compression selection
    And FileType values are used for handler selection

  @REQ-FILETYPES-013 @error
  Scenario: FileType handles boundary values correctly
    Given a NovusPack package
    When FileType is set to 0
    Then FileType represents FileTypeBinary
    When FileType is set to 65535
    Then FileType represents FileTypeSpecialEnd
    And FileType boundary values are valid
