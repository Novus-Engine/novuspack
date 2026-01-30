@domain:basic_ops @m1 @REQ-API_BASIC-004 @spec(api_basic_operations.md#11-package-structure)
Feature: Package structure and loading process

  @happy
  Scenario: Package structure is correctly initialized
    Given a new Package instance
    When the package is created
    Then Info is initialized to a PackageInfo structure
    And FileEntries is initialized to an empty slice
    And PathMetadataEntries is initialized to an empty slice
    And SpecialFiles is initialized to an empty map
    And IsOpen is set to false
    And FilePath is empty
    And header, index, and fileHandle are initialized

  @happy
  Scenario: Package loading process follows correct sequence
    Given an existing NovusPack package file
    When the package is opened
    Then package header is loaded first
    And package header magic number and version are validated
    Then package info is synchronized from header using PackageInfo.FromHeader
    Then file entries are loaded from file index
    Then special metadata files are loaded (file types 65000-65535)
    Then path metadata is loaded from special files
    Then file-path associations are updated
    And IsOpen is set to true

  @happy
  Scenario: loadSpecialMetadataFiles processes special file types
    Given a NovusPack package with special metadata files
    When loadSpecialMetadataFiles is called
    Then all files with types 65000-65535 are identified
    And special files are stored in SpecialFiles map
    And fileType maps to FileEntry in SpecialFiles

  @happy
  Scenario: loadPathMetadata parses path structure
    Given a NovusPack package with path metadata files
    When loadPathMetadata is called
    Then path metadata is parsed from YAML special files
    And PathMetadataEntries are populated
    And path hierarchy is established

  @happy
  Scenario: updateFilePathAssociations links files to path metadata
    Given a NovusPack package with files and path metadata
    When updateFilePathAssociations is called
    Then each file path is linked to its corresponding PathMetadataEntry
    And PathMetadataEntries map is populated correctly
    And file-path relationships are established

  @error
  Scenario: Package loading fails if header is invalid
    Given a corrupted NovusPack package file with invalid header
    When the package is opened
    Then a structured invalid format error is returned
    And package loading stops at header validation

  @error
  Scenario: Package loading fails if file index is corrupted
    Given a NovusPack package with corrupted file index
    When the package is opened
    Then a structured corruption error is returned
    And file entries are not loaded
