@domain:file_format @m1 @REQ-FILEFMT-019 @spec(package_file_format.md#5-file-index-section)
Feature: File index section complete specification

  @happy
  Scenario: File index section structure is correct
    Given a NovusPack package
    When file index section is examined
    Then FileIndexBinary header is present (16 bytes)
    And entry references follow header
    And entry count matches file count
    And index size matches IndexSize in header

  @happy
  Scenario: Entry references provide file metadata
    Given a NovusPack package with indexed files
    When entry references are examined
    Then each entry provides file path
    And each entry provides original size
    And each entry provides stored size
    And each entry provides compression type
    And each entry provides encryption type
    And each entry provides data offset

  @happy
  Scenario: Index enables efficient file lookup
    Given a NovusPack package
    When file is looked up by path
    Then index provides direct offset to file data
    And file can be accessed efficiently
    And index eliminates need to scan all entries

  @error
  Scenario: Invalid index structure is detected
    Given a corrupted NovusPack package
    When index structure is validated
    Then validation fails
    And structured corruption error is returned
    And error indicates index structure issue

  @happy
  Scenario: NewFileIndex creates index with zero values
    Given NewFileIndex is called
    Then a FileIndex is returned
    And FileIndex is in initialized state
    And file index all fields are zero or empty

  @happy
  Scenario: WriteTo serializes file index to binary format
    Given a FileIndex with values
    When file index WriteTo is called with writer
    Then file index is written to writer
    And header is written first (16 bytes)
    And entries follow header
    And written data matches file index content

  @happy
  Scenario: ReadFrom deserializes file index from binary format
    Given a reader with valid file index data
    When file index ReadFrom is called with reader
    Then file index is read from reader
    And file index fields match reader data
    And file index is valid

  @happy
  Scenario: File index round-trip serialization preserves all fields
    Given a FileIndex with all fields set
    When file index WriteTo is called with writer
    And ReadFrom is called with written data
    Then all file index fields are preserved
    And file index is valid

  @error
  Scenario: ReadFrom fails with incomplete header
    Given a reader with less than 16 bytes of file index data
    When file index ReadFrom is called with reader
    Then structured IO error is returned
    And error indicates read failure
