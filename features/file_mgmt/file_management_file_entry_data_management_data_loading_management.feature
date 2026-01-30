@domain:file_mgmt @m2 @REQ-FILEMGMT-186 @spec(api_file_mgmt_file_entry.md#9-fileentry-data-management)
Feature: File Management: File Entry Data Management (Data Loading)

  @REQ-FILEMGMT-186 @happy
  Scenario: File entry data management loads data into memory
    Given an open NovusPack package
    And a valid context
    And a FileEntry in the package
    When LoadData is called on the FileEntry
    Then file content is loaded from package into memory
    And data is prepared for access and processing
    And decompression or decryption may be triggered if needed

  @REQ-FILEMGMT-186 @happy
  Scenario: File entry data management processes file data
    Given an open NovusPack package
    And a valid context
    And a FileEntry with loaded data
    When ProcessData is called on the FileEntry
    Then compression is applied to file data if requested
    And encryption is applied to file data if requested
    And file entry metadata is updated
    And data is prepared for storage

  @REQ-FILEMGMT-186 @happy
  Scenario: File entry data management manages memory efficiently
    Given an open NovusPack package
    And a valid context
    And a FileEntry in the package
    When LoadData and UnloadData are called
    Then data can be loaded into memory when needed
    And data can be unloaded from memory when done
    And memory usage is managed efficiently

  @REQ-FILEMGMT-186 @happy
  Scenario: File entry data management supports temp file operations
    Given an open NovusPack package
    And a valid context
    And a FileEntry in the package
    When temp file operations are performed
    Then temp file can be created
    And data can be streamed to temp file
    And data can be written to temp file
    And data can be read from temp file
    And temp file can be cleaned up

  @REQ-FILEMGMT-186 @error
  Scenario: File entry data management handles I/O errors
    Given an open NovusPack package
    And a valid context
    And a FileEntry in the package
    And I/O operation failure occurs
    When data operations are attempted
    Then I/O error is returned
    And error follows structured error format
    And error indicates I/O failure

  @REQ-FILEMGMT-186 @error
  Scenario: File entry data management handles decryption failures
    Given an open NovusPack package
    And a valid context
    And an encrypted FileEntry
    And decryption failure occurs
    When LoadData is called
    Then decryption error is returned
    And error indicates decryption failure
    And error follows structured error format

  @REQ-FILEMGMT-186 @error
  Scenario: File entry data management handles decompression failures
    Given an open NovusPack package
    And a valid context
    And a compressed FileEntry
    And decompression failure occurs
    When LoadData is called
    Then decompression error is returned
    And error indicates decompression failure
    And error follows structured error format
