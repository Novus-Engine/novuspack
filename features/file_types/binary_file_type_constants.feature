@domain:file_types @m2 @REQ-FILETYPES-016 @spec(file_type_system.md#31121-binary-file-types-0-999)
Feature: Binary File Type Constants

  @REQ-FILETYPES-016 @happy
  Scenario: Binary file type constants define binary file type range
    Given a NovusPack package
    When binary file type constants are examined
    Then FileTypeBinaryStart is 0
    And FileTypeBinaryEnd is 999
    And FileTypeBinary is 0
    And binary file types are within range 0-999

  @REQ-FILETYPES-016 @happy
  Scenario: Specific binary file type constants are defined
    Given a NovusPack package
    When binary file type constants are examined
    Then FileTypeExecutable is 1
    And FileTypeLibrary is 2
    And FileTypeArchive is 3
    And FileTypeBinary is 0

  @REQ-FILETYPES-016 @happy
  Scenario: Binary file types are recognized by IsBinaryFile
    Given a NovusPack package
    When FileTypeBinary is checked with IsBinaryFile
    Then IsBinaryFile returns true
    When FileTypeExecutable is checked with IsBinaryFile
    Then IsBinaryFile returns true
    When FileTypeLibrary is checked with IsBinaryFile
    Then IsBinaryFile returns true
    When FileTypeArchive is checked with IsBinaryFile
    Then IsBinaryFile returns true

  @REQ-FILETYPES-016 @error
  Scenario: Non-binary file types are not recognized by IsBinaryFile
    Given a NovusPack package
    When FileTypeText is checked with IsBinaryFile
    Then IsBinaryFile returns false
    When FileTypeImage is checked with IsBinaryFile
    Then IsBinaryFile returns false
    When FileTypeScript is checked with IsBinaryFile
    Then IsBinaryFile returns false

  @REQ-FILETYPES-016 @happy
  Scenario: Binary file types support security scanning
    Given a NovusPack package
    And a binary file with type FileTypeExecutable
    When the binary file is processed
    Then security scanning is performed
    And execution validation is appropriate for binary executables
