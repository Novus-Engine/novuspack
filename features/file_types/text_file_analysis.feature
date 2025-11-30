@domain:file_types @m2 @REQ-FILETYPES-029 @spec(file_type_system.md#4113-text-file-analysis)
Feature: Text File Analysis

  @REQ-FILETYPES-029 @happy
  Scenario: Text file analysis performs data validation
    Given a NovusPack package
    And file data
    When text file analysis is performed
    Then data length is checked and must be greater than 0
    And each byte is analyzed to determine if content is text
    And analysis proceeds to text detection if data is valid

  @REQ-FILETYPES-029 @happy
  Scenario: Text file analysis validates text content
    Given a NovusPack package
    And file data with text content
    When text file analysis processes the data
    Then control characters are handled appropriately
    And newlines are allowed
    And carriage returns are allowed
    And tabs are allowed
    And content must contain printable ASCII characters
    And FileTypeText is returned if content is valid text

  @REQ-FILETYPES-029 @happy
  Scenario: Text file analysis detects valid text files
    Given a NovusPack package
    And file data containing printable ASCII text
    When text file analysis processes the data
    Then text detection succeeds
    And FileTypeText is returned
    And analysis completes successfully

  @REQ-FILETYPES-029 @error
  Scenario: Text file analysis handles invalid text content
    Given a NovusPack package
    And file data containing non-printable binary characters
    When text file analysis processes the data
    Then text detection fails
    And FileTypeText is not returned
    And analysis continues to next detection stage
    When empty file data is processed
    Then data validation fails
    And FileTypeText is not returned
