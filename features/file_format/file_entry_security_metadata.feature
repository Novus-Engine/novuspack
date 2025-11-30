@domain:file_format @m2 @REQ-FILEFMT-052 @spec(package_file_format.md#4125-security-metadata)
Feature: File Entry Security Metadata

  @REQ-FILEFMT-052 @happy
  Scenario: File entry includes encryption metadata
    Given a NovusPack package
    And a file entry
    When security metadata is examined
    Then encryption metadata is present
    And EncryptionType field indicates encryption algorithm
    And per-file encryption settings are available

  @REQ-FILEFMT-052 @happy
  Scenario: File entry includes compression metadata
    Given a NovusPack package
    And a file entry
    When security metadata is examined
    Then compression metadata is present
    And CompressionType field indicates compression algorithm
    And CompressionLevel field indicates compression level
    And per-file compression settings are available

  @REQ-FILEFMT-052 @happy
  Scenario: Security metadata allows per-file security settings
    Given a NovusPack package
    And multiple file entries
    When security metadata is examined for each file
    Then each file can have different encryption settings
    And each file can have different compression settings
    And per-file security configuration is supported

  @REQ-FILEFMT-052 @happy
  Scenario: Security metadata allows per-file optimization settings
    Given a NovusPack package
    And a file entry
    When security metadata is examined
    Then optimization settings are available per file
    And compression can be configured per file
    And security settings enable per-file optimization

  @REQ-FILEFMT-052 @happy
  Scenario: Security metadata is stored in file entry fixed structure
    Given a NovusPack package
    And a file entry
    When file entry structure is examined
    Then EncryptionType is in fixed structure (1 byte)
    And CompressionType is in fixed structure (1 byte)
    And CompressionLevel is in fixed structure (1 byte)
    And security metadata is part of 64-byte fixed structure
