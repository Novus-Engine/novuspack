@domain:file_types @m2 @REQ-FILETYPES-005 @spec(file_type_system.md#11-file-type-range-architecture)
Feature: File Type Range Architecture

  @REQ-FILETYPES-005 @happy
  Scenario: File type range architecture defines all file type categories
    Given a NovusPack package
    When file type range architecture is examined
    Then binary files are in range 0-999
    And text files are in range 1000-1999
    And script files are in range 2000-3999
    And config files are in range 4000-4999
    And image files are in range 5000-6999
    And audio files are in range 7000-7999
    And video files are in range 8000-9999
    And system files are in range 10000-10999
    And reserved range is 11000-64999
    And special files are in range 65000-65535

  @REQ-FILETYPES-005 @happy
  Scenario: File type range architecture enables specialized handling
    Given a NovusPack package
    When file type ranges are used
    Then binary files support security scanning and execution validation
    And text files support text processing and encoding validation
    And script files support syntax validation and security analysis
    And config files support schema validation and config parsing
    And image files support format validation and image processing
    And audio files support format validation and audio processing
    And video files support format validation and video processing
    And system files support system validation and path handling
    And special files support special processing and reserved handling

  @REQ-FILETYPES-005 @happy
  Scenario: File type range architecture is authoritative definition
    Given a NovusPack package
    When file type system is used
    Then file type range architecture is the authoritative definition
    And all other references to file types link to this architecture
    And range architecture provides consistent categorization

  @REQ-FILETYPES-005 @error
  Scenario: File type range architecture handles boundary values
    Given a NovusPack package
    When file type value is at range boundary
    Then range boundaries are inclusive
    And FileTypeBinaryStart is 0
    And FileTypeBinaryEnd is 999
    And FileTypeSpecialEnd is 65535
