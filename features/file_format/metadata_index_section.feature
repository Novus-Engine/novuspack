@domain:file_format @m2 @REQ-FILEFMT-077 @spec(package_file_format.md#5-metadata-index-section)
Feature: Metadata Index Section
    As a package format implementer
    I want a metadata index section in package files
    So that special metadata files can be discovered and accessed efficiently

    Background:
        Given a NovusPack package file
        And the metadata index section provides lookup information for special metadata files

    @REQ-FILEFMT-077 @happy
    Scenario: Metadata index section enables fast metadata file discovery
        When the package contains special metadata files
        Then the metadata index section MUST provide file type and offset information
        And readers MUST be able to locate metadata files without scanning all file entries

    @REQ-FILEFMT-077 @happy
    Scenario: Metadata index is optional
        When a package does not contain special metadata files
        Then the metadata index section MAY be omitted
        And the absence of metadata index MUST NOT cause errors

    @REQ-FILEFMT-077 @happy
    Scenario: Metadata index entries reference special files
        When a special metadata file is added to the package
        Then the metadata index MUST contain an entry for that file
        And the entry MUST include file type identifier and byte offset

    @REQ-FILEFMT-077 @happy
    Scenario: Metadata index supports multiple metadata files
        When a package contains multiple special metadata files
        Then the metadata index MUST contain entries for all special files
        And each entry MUST be uniquely identifiable by file type
