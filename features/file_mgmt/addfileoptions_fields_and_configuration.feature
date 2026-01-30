@domain:file_mgmt @m2 @REQ-FILEMGMT-067 @REQ-FILEMGMT-068 @REQ-FILEMGMT-007 @spec(api_file_mgmt_addition.md#252-fields)
Feature: AddFileOptions Fields and Configuration

  @REQ-FILEMGMT-067 @REQ-FILEMGMT-068 @REQ-FILEMGMT-007 @happy
  Scenario: Addfileoptions fields define processing, encryption, and pattern options
    Given AddFileOptions structure
    When structure is examined
    Then file processing options fields exist
    And encryption options fields exist
    And pattern-specific options fields exist
    And unified configuration is provided

  @REQ-FILEMGMT-067 @REQ-FILEMGMT-068 @REQ-FILEMGMT-007 @happy
  Scenario: Addfileoptions file processing options control compression and deduplication
    Given AddFileOptions structure
    When file processing options are examined
    Then Compress field controls compression
    And CompressionType field specifies algorithm (0=none, 1=Zstd, 2=LZ4, 3=LZMA)
    And CompressionLevel field specifies level (1-9, 0=default)
    And FileType field specifies file type identifier
    And Tags field specifies per-file tags (key-value pairs)

  @REQ-FILEMGMT-067 @REQ-FILEMGMT-068 @REQ-FILEMGMT-007 @happy
  Scenario: AddFileOptions encryption options control file encryption
    Given AddFileOptions structure
    When encryption options are examined
    Then Encrypt field controls encryption
    And EncryptionType field specifies algorithm type
    And EncryptionKey field specifies specific key (overrides EncryptionType)
    And encryption options are configurable

  @REQ-FILEMGMT-067 @REQ-FILEMGMT-068 @REQ-FILEMGMT-007 @happy
  Scenario: AddFileOptions pattern-specific options control pattern behavior
    Given AddFileOptions structure
    When pattern-specific options are examined
    Then ExcludePatterns field specifies patterns to exclude
    And MaxFileSize field specifies maximum file size (0=no limit)
    And PreservePaths field controls directory structure preservation
    And pattern options work for bulk operations

  @REQ-FILEMGMT-067 @REQ-FILEMGMT-068 @REQ-FILEMGMT-007 @happy
  Scenario: AddFileOptions provides unified configuration for file addition
    Given AddFileOptions structure
    When options are used for file addition
    Then options work for individual files via AddFile
    And options work for pattern operations via AddFilePattern
    And unified configuration is consistent

  @REQ-FILEMGMT-067 @REQ-FILEMGMT-068 @REQ-FILEMGMT-007 @happy
  Scenario: AddFileOptions supports Option type for optional fields
    Given AddFileOptions structure
    When optional fields are examined
    Then fields use Option type for optional values
    And unset options use default values
    And option handling works correctly

  @REQ-FILEMGMT-067 @REQ-FILEMGMT-068 @error
  Scenario: AddFileOptions validates compression type values
    Given AddFileOptions with invalid CompressionType
    When options are used
    Then structured validation error is returned
    And error indicates invalid compression type
    And error follows structured error format
