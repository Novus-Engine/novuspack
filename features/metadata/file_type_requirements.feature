@domain:metadata @m2 @REQ-META-092 @spec(api_metadata.md#8311-file-type-requirements)
Feature: File Type Requirements

  @REQ-META-092 @happy
  Scenario: Special metadata files use special file types
    Given an open NovusPack package
    And a special metadata file
    When file type is examined
    Then file uses special file types in range 65000-65535
    And file type follows special file type specifications

  @REQ-META-092 @happy
  Scenario: Special metadata files have reserved file names
    Given an open NovusPack package
    And a special metadata file
    When file name is examined
    Then file has reserved file name
    And file name follows naming convention (e.g., "__NPK_DIR_65001__.npkdir")
    And reserved name ensures uniqueness

  @REQ-META-092 @happy
  Scenario: Special metadata files are uncompressed
    Given an open NovusPack package
    And a special metadata file
    When file compression is examined
    Then file is uncompressed
    And CompressionType is set to 0
    And FastWrite compatibility is maintained

  @REQ-META-092 @happy
  Scenario: Special metadata files have proper package header flags
    Given an open NovusPack package
    And a special metadata file
    When package header flags are examined
    Then bit 6 is set to 1 when special files exist
    And bit 5 is set to 1 if directory metadata provides inheritance
    And flags indicate special file presence

  @REQ-META-092 @happy
  Scenario: File type requirements validate special file type range
    Given an open NovusPack package
    And a file entry
    When special file type is validated
    Then file type must be in range 65000-65535
    And file type must match special file specifications
    And validation ensures correct type assignment

  @REQ-META-092 @happy
  Scenario: File type requirements ensure reserved name usage
    Given an open NovusPack package
    And a special metadata file
    When file name is validated
    Then file name must match reserved name pattern
    And file name must be case-sensitive
    And reserved name must follow convention

  @REQ-META-011 @error
  Scenario: File type validation fails with invalid file type
    Given an open NovusPack package
    And a file entry with invalid special file type
    When file type is validated
    Then structured validation error is returned
    And error indicates invalid file type

  @REQ-META-011 @error
  Scenario: File type validation fails with wrong reserved name
    Given an open NovusPack package
    And a special metadata file with incorrect name
    When file name is validated
    Then structured validation error is returned
    And error indicates invalid file name

  @REQ-META-011 @error
  Scenario: File type validation fails if file is compressed
    Given an open NovusPack package
    And a special metadata file with compression
    When file compression is validated
    Then structured validation error is returned
    And error indicates compression not allowed
