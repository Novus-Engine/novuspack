@domain:file_mgmt @m2 @REQ-FILEMGMT-007 @spec(api_file_management.md#25-addfileoptions-configuration)
Feature: AddFileOptions configuration

  @happy
  Scenario: AddFileOptions contains file processing options
    Given AddFileOptions structure
    When options are examined
    Then CompressionType option exists
    And CompressionLevel option exists
    And FileType option exists
    And PreservePermissions option exists
    And PreserveTimestamps option exists

  @happy
  Scenario: AddFileOptions contains encryption options
    Given AddFileOptions structure
    When encryption options are examined
    Then EncryptionType option exists
    And EncryptionKey option exists
    And KeyID option exists
    And encryption options configure file encryption

  @happy
  Scenario: AddFileOptions contains pattern-specific options
    Given AddFileOptions structure
    When pattern options are examined
    Then Pattern option exists
    And Recursive option exists
    And IncludePatterns option exists
    And ExcludePatterns option exists
    And FollowSymlinks option exists

  @happy
  Scenario: AddFileOptions applies compression settings
    Given AddFileOptions with CompressionType and CompressionLevel
    When file is added with options
    Then compression settings are applied
    And file is compressed according to options

  @happy
  Scenario: AddFileOptions applies encryption settings
    Given AddFileOptions with EncryptionType and EncryptionKey
    When file is added with options
    Then encryption settings are applied
    And file is encrypted according to options

  @happy
  Scenario: AddFileOptions preserves file metadata
    Given AddFileOptions with PreservePermissions and PreserveTimestamps
    When file is added with options
    Then file permissions are preserved
    And file timestamps are preserved
    And metadata matches original file

  @happy
  Scenario: AddFileOptions respects pattern options
    Given AddFileOptions with pattern options
    When files are added with pattern
    Then recursive option is respected
    And include/exclude patterns are applied
    And symlink following is controlled

  @error
  Scenario: Invalid AddFileOptions are rejected
    Given AddFileOptions with invalid configuration
    When file is added
    Then structured validation error is returned
    And error indicates invalid option
