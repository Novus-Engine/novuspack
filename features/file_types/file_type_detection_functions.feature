@domain:file_types @m2 @REQ-FILETYPES-025 @spec(file_type_system.md#32-file-type-detection-functions)
Feature: File Type Detection Functions

  @REQ-FILETYPES-025 @happy
  Scenario: File type detection functions provide type detection operations
    Given a NovusPack package
    When file type detection functions are used
    Then DetermineFileType detects file type from name and content
    And SelectCompressionType selects compression based on file type
    And detection functions return appropriate FileType values

  @REQ-FILETYPES-025 @happy
  Scenario: DetermineFileType function signature and behavior
    Given a NovusPack package
    And a file with name string and data bytes
    When DetermineFileType is called with name and data
    Then function accepts string name parameter
    And function accepts []byte data parameter
    And function returns FileType value
    And function uses multi-stage detection process

  @REQ-FILETYPES-025 @happy
  Scenario: SelectCompressionType function signature and behavior
    Given a NovusPack package
    And file data
    And a file type
    When SelectCompressionType is called with data and file type
    Then function accepts []byte data parameter
    And function accepts FileType parameter
    And function returns uint8 compression type
    And function selects compression based on file type rules

  @REQ-FILETYPES-025 @error
  Scenario: File type detection functions handle edge cases
    Given a NovusPack package
    And empty file data
    When DetermineFileType is called
    Then function handles empty data gracefully
    And appropriate file type is returned
    When invalid file type is provided to SelectCompressionType
    Then function handles invalid file type gracefully
    And default compression type is returned
