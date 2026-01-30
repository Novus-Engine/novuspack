@skip @domain:file_types @m2 @REQ-FILETYPES-004 @spec(file_type_system.md#1-filetype-system-specification)
Feature: File Type System Definitions

# This feature captures file type system architecture expectations from the file type system spec.
# More detailed runnable scenarios for type detection and behavior live in dedicated file_types feature files.

  @REQ-FILETYPES-005 @architecture
  Scenario: File type IDs are organized into the defined ranges
    Given a file has an assigned FileType value
    When the FileType is compared to the published range architecture
    Then the FileType falls into exactly one category range
    And the special file range remains reserved for internal package files

  @REQ-FILETYPES-006 @REQ-FILETYPES-007 @constraint
  Scenario: Special metadata files use unique NovusPack extensions
    Given a special metadata file is created for a package
    When the file is named for storage
    Then the file name uses the NovusPack special prefix and a unique extension
    And metadata, manifest, index, and signature files use distinct extensions
