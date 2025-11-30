@domain:basic_ops @m1 @REQ-API_BASIC-004 @spec(api_basic_operations.md#11-package-structure)
Feature: Package structure and loading process

  @happy
  Scenario: Package structure is correctly initialized
    Given a new Package instance
    When the package is created
    Then Info is initialized to a PackageInfo structure
    And FileEntries is initialized to an empty slice
    And DirectoryEntries is initialized to an empty slice
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
    Then package info is loaded (comment, VendorID, AppID)
    Then file entries are loaded from file index
    Then special metadata files are loaded (file types 65000-65535)
    Then directory metadata is loaded from special files
    Then file-directory associations are updated
    And IsOpen is set to true

  @happy
  Scenario: loadSpecialMetadataFiles processes special file types
    Given a NovusPack package with special metadata files
    When loadSpecialMetadataFiles is called
    Then all files with types 65000-65535 are identified
    And special files are stored in SpecialFiles map
    And fileType maps to FileEntry in SpecialFiles

  @happy
  Scenario: loadDirectoryMetadata parses directory structure
    Given a NovusPack package with directory metadata files
    When loadDirectoryMetadata is called
    Then directory metadata is parsed from YAML special files
    And DirectoryEntries are populated
    And directory hierarchy is established

  @happy
  Scenario: updateFileDirectoryAssociations links files to directories
    Given a NovusPack package with files and directory metadata
    When updateFileDirectoryAssociations is called
    Then each file is linked to its parent directory
    And ParentDirectory pointer is set correctly
    And file-directory relationships are established

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
