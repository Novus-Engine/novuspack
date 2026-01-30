@domain:file_format @m2 @REQ-FILEFMT-079 @spec(package_file_format.md#52-metadata-index-detection)
Feature: Metadata Index Detection
    As a package reader implementation
    I want to detect the presence and location of the metadata index
    So that special metadata files can be discovered without full package scanning

    Background:
        Given a NovusPack package file is being opened
        And the metadata index may or may not be present

    @REQ-FILEFMT-079 @happy
    Scenario: Detect metadata index presence during package open
        When the package is opened
        Then the reader MUST detect if a metadata index is present
        And detection MUST occur during header parsing

    @REQ-FILEFMT-079 @happy
    Scenario: Locate metadata index via package header flags
        When the package header is read
        Then the metadata index presence flag MUST indicate if index exists
        And the header MUST provide the byte offset to the metadata index

    @REQ-FILEFMT-079 @happy
    Scenario: Handle missing metadata index gracefully
        When a package does not have a metadata index
        Then the reader MUST NOT fail to open the package
        And metadata file discovery MUST fall back to file entry scanning

    @REQ-FILEFMT-079 @happy
    Scenario: Validate metadata index integrity
        When a metadata index is detected
        Then the index structure MUST be validated
        And corrupted index data MUST be detected and reported

    @REQ-FILEFMT-079 @happy
    Scenario: Cache metadata index for fast access
        When a metadata index is loaded
        Then the index data SHOULD be cached in memory
        And subsequent metadata lookups SHOULD use the cached index
