@domain:file_mgmt @m2 @REQ-FILEMGMT-138 @spec(api_file_management.md#51-extract-file)
Feature: ExtractFile extracts file content

  @REQ-FILEMGMT-138 @happy
  Scenario: ExtractFile returns file content as byte slice
    Given an open NovusPack package
    And a file exists in the package at a specific path
    When ExtractFile is called with the file path
    Then file content is returned as byte slice
    And content matches original file data
    And content is complete and unmodified

  @REQ-FILEMGMT-138 @happy
  Scenario: ExtractFile handles compressed files
    Given an open NovusPack package
    And a compressed file exists in the package
    When ExtractFile is called with the file path
    Then file content is automatically decompressed
    And decompressed content is returned as byte slice
    And content matches original uncompressed data

  @REQ-FILEMGMT-138 @happy
  Scenario: ExtractFile handles encrypted files
    Given an open NovusPack package
    And an encrypted file exists in the package
    And encryption key is available
    When ExtractFile is called with the file path
    Then file content is automatically decrypted
    And decrypted content is returned as byte slice
    And content matches original unencrypted data

  @REQ-FILEMGMT-138 @error
  Scenario: ExtractFile returns error when package is not open
    Given a NovusPack package that is not open
    When ExtractFile is called with a file path
    Then package not open error is returned
    And error indicates package must be open

  @REQ-FILEMGMT-138 @error
  Scenario: ExtractFile returns error when file does not exist
    Given an open NovusPack package
    And a file does not exist at the specified path
    When ExtractFile is called with the non-existent path
    Then file not found error is returned
    And error indicates file does not exist

  @REQ-FILEMGMT-138 @error
  Scenario: ExtractFile returns error for invalid path
    Given an open NovusPack package
    When ExtractFile is called with invalid or malformed path
    Then invalid path error is returned
    And error indicates path format issue

  @REQ-FILEMGMT-138 @error
  Scenario: ExtractFile handles decompression failures
    Given an open NovusPack package
    And a file with corrupted compression data
    When ExtractFile is called
    Then decompression failed error is returned
    And error indicates compression issue

  @REQ-FILEMGMT-138 @error
  Scenario: ExtractFile handles decryption failures
    Given an open NovusPack package
    And an encrypted file with incorrect key
    When ExtractFile is called
    Then decryption failed error is returned
    And error indicates decryption issue

  @REQ-FILEMGMT-138 @error
  Scenario: ExtractFile respects context cancellation
    Given an open NovusPack package
    And a cancelled context
    When ExtractFile is called with cancelled context
    Then context cancelled error is returned
    And error type is context cancellation

  @REQ-FILEMGMT-138 @error
  Scenario: ExtractFile handles context timeout
    Given an open NovusPack package
    And a context with timeout
    And a large file that exceeds timeout
    When ExtractFile is called
    Then context timeout error is returned
    And error type is context timeout
