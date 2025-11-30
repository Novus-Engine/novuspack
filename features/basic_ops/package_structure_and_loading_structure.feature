@domain:basic_ops @m2 @REQ-API_BASIC-021 @REQ-API_BASIC-028 @spec(api_basic_operations.md#1-package-structure-and-loading)
Feature: Package Structure and Loading

  @REQ-API_BASIC-021 @happy
  Scenario: Package structure contains all required fields
    Given a NovusPack package structure
    When the package is examined
    Then package contains Info metadata field
    And package contains FileEntries array
    And package contains DirectoryEntries array
    And package contains SpecialFiles map indexed by file type
    And package contains IsOpen state field
    And package contains FilePath string field
    And package contains internal header, index, and fileHandle fields

  @REQ-API_BASIC-021 @happy
  Scenario: Package loading process initializes all components
    Given a NovusPack file on disk
    When package is opened
    Then package header is loaded and validated
    And package info metadata is loaded
    And file entries index is loaded
    And special metadata files are loaded
    And directory metadata is parsed from YAML
    And file-directory associations are updated
    And package IsOpen state is set to true

  @REQ-API_BASIC-021 @happy
  Scenario: Package structure supports on-demand file entry loading
    Given a NovusPack package structure
    When package is opened
    Then FileEntries array is initialized
    And file entries can be loaded on demand
    And package structure supports lazy loading of file content

  @REQ-API_BASIC-028 @happy
  Scenario: NewPackage creates empty package in memory
    Given a valid context
    When NewPackage is called
    Then a new Package instance is returned
    And package has default header values with magic number 0x4E56504B
    And package has format version 1
    And package has empty file index
    And package has empty comment
    And package IsOpen state is false
    And package exists only in memory with no file I/O

  @REQ-API_BASIC-028 @happy
  Scenario: NewPackage initializes package with current timestamp
    Given a valid context
    When NewPackage is called
    Then package creation timestamp is set to current time
    And package modification timestamp is initialized

  @REQ-API_BASIC-028 @happy
  Scenario: NewPackage usage follows standard pattern with defer cleanup
    Given a valid context
    When NewPackage is called and assigned to variable
    And defer package.Close() is set up
    Then package instance is ready for operations
    And package will be properly closed when function exits
    And package follows standard Go resource management pattern

  @REQ-API_BASIC-028 @happy
  Scenario: NewPackage requires Write function to persist to disk
    Given a package created with NewPackage
    When package operations are performed in memory
    Then package remains in memory only
    And package is not written to disk
    And Write, SafeWrite, or FastWrite must be called to persist package
