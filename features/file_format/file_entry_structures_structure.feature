@domain:file_format @m2 @REQ-FILEFMT-045 @REQ-FILEFMT-047 @REQ-FILEFMT-053 @spec(package_file_format.md#4-file-entries-and-data-section)
Feature: File Entry Structures

  @REQ-FILEFMT-045 @happy
  Scenario: File entries and data section define file storage structure
    Given a NovusPack package
    When file entries and data section is examined
    Then file storage structure is defined
    And section contains interleaved file entries and their data
    And each file entry immediately precedes its related data

  @REQ-FILEFMT-045 @happy
  Scenario: File entry structure is 64-byte binary format plus extended data
    Given a NovusPack package
    And a file entry
    When file entry structure is examined
    Then file entry is 64-byte binary format
    And extended data (paths, hashes, optional data) follows
    And structure enables efficient streaming and processing

  @REQ-FILEFMT-045 @happy
  Scenario: File data follows file entry immediately
    Given a NovusPack package
    And file entries are present
    When file entry layout is examined
    Then file data follows file entry immediately
    And interleaved layout is: Entry 1 => Data 1 => Entry 2 => Data 2
    And layout enables efficient processing

  @REQ-FILEFMT-045 @happy
  Scenario: Variable length structure supports different content sizes
    Given a NovusPack package
    And file entries have different sizes
    When file entry structure is examined
    Then structure length is variable based on content
    And variable length supports paths, hashes, and optional data
    And structure adapts to file entry complexity

  @REQ-FILEFMT-047 @happy
  Scenario: File entry structure requirements define entry format rules
    Given a file entry
    When file entry structure requirements are examined
    Then entry format rules are defined
    And structure supports unique file identification
    And structure supports version tracking and metadata

  @REQ-FILEFMT-053 @happy
  Scenario: Fixed structure provides 64-byte optimized structure
    Given a file entry
    When fixed structure is examined
    Then fixed structure is 64-byte optimized structure
    And structure is optimized for 8-byte alignment
    And structure minimizes padding and improves performance
