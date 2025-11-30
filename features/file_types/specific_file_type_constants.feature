@domain:file_types @m2 @REQ-FILETYPES-015 @spec(file_type_system.md#3112-specific-file-type-constants)
Feature: Specific File Type Constants

  @REQ-FILETYPES-015 @happy
  Scenario: Specific file type constants define type values
    Given a NovusPack package
    When specific file type constants are examined
    Then each file type category has specific type constants defined
    And specific constants have assigned numeric values
    And specific constants are within their category ranges

  @REQ-FILETYPES-015 @happy
  Scenario: Binary file type constants define specific binary types
    Given a NovusPack package
    When binary file type constants are examined
    Then FileTypeBinary is 0
    And FileTypeExecutable is 1
    And FileTypeLibrary is 2
    And FileTypeArchive is 3
    And binary constants are within range 0-999

  @REQ-FILETYPES-015 @happy
  Scenario: Text file type constants define specific text types
    Given a NovusPack package
    When text file type constants are examined
    Then FileTypeText is 1000
    And FileTypeMarkdown is 1001
    And FileTypeHTML is 1002
    And FileTypeCSS is 1003
    And FileTypeCSV is 1004
    And text constants are within range 1000-1999

  @REQ-FILETYPES-015 @happy
  Scenario: Script file type constants define specific script types
    Given a NovusPack package
    When script file type constants are examined
    Then FileTypeScript is 2000
    And FileTypePython is 2001
    And FileTypeJavaScript is 2002
    And FileTypeTypeScript is 2003
    And FileTypeShell is 2004
    And script constants are within range 2000-3999

  @REQ-FILETYPES-015 @happy
  Scenario: Specific file type constants cover all major file categories
    Given a NovusPack package
    When specific file type constants are examined
    Then binary file types have specific constants
    And text file types have specific constants
    And script file types have specific constants
    And config file types have specific constants
    And image file types have specific constants
    And audio file types have specific constants
    And video file types have specific constants
    And system file types have specific constants
    And special file types have specific constants
