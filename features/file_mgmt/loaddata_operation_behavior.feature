@domain:file_mgmt @m2 @REQ-FILEMGMT-106 @spec(api_file_management.md#1132-loaddata-behavior)
Feature: LoadData Operation Behavior

  @REQ-FILEMGMT-106 @happy
  Scenario: LoadData loads file content from package into memory
    Given an open NovusPack package
    And a valid context
    And a FileEntry in the package
    When LoadData is called on the FileEntry
    Then file content is loaded from package into memory
    And data is prepared for access and processing
    And file content becomes available

  @REQ-FILEMGMT-106 @happy
  Scenario: LoadData may trigger decompression if needed
    Given an open NovusPack package
    And a valid context
    And a compressed FileEntry
    When LoadData is called on compressed FileEntry
    Then decompression may be triggered
    And decompressed content is loaded into memory
    And data is ready for access

  @REQ-FILEMGMT-106 @happy
  Scenario: LoadData may trigger decryption if needed
    Given an open NovusPack package
    And a valid context
    And an encrypted FileEntry
    When LoadData is called on encrypted FileEntry
    Then decryption may be triggered
    And decrypted content is loaded into memory
    And data is ready for access

  @REQ-FILEMGMT-106 @error
  Scenario: LoadData handles I/O errors during data loading
    Given an open NovusPack package
    And a valid context
    And a FileEntry in the package
    And I/O operation failure occurs
    When LoadData is called and I/O error occurs
    Then I/O error is returned
    And error indicates data loading failure
    And error follows structured error format
