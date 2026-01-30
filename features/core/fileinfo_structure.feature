@domain:core @m1 @REQ-CORE-064 @spec(api_core.md#1147-fileinfo-structure)
Feature: FileInfo structure provides lightweight file information

  Background:
    Given a NovusPack package with multiple files

  @REQ-CORE-065 @happy
  Scenario: FileInfo includes path identification fields
    Given a file with primary path "textures/diffuse.dds"
    When ListFiles is called
    Then FileInfo for the file includes PrimaryPath "textures/diffuse.dds"
    And FileInfo includes Paths array with at least one entry
    And FileInfo includes unique FileID

  @REQ-CORE-065 @happy
  Scenario: FileInfo exposes all aliased paths for multi-path files
    Given a file with primary path "models/player.obj"
    And the file has alias path "assets/player.obj"
    When ListFiles is called
    Then FileInfo for the file includes PrimaryPath "models/player.obj"
    And FileInfo Paths array contains "models/player.obj"
    And FileInfo Paths array contains "assets/player.obj"
    And FileInfo PathCount is 2

  @REQ-CORE-066 @happy
  Scenario: FileInfo includes file type identification
    Given a file with FileType 1000
    And the file type name is "Texture"
    When ListFiles is called
    Then FileInfo for the file includes FileType 1000
    And FileInfo includes FileTypeName "Texture"

  @REQ-CORE-066 @happy
  Scenario: FileInfo includes human-readable file type names
    Given multiple files with different FileTypes
    When ListFiles is called
    Then each FileInfo includes FileType numeric value
    And each FileInfo includes FileTypeName string
    And FileTypeName is derived from FileType via type system lookup

  @REQ-CORE-067 @happy
  Scenario: FileInfo includes size information for uncompressed files
    Given an uncompressed file with original size 4096 bytes
    When ListFiles is called
    Then FileInfo Size is 4096
    And FileInfo StoredSize is 4096

  @REQ-CORE-067 @happy
  Scenario: FileInfo distinguishes original and stored sizes for compressed files
    Given a compressed file with original size 4096 bytes
    And stored size 1024 bytes after compression
    When ListFiles is called
    Then FileInfo Size is 4096
    And FileInfo StoredSize is 1024

  @REQ-CORE-068 @happy
  Scenario: FileInfo indicates compression status
    Given a compressed file using Zstd compression
    When ListFiles is called
    Then FileInfo IsCompressed is true
    And FileInfo CompressionType is 1

  @REQ-CORE-068 @happy
  Scenario: FileInfo indicates encryption status
    Given an encrypted file
    When ListFiles is called
    Then FileInfo IsEncrypted is true

  @REQ-CORE-068 @happy
  Scenario: FileInfo indicates unprocessed file status
    Given an uncompressed and unencrypted file
    When ListFiles is called
    Then FileInfo IsCompressed is false
    And FileInfo IsEncrypted is false
    And FileInfo CompressionType is 0

  @REQ-CORE-069 @happy
  Scenario: FileInfo includes content verification checksums
    Given a file with RawChecksum 0x12345678
    And StoredChecksum 0x12345678
    When ListFiles is called
    Then FileInfo RawChecksum is 0x12345678
    And FileInfo StoredChecksum is 0x12345678

  @REQ-CORE-069 @happy
  Scenario: FileInfo checksums differ for compressed files
    Given a compressed file with RawChecksum 0x12345678
    And StoredChecksum 0x87654321 after compression
    When ListFiles is called
    Then FileInfo RawChecksum is 0x12345678
    And FileInfo StoredChecksum is 0x87654321

  @REQ-CORE-070 @happy
  Scenario: FileInfo indicates single-path files
    Given a file with one path
    When ListFiles is called
    Then FileInfo PathCount is 1
    And FileInfo Paths array has length 1

  @REQ-CORE-070 @happy
  Scenario: FileInfo indicates multi-path files
    Given a file with three aliased paths
    When ListFiles is called
    Then FileInfo PathCount is 3
    And FileInfo Paths array has length 3

  @REQ-CORE-071 @happy
  Scenario: FileInfo includes version tracking
    Given a file with FileVersion 5
    And MetadataVersion 3
    When ListFiles is called
    Then FileInfo FileVersion is 5
    And FileInfo MetadataVersion is 3

  @REQ-CORE-072 @happy
  Scenario: FileInfo indicates presence of tags
    Given a file with custom tags
    When ListFiles is called
    Then FileInfo HasTags is true

  @REQ-CORE-072 @happy
  Scenario: FileInfo indicates absence of tags
    Given a file without tags
    When ListFiles is called
    Then FileInfo HasTags is false

  @REQ-CORE-073 @happy
  Scenario: FileInfo provides lightweight information without full FileEntry
    Given a package with 1000 files
    When ListFiles is called
    Then FileInfo structures are returned quickly
    And no variable-length FileEntry data is loaded
    And only static FileEntry fields are included
    And ListFiles remains a pure in-memory operation

  @REQ-CORE-073 @happy
  Scenario: FileInfo enables filtering without loading full FileEntry
    Given a package with mixed file types
    When ListFiles is called
    And files are filtered by FileType
    Then filtering uses only FileInfo data
    And no full FileEntry objects are loaded

  @happy
  Scenario: FileInfo supports compression ratio calculation
    Given a compressed file with Size 4096
    And StoredSize 1024
    When compression ratio is calculated
    Then ratio is 25.0 percent

  @happy
  Scenario: FileInfo supports deduplication by RawChecksum
    Given multiple files with same RawChecksum
    When files are deduplicated by RawChecksum
    Then duplicate files are identified
    And unique files are preserved

  @happy
  Scenario: FileInfo PrimaryPath is first path lexicographically
    Given a file with paths "models/player.obj", "assets/player.obj", "data/player.obj"
    When ListFiles is called
    Then FileInfo PrimaryPath is "assets/player.obj"
    And Paths array is sorted lexicographically
