@domain:file_mgmt @m2 @REQ-FILEMGMT-222 @REQ-FILEMGMT-223 @REQ-FILEMGMT-224 @REQ-FILEMGMT-225 @REQ-FILEMGMT-226 @spec(api_file_management.md#925-findentriesbytype)
Feature: File search by file type

  @REQ-FILEMGMT-223 @REQ-FILEMGMT-224 @happy
  Scenario: FindEntriesByType finds files of specified type
    Given an open package
    And file "texture1.png" has type FileTypeImagePNG
    And file "texture2.png" has type FileTypeImagePNG
    And file "config.yaml" has type FileTypeConfigYAML
    When FindEntriesByType is called with type FileTypeImagePNG
    Then files "texture1.png" and "texture2.png" are returned
    And "config.yaml" is not included

  @REQ-FILEMGMT-223 @REQ-FILEMGMT-224 @happy
  Scenario: FindEntriesByType returns empty slice when no matches
    Given an open package without files of type FileTypeAudioMP3
    When FindEntriesByType is called with type FileTypeAudioMP3
    Then an empty slice is returned
    And no error occurs

  @REQ-FILEMGMT-226 @happy
  Scenario: FindEntriesByType supports filtering by file category
    Given an open package with various file types
    And files include images, audio, and config files
    When FindEntriesByType is called for image types
    Then only image files are returned
    When FindEntriesByType is called for audio types
    Then only audio files are returned

  @REQ-FILEMGMT-226 @happy
  Scenario: FindEntriesByType works with type detection
    Given an open package with files
    And file types are determined automatically
    When FindEntriesByType is called with detected types
    Then files with matching types are returned
    And type detection and lookup are consistent

  @REQ-FILEMGMT-224 @error
  Scenario: FindEntriesByType validates file type parameter
    Given an open package
    When FindEntriesByType is called with invalid file type
    Then structured validation error is returned
    And error indicates invalid file type

  @REQ-FILEMGMT-224 @error
  Scenario: FindEntriesByType respects context cancellation
    Given an open package
    And a cancelled context
    When FindEntriesByType is called
    Then a structured context error is returned
    And error type is context cancellation

  @REQ-FILEMGMT-225 @happy
  Scenario: FindEntriesByType returns complete FileEntry objects
    Given an open package with files of various types
    When FindEntriesByType is called
    Then returned FileEntry objects contain file type information
    And all file metadata is accessible
    And file types match the requested type
