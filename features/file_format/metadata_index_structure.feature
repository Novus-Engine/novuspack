@domain:file_format @m2 @REQ-FILEFMT-078 @spec(package_file_format.md#51-metadata-index-structure)
Feature: Metadata Index Structure
    As a package format implementer
    I want a defined binary format for the metadata index
    So that metadata files can be located efficiently during package reading

    Background:
        Given the metadata index has a specific binary format
        And the format supports multiple metadata file entries

    @REQ-FILEFMT-078 @happy
    Scenario: Metadata index entry contains file type identifier
        When a metadata index entry is created
        Then the entry MUST include a file type identifier field
        And the file type MUST identify the special metadata file category

    @REQ-FILEFMT-078 @happy
    Scenario: Metadata index entry contains byte offset
        When a metadata index entry is created
        Then the entry MUST include a byte offset field
        And the offset MUST point to the start of the metadata file content

    @REQ-FILEFMT-078 @happy
    Scenario: Metadata index entry contains file size
        When a metadata index entry is created
        Then the entry MUST include a file size field
        And the size MUST represent the byte length of the metadata file

    @REQ-FILEFMT-078 @happy
    Scenario: Metadata index maintains entry order
        When multiple metadata files are indexed
        Then entries MUST be stored in a consistent order
        And the order MUST support efficient lookup operations

    @REQ-FILEFMT-078 @happy
    Scenario: Metadata index binary format is versioned
        When the metadata index format changes
        Then the package format version MUST be updated
        And readers MUST handle version differences appropriately
