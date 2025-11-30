@domain:file_mgmt @m2 @REQ-FILEMGMT-135 @spec(api_file_management.md#5-extract-file-operations)
Feature: Extract File Operations

  @REQ-FILEMGMT-135 @happy
  Scenario: Extract file operations support file extraction from packages
    Given an open NovusPack package
    And file "documents/data.txt" exists in the package
    When ExtractFile is called with the file path
    Then file content is extracted and returned
    And content is returned as byte slice
    And extraction operation completes successfully

  @REQ-FILEMGMT-135 @happy
  Scenario: ExtractFile automatically handles decompression
    Given an open package
    And compressed file "archive.zip" exists
    When ExtractFile is called
    Then file content is decompressed automatically
    And decompressed content is returned
    And decompression is transparent to caller

  @REQ-FILEMGMT-135 @happy
  Scenario: ExtractFile automatically handles decryption
    Given an open package
    And encrypted file "secret.txt" exists
    And decryption keys are available
    When ExtractFile is called
    Then file content is decrypted automatically
    And decrypted content is returned
    And decryption is transparent to caller

  @REQ-FILEMGMT-135 @happy
  Scenario: ExtractFile handles combined compression and encryption
    Given an open package
    And file exists with both compression and encryption
    When ExtractFile is called
    Then file is decrypted first
    Then file is decompressed
    And final content is returned
    And operations occur in correct order

  @REQ-FILEMGMT-135 @happy
  Scenario: ExtractFile locates file entry in package index
    Given an open package
    And file entry exists in package index
    When ExtractFile is called with valid path
    Then file entry is located in package index
    And file content is read from data section
    And extraction completes successfully

  @REQ-FILEMGMT-135 @error
  Scenario: ExtractFile returns error when package not open
    Given a package that is not open
    When ExtractFile is called
    Then ErrPackageNotOpen error is returned
    And error follows structured error format

  @REQ-FILEMGMT-135 @error
  Scenario: ExtractFile returns error when file not found
    Given an open package
    And file does not exist at specified path
    When ExtractFile is called with non-existent path
    Then ErrFileNotFound error is returned
    And error indicates file not found
    And error follows structured error format

  @REQ-FILEMGMT-135 @error
  Scenario: ExtractFile returns error when decryption fails
    Given an open package
    And encrypted file exists
    And decryption keys are invalid or missing
    When ExtractFile is called
    Then ErrDecryptionFailed error is returned
    And error indicates decryption failure
    And error follows structured error format

  @REQ-FILEMGMT-135 @error
  Scenario: ExtractFile respects context cancellation
    Given an open package
    And a cancelled context
    When ExtractFile is called
    Then ErrContextCancelled error is returned
    And error follows structured error format
