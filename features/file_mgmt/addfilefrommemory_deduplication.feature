@domain:file_mgmt @m2 @REQ-FILEMGMT-205 @spec(api_file_mgmt_addition.md#2273-addfilefrommemory-deduplication)
Feature: AddFileFromMemory Deduplication

  @REQ-FILEMGMT-205 @happy
  Scenario: Deduplication same as AddFile
    Given in-memory byte data
    When AddFileFromMemory is called
    Then deduplication is performed
    And deduplication behavior matches AddFile
    And content-based deduplication is used
    And duplicate content shares data blocks

  @REQ-FILEMGMT-205 @happy
  Scenario: Content-based deduplication
    Given two files with identical content in memory
    When both files added via AddFileFromMemory
    Then content is stored once
    And both FileEntry objects reference same data block
    And data block offset is shared
    And package size is optimized

  @REQ-FILEMGMT-205 @happy
  Scenario: Shared data blocks
    Given multiple AddFileFromMemory calls with same content
    When files are added to package
    Then single data block is created
    And all FileEntry objects share data block
    And DataOffset points to same location
    And deduplication reduces package size

  @REQ-FILEMGMT-205 @happy
  Scenario: Deduplication across AddFile and AddFileFromMemory
    Given file added via AddFile from filesystem
    And identical content in memory
    When AddFileFromMemory is called with identical content
    Then content is deduplicated
    And memory-added file shares data with filesystem-added file
    And deduplication works across add methods

  @REQ-FILEMGMT-205 @happy
  Scenario: Hash-based deduplication detection
    Given byte data with specific content hash
    When AddFileFromMemory is called
    Then content hash is computed
    And hash is compared against existing data blocks
    And matching hash triggers deduplication
    And new data block only if hash differs

  @REQ-FILEMGMT-205 @happy
  Scenario: Deduplication preserves unique file metadata
    Given two files with identical content
    When both added via AddFileFromMemory
    Then separate FileEntry objects created
    And each FileEntry has unique FileID
    And each FileEntry can have different paths
    And each FileEntry can have different metadata
    And only data block is shared

  @REQ-FILEMGMT-205 @happy
  Scenario: Compression and deduplication interaction
    Given identical content added with different compression
    When AddFileFromMemory is called with compression options
    Then deduplication occurs before compression
    And uncompressed content hash is used for deduplication
    And compressed storage may differ per file
    And deduplication based on raw content

  @REQ-FILEMGMT-205 @happy
  Scenario: Deduplication efficiency for memory operations
    Given many in-memory files with duplicate content
    When files added via AddFileFromMemory
    Then significant space savings from deduplication
    And package size much smaller than total content size
    And deduplication works efficiently with memory operations
