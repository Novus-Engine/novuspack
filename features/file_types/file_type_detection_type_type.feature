@domain:file_types @m2 @REQ-FILETYPES-012 @REQ-FILETYPES-014 @REQ-FILETYPES-025 @spec(file_type_system.md#3-file-type-api)
Feature: File Type Detection

  @REQ-FILETYPES-012 @happy
  Scenario: File type API provides file type operations
    Given a NovusPack package
    When file type API is used
    Then file type management operations are available
    And file type detection functions are available
    And file type range constants are available
    And category checking functions are available

  @REQ-FILETYPES-012 @happy
  Scenario: File type API includes FileType definition
    Given a NovusPack package
    When file type API is used
    Then FileType type definition is available
    And FileType represents file type identifier
    And FileType is uint16 type

  @REQ-FILETYPES-014 @happy
  Scenario: File type range constants define type range values
    Given a NovusPack package
    When file type range constants are used
    Then FileTypeBinaryStart and FileTypeBinaryEnd define binary range
    And FileTypeTextStart and FileTypeTextEnd define text range
    And FileTypeScriptStart and FileTypeScriptEnd define script range
    And FileTypeConfigStart and FileTypeConfigEnd define config range
    And FileTypeImageStart and FileTypeImageEnd define image range
    And FileTypeAudioStart and FileTypeAudioEnd define audio range
    And FileTypeVideoStart and FileTypeVideoEnd define video range
    And FileTypeSystemStart and FileTypeSystemEnd define system range
    And FileTypeSpecialStart and FileTypeSpecialEnd define special range

  @REQ-FILETYPES-025 @happy
  Scenario: File type detection functions provide type detection operations
    Given a NovusPack package
    When file type detection functions are used
    Then DetermineFileType function is available
    And SelectCompressionType function is available
    And detection functions accept appropriate parameters
    And detection functions return appropriate values

  @REQ-FILETYPES-012 @error
  Scenario: File type API handles invalid operations
    Given a NovusPack package
    When invalid file type operation is performed
    Then appropriate error is returned
    And error follows structured error format
