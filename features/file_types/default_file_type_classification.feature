@domain:file_types @m2 @REQ-FILETYPES-030 @spec(file_type_system.md#4114-default-classification)
Feature: Default File Type Classification

  @REQ-FILETYPES-030 @happy
  Scenario: Default classification provides fallback type assignment
    Given a NovusPack package
    And file data with unknown type
    When DetermineFileType cannot identify file type
    Then default classification assigns FileTypeBinary
    And fallback type assignment completes successfully

  @REQ-FILETYPES-030 @happy
  Scenario: Default classification handles unrecognized files
    Given a NovusPack package
    And file data with no recognizable signature
    And file extension is unknown or missing
    When DetermineFileType processes the file
    Then content-based detection fails
    And extension-based detection fails
    And text file analysis fails
    And default classification returns FileTypeBinary

  @REQ-FILETYPES-030 @happy
  Scenario: Default classification is last stage in detection process
    Given a NovusPack package
    And file data
    When DetermineFileType processes the file
    Then extension-based detection is attempted first
    And content-based detection is attempted second
    And MIME type mapping is attempted third
    And extension fallback is attempted fourth
    And text file analysis is attempted fifth
    And default classification is used as final fallback

  @REQ-FILETYPES-030 @error
  Scenario: Default classification handles empty or invalid data
    Given a NovusPack package
    And empty file data
    When DetermineFileType processes the file
    Then default classification handles empty data
    And appropriate file type is assigned
