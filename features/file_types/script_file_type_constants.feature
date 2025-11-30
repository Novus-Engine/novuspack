@domain:file_types @m2 @REQ-FILETYPES-018 @spec(file_type_system.md#31123-script-file-types-2000-3999)
Feature: Script File Type Constants

  @REQ-FILETYPES-018 @happy
  Scenario: Script file type constants define script file type range
    Given a NovusPack package
    When script file type constants are examined
    Then FileTypeScriptStart is 2000
    And FileTypeScriptEnd is 3999
    And FileTypeScript is 2000
    And script file types are within range 2000-3999

  @REQ-FILETYPES-018 @happy
  Scenario: Specific script file type constants are defined
    Given a NovusPack package
    When script file type constants are examined
    Then FileTypePython is 2001
    And FileTypeJavaScript is 2002
    And FileTypeTypeScript is 2003
    And FileTypeShell is 2004
    And FileTypeLua is 2005
    And FileTypePerl is 2006
    And FileTypeRuby is 2007
    And FileTypePHP is 2008
    And FileTypeGo is 2011
    And FileTypeRust is 2012

  @REQ-FILETYPES-018 @happy
  Scenario: Script file types are recognized by IsScriptFile
    Given a NovusPack package
    When FileTypeScript is checked with IsScriptFile
    Then IsScriptFile returns true
    When FileTypePython is checked with IsScriptFile
    Then IsScriptFile returns true
    When FileTypeJavaScript is checked with IsScriptFile
    Then IsScriptFile returns true
    When FileTypeShell is checked with IsScriptFile
    Then IsScriptFile returns true

  @REQ-FILETYPES-018 @error
  Scenario: Non-script file types are not recognized by IsScriptFile
    Given a NovusPack package
    When FileTypeText is checked with IsScriptFile
    Then IsScriptFile returns false
    When FileTypeImage is checked with IsScriptFile
    Then IsScriptFile returns false
    When FileTypeBinary is checked with IsScriptFile
    Then IsScriptFile returns false

  @REQ-FILETYPES-018 @happy
  Scenario: Script file types support syntax validation
    Given a NovusPack package
    And a script file with type FileTypePython
    When the script file is processed
    Then syntax validation is performed
    And security analysis is appropriate for script files
