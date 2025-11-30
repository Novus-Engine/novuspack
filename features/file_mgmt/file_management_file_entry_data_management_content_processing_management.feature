@domain:file_mgmt @m2 @REQ-FILEMGMT-045 @spec(api_file_management.md#122-data-management)
Feature: File Management: File Entry Data Management (Content Processing)

  @REQ-FILEMGMT-045 @happy
  Scenario: Data management methods handle file content loading and processing
    Given a file entry
    When data management methods are examined
    Then file content loading is handled
    And file content processing is handled
    And data management supports file operations

  @REQ-FILEMGMT-045 @happy
  Scenario: LoadData loads file content into memory
    Given a file entry
    And file content exists in package
    When LoadData is called
    Then file content is loaded into memory
    And data is prepared for access and processing
    And decompression or decryption may be triggered
    And IsDataLoaded is set to true

  @REQ-FILEMGMT-045 @happy
  Scenario: LoadData prepares data for access and processing
    Given a file entry
    And file content needs to be accessed
    When LoadData is called
    Then data is prepared for access
    And data is prepared for processing
    And file content is available in memory

  @REQ-FILEMGMT-045 @happy
  Scenario: ProcessData applies compression and encryption to file data
    Given a file entry
    And file data needs processing
    When ProcessData is called
    Then compression is applied to file data if configured
    And encryption is applied to file data if configured
    And file entry metadata is updated
    And data is prepared for storage

  @REQ-FILEMGMT-045 @happy
  Scenario: ProcessData updates file entry metadata
    Given a file entry
    And ProcessData is called
    When file entry metadata is examined
    Then file entry metadata is updated
    And compression metadata reflects processing
    And encryption metadata reflects processing
    And metadata matches processed data state

  @REQ-FILEMGMT-045 @happy
  Scenario: Data management supports in-memory and streaming operations
    Given a file entry
    When data management is used
    Then in-memory operations are supported
    And streaming operations are supported
    And data source flexibility is enabled

  @REQ-FILEMGMT-045 @error
  Scenario: LoadData returns error for I/O failures
    Given a file entry
    And I/O error occurs during loading
    When LoadData is called
    Then structured I/O error is returned
    And error indicates data loading failure

  @REQ-FILEMGMT-045 @error
  Scenario: ProcessData returns error for decryption failures
    Given a file entry
    And encrypted file data
    And decryption key is invalid or missing
    When ProcessData is called
    Then structured decryption error is returned
    And error indicates decryption failure

  @REQ-FILEMGMT-045 @error
  Scenario: ProcessData returns error for decompression failures
    Given a file entry
    And compressed file data
    And decompression fails
    When ProcessData is called
    Then structured decompression error is returned
    And error indicates decompression failure
