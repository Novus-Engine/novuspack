@domain:file_types @m2 @REQ-FILETYPES-026 @spec(file_type_system.md#41-detection-process)
Feature: File Type Detection Process

  @REQ-FILETYPES-026 @happy
  Scenario: Detection process uses multi-stage detection algorithm
    Given a NovusPack package
    And a file with name and content
    When DetermineFileType processes the file
    Then extension-based detection is attempted first
    And content-based detection using mimetype library is attempted second
    And MIME type mapping is attempted third
    And extension fallback is attempted fourth
    And text file analysis is attempted fifth
    And default classification is used as final fallback

  @REQ-FILETYPES-026 @happy
  Scenario: Detection process stops at first successful stage
    Given a NovusPack package
    And a file with extension ".ogg"
    When DetermineFileType processes the file
    Then extension-based detection identifies FileTypeOGG
    And no further detection stages are executed
    And file type is returned immediately

  @REQ-FILETYPES-026 @happy
  Scenario: Detection process falls through stages when earlier stages fail
    Given a NovusPack package
    And a file with unknown extension
    And content that can be detected by mimetype library
    When DetermineFileType processes the file
    Then extension-based detection fails
    And content-based detection succeeds
    And MIME type is mapped to file type constant
    And file type is returned

  @REQ-FILETYPES-026 @error
  Scenario: Detection process handles invalid inputs
    Given a NovusPack package
    And empty file name
    And empty file content
    When DetermineFileType processes the file
    Then detection process handles invalid inputs gracefully
    And default classification returns FileTypeBinary
