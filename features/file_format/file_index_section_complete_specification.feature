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
