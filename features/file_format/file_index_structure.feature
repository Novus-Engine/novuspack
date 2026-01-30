@domain:file_format @m2 @REQ-FILEFMT-080 @spec(package_file_format.md#61-file-index-structure)
Feature: File Index Structure
    As a package format implementer
    I want a defined file index structure
    So that file entries are organized and accessible efficiently

    Background:
        Given the file index structure organizes all file entries
        And the structure supports efficient lookup and traversal

    @REQ-FILEFMT-080 @happy
    Scenario: File index stores file entry count
        When a file index is created
        Then it MUST contain a field for the total number of file entries
        And the count MUST match the actual number of entries stored

    @REQ-FILEFMT-080 @happy
    Scenario: File index provides sequential entry access
        When file entries are read from the index
        Then entries MUST be accessible in sequential order
        And the order MUST be deterministic and reproducible

    @REQ-FILEFMT-080 @happy
    Scenario: File index supports offset-based entry lookup
        When a file entry needs to be accessed
        Then the index MUST provide byte offset information
        And offset information MUST enable direct entry access without scanning

    @REQ-FILEFMT-080 @happy
    Scenario: File index NewFileIndex creates initialized structure
        When NewFileIndex is called
        Then it MUST return a FileIndex with zero values
        And all fields MUST be properly initialized for first use

    @REQ-FILEFMT-080 @happy
    Scenario: File index maintains integrity during updates
        When file entries are added or removed
        Then the index structure MUST remain valid
        And entry count MUST be updated accordingly

    @REQ-FILEFMT-080 @happy
    Scenario: File index binary structure has defined size
        When FileIndexBinary is serialized
        Then the base structure MUST be 16 bytes
        And total size MUST be 16 bytes plus entry references
        And each entry reference adds to the total size
        And the binary format MUST be well-defined and consistent

    @REQ-FILEFMT-080 @happy
    Scenario: File index entry references are fixed size
        When entry references are stored in the file index
        Then each entry reference MUST be exactly 16 bytes
        And entry reference size MUST be consistent across all entries
        And the fixed size MUST enable direct offset calculations
        And total index size MUST be predictable from entry count
