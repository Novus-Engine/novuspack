@domain:file_types @m2 @REQ-FILETYPES-017 @REQ-FILETYPES-029 @spec(file_type_system.md#31122-text-file-types-1000-1999)
Feature: Text File Type Constants

  @REQ-FILETYPES-017 @happy
  Scenario: Text file type constants define text file type range
    Given a NovusPack package
    When text file type constants are examined
    Then FileTypeTextStart is 1000
    And FileTypeTextEnd is 1999
    And FileTypeText is 1000
    And text file types are within range 1000-1999

  @REQ-FILETYPES-017 @happy
  Scenario: Specific text file type constants are defined
    Given a NovusPack package
    When text file type constants are examined
    Then FileTypeText is 1000
    And FileTypeMarkdown is 1001
    And FileTypeHTML is 1002
    And FileTypeCSS is 1003
    And FileTypeCSV is 1004
    And FileTypeSQL is 1005
    And FileTypeDiff is 1006
    And FileTypeTeX is 1007

  @REQ-FILETYPES-017 @happy
  Scenario: Text file types are recognized by IsTextFile
    Given a NovusPack package
    When FileTypeText is checked with IsTextFile
    Then IsTextFile returns true
    When FileTypeMarkdown is checked with IsTextFile
    Then IsTextFile returns true
    When FileTypeHTML is checked with IsTextFile
    Then IsTextFile returns true
    When FileTypeCSS is checked with IsTextFile
    Then IsTextFile returns true

  @REQ-FILETYPES-017 @error
  Scenario: Non-text file types are not recognized by IsTextFile
    Given a NovusPack package
    When FileTypeBinary is checked with IsTextFile
    Then IsTextFile returns false
    When FileTypeImage is checked with IsTextFile
    Then IsTextFile returns false
    When FileTypeScript is checked with IsTextFile
    Then IsTextFile returns false

  @REQ-FILETYPES-029 @happy
  Scenario: Text file analysis provides content-based type detection
    Given a NovusPack package
    And file data with text content
    When text file analysis is performed
    Then data length is checked and must be greater than 0
    And each byte is analyzed to determine if content is text
    And control characters are handled appropriately
    And newlines, carriage returns, and tabs are allowed
    And content must contain printable ASCII characters
    And FileTypeText is returned if content is valid text

  @REQ-FILETYPES-029 @happy
  Scenario: Text file analysis validates text content correctly
    Given a NovusPack package
    And file data containing printable ASCII text
    When text file analysis processes the data
    Then text detection succeeds
    And FileTypeText is returned
    When file data contains non-printable binary characters
    Then text detection fails
    And analysis continues to next detection stage

  @REQ-FILETYPES-029 @error
  Scenario: Text file analysis handles empty or invalid data
    Given a NovusPack package
    And empty file data
    When text file analysis is performed
    Then data validation fails
    And FileTypeText is not returned
    And analysis continues to default classification
