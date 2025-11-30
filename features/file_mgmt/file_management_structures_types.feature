@domain:file_mgmt @m2 @REQ-FILEMGMT-042 @REQ-FILEMGMT-005 @spec(api_file_management.md#1-core-data-structures)
Feature: File Management Structures

  @REQ-FILEMGMT-042 @REQ-FILEMGMT-005 @happy
  Scenario: Core data structures define fileentry and related types
    Given core data structures
    When structures are examined
    Then FileEntry structure is defined
    And FileType type is defined
    And related types are defined
    And structures provide file metadata representation

  @REQ-FILEMGMT-042 @REQ-FILEMGMT-005 @happy
  Scenario: FileEntry structure provides file metadata representation
    Given FileEntry structure
    When structure is examined
    Then structure contains static fields (64 bytes total)
    And structure contains FileID, size fields, checksums
    And structure contains version fields, counts, types
    And structure contains compression and encryption fields
    And structure contains offset and length fields

  @REQ-FILEMGMT-042 @REQ-FILEMGMT-005 @happy
  Scenario: FileEntry structure supports variable-length data
    Given FileEntry structure
    When structure is examined
    Then structure supports path entries
    And structure supports hash data
    And structure supports optional data
    And variable-length sections are supported

  @REQ-FILEMGMT-042 @REQ-FILEMGMT-005 @happy
  Scenario: FileEntry structure provides comprehensive file information
    Given FileEntry structure
    When structure fields are examined
    Then file identification is provided via FileID
    And file size information is provided (OriginalSize, StoredSize)
    And file checksums are provided (RawChecksum, StoredChecksum)
    And version tracking is provided (FileVersion, MetadataVersion)
    And processing information is provided (compression, encryption)

  @REQ-FILEMGMT-042 @REQ-FILEMGMT-005 @happy
  Scenario: FileType represents file type identifier
    Given FileType type
    When type is examined
    Then FileType is uint16
    And type represents file type identifier from file type system
    And type is used in FileEntry structure

  @REQ-FILEMGMT-042 @happy
  Scenario: Core data structures support file management operations
    Given core data structures
    When structures are used in file operations
    Then AddFile operations use FileEntry
    And UpdateFile operations use FileEntry
    And RemoveFile operations use FileEntry
    And structures enable comprehensive file management
