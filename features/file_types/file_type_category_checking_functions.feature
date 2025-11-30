@domain:file_types @m2 @REQ-FILETYPES-009 @spec(file_type_system.md#21-category-checking-functions)
Feature: File Type Category Checking Functions

  @REQ-FILETYPES-009 @happy
  Scenario: Category checking functions provide type category access
    Given a NovusPack package
    And a file type value
    When category checking functions are called
    Then IsBinaryFile checks if file type is in range 0-999
    And IsTextFile checks if file type is in range 1000-1999
    And IsScriptFile checks if file type is in range 2000-3999
    And IsConfigFile checks if file type is in range 4000-4999
    And IsImageFile checks if file type is in range 5000-6999
    And IsAudioFile checks if file type is in range 7000-7999
    And IsVideoFile checks if file type is in range 8000-9999
    And IsSystemFile checks if file type is in range 10000-10999
    And IsSpecialFile checks if file type is in range 65000-65535

  @REQ-FILETYPES-009 @happy
  Scenario: IsBinaryFile correctly identifies binary file types
    Given a NovusPack package
    When IsBinaryFile is called with FileTypeBinary
    Then IsBinaryFile returns true
    When IsBinaryFile is called with FileTypeExecutable
    Then IsBinaryFile returns true
    When IsBinaryFile is called with FileTypeText
    Then IsBinaryFile returns false

  @REQ-FILETYPES-009 @happy
  Scenario: IsTextFile correctly identifies text file types
    Given a NovusPack package
    When IsTextFile is called with FileTypeText
    Then IsTextFile returns true
    When IsTextFile is called with FileTypeMarkdown
    Then IsTextFile returns true
    When IsTextFile is called with FileTypeBinary
    Then IsTextFile returns false

  @REQ-FILETYPES-009 @happy
  Scenario: IsScriptFile correctly identifies script file types
    Given a NovusPack package
    When IsScriptFile is called with FileTypeScript
    Then IsScriptFile returns true
    When IsScriptFile is called with FileTypePython
    Then IsScriptFile returns true
    When IsScriptFile is called with FileTypeJavaScript
    Then IsScriptFile returns true
    When IsScriptFile is called with FileTypeText
    Then IsScriptFile returns false

  @REQ-FILETYPES-009 @happy
  Scenario: Category checking functions handle boundary values
    Given a NovusPack package
    When IsBinaryFile is called with FileTypeBinaryStart
    Then IsBinaryFile returns true
    When IsBinaryFile is called with FileTypeBinaryEnd
    Then IsBinaryFile returns true
    When IsTextFile is called with FileTypeTextStart
    Then IsTextFile returns true
    When IsTextFile is called with FileTypeTextEnd
    Then IsTextFile returns true
    When IsSpecialFile is called with FileTypeSpecialStart
    Then IsSpecialFile returns true
    When IsSpecialFile is called with FileTypeSpecialEnd
    Then IsSpecialFile returns true
