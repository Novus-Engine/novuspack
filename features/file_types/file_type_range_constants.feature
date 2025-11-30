@domain:file_types @m2 @REQ-FILETYPES-014 @spec(file_type_system.md#3111-file-type-range-constants)
Feature: File Type Range Constants

  @REQ-FILETYPES-014 @happy
  Scenario: File type range constants define all category ranges
    Given a NovusPack package
    When file type range constants are examined
    Then FileTypeBinaryStart is 0
    And FileTypeBinaryEnd is 999
    And FileTypeTextStart is 1000
    And FileTypeTextEnd is 1999
    And FileTypeScriptStart is 2000
    And FileTypeScriptEnd is 3999
    And FileTypeConfigStart is 4000
    And FileTypeConfigEnd is 4999
    And FileTypeImageStart is 5000
    And FileTypeImageEnd is 6999
    And FileTypeAudioStart is 7000
    And FileTypeAudioEnd is 7999
    And FileTypeVideoStart is 8000
    And FileTypeVideoEnd is 9999
    And FileTypeSystemStart is 10000
    And FileTypeSystemEnd is 10999
    And FileTypeSpecialStart is 65000
    And FileTypeSpecialEnd is 65535

  @REQ-FILETYPES-014 @happy
  Scenario: File type range constants are 2-byte values
    Given a NovusPack package
    When file type range constants are examined
    Then all range constants are uint16 values
    And constants can hold values from 0 to 65535
    And constants represent valid file type ranges

  @REQ-FILETYPES-014 @happy
  Scenario: File type range constants support category checking
    Given a NovusPack package
    When category checking functions are used
    Then range constants enable efficient range checking
    And IsBinaryFile uses FileTypeBinaryStart and FileTypeBinaryEnd
    And IsTextFile uses FileTypeTextStart and FileTypeTextEnd
    And IsScriptFile uses FileTypeScriptStart and FileTypeScriptEnd
    And IsConfigFile uses FileTypeConfigStart and FileTypeConfigEnd
    And IsImageFile uses FileTypeImageStart and FileTypeImageEnd
    And IsAudioFile uses FileTypeAudioStart and FileTypeAudioEnd
    And IsVideoFile uses FileTypeVideoStart and FileTypeVideoEnd
    And IsSystemFile uses FileTypeSystemStart and FileTypeSystemEnd
    And IsSpecialFile uses FileTypeSpecialStart and FileTypeSpecialEnd

  @REQ-FILETYPES-014 @error
  Scenario: File type range constants handle reserved range
    Given a NovusPack package
    When file type value is in reserved range 11000-64999
    Then reserved range is not assigned to any category
    And reserved range is available for future use
