@domain:basic_ops @m1 @REQ-API_BASIC-001 @spec(api_basic_operations.md#41-package-constructor)
Feature: Create New Package

  @REQ-API_BASIC-001 @happy
  Scenario: NewPackage creates empty valid container in memory
    When NewPackage is invoked
    Then a valid NovusPack container exists in memory
    And the container has default header values
    And the container has an empty file index
    And no file I/O operations are performed

  @REQ-API_BASIC-001 @happy
  Scenario: NewPackage initializes package with standard header
    When NewPackage is called
    Then package has standard NovusPack header structure
    And header magic number is set to 0x4E56504B
    And header version is set to 1
    And header timestamps are initialized

  @REQ-API_BASIC-001 @happy
  Scenario: NewPackage prepares package for file operations
    When NewPackage is called
    Then package is ready for file operations
    And package can accept files via AddFile
    And package structure is initialized for use

  @REQ-API_BASIC-001 @happy
  Scenario: NewPackage creates package in memory only
    When NewPackage is called
    Then package exists only in memory
    And no file is created on disk
    And package must be written using Write methods to persist

  @REQ-API_BASIC-001 @error
  Scenario: NewPackage returns error on failure
    When NewPackage is called
    And package initialization fails
    Then a structured error is returned
    And error indicates initialization failure
    And no package instance is created
