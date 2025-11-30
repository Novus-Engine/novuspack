@domain:file_types @m2 @REQ-FILETYPES-027 @spec(file_type_system.md#411-determinefiletype)
Feature: DetermineFileType

  @REQ-FILETYPES-027 @happy
  Scenario: DetermineFileType identifies file type from content and extension
    Given a NovusPack package
    And a file with name and content
    When DetermineFileType is called with name and data
    Then file type is identified using multi-stage detection process
    And file type is returned as FileType

  @REQ-FILETYPES-027 @happy
  Scenario: DetermineFileType uses extension-based detection first
    Given a NovusPack package
    And a file with extension ".ogg"
    When DetermineFileType processes the file
    Then extension-based detection identifies FileTypeOGG
    And file type is returned without further processing

  @REQ-FILETYPES-027 @happy
  Scenario: DetermineFileType uses content-based detection with mimetype library
    Given a NovusPack package
    And a file with content that can be detected by mimetype library
    When DetermineFileType processes the file
    Then mimetype.Detect analyzes file content
    And MIME type is mapped to specific file type constant
    And appropriate file type is returned

  @REQ-FILETYPES-027 @happy
  Scenario: DetermineFileType uses extension fallback when content detection fails
    Given a NovusPack package
    And a file with extension ".txt"
    And content that cannot be detected by mimetype library
    When DetermineFileType processes the file
    Then content-based detection fails
    And extension fallback mapping identifies FileTypeText
    And file type is returned

  @REQ-FILETYPES-027 @happy
  Scenario: DetermineFileType performs text file analysis for unknown files
    Given a NovusPack package
    And a file with text content
    And no recognizable extension or MIME type
    When DetermineFileType processes the file
    Then text file analysis checks if content is text
    And text detection validates printable ASCII characters
    And FileTypeText is returned if content is valid text

  @REQ-FILETYPES-027 @happy
  Scenario: DetermineFileType defaults to binary for unrecognized files
    Given a NovusPack package
    And a file with unrecognized content and extension
    When DetermineFileType processes the file
    Then all detection stages fail
    And default classification returns FileTypeBinary
