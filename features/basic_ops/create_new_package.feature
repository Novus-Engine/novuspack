@domain:basic_ops @m1 @REQ-API_BASIC-001 @spec(api_basic_operations.md#41-package-constructor)
Feature: Create New Package

  @REQ-API_BASIC-001 @happy
  Scenario: NewPackage creates empty valid container in memory
    Given a valid context
    When NewPackage is invoked
    Then a valid NovusPack container exists in memory
    And the container has default header values
    And the container has an empty file index
    And no file I/O operations are performed

  @REQ-API_BASIC-001 @happy
  Scenario: NewPackage initializes package with standard header
    Given a valid context
    When NewPackage is called
    Then package has standard NovusPack header structure
    And header magic number is set to 0x4E56504B
    And header version is set to 1
    And header timestamps are initialized

  @REQ-API_BASIC-001 @happy
  Scenario: NewPackage prepares package for file operations
    Given a valid context
    When NewPackage is called
    Then package is ready for file operations
    And package can accept files via AddFile
    And package structure is initialized for use

  @REQ-API_BASIC-001 @happy
  Scenario: NewPackage creates package in memory only
    Given a valid context
    When NewPackage is called
    Then package exists only in memory
    And no file is created on disk
    And package must be written using Write methods to persist

  @REQ-API_BASIC-017 @happy
  Scenario: NewPackage accepts context parameter
    Given a valid context
    When NewPackage is called with context
    Then context is accepted as parameter
    And context supports cancellation
    And context supports timeout handling

  @REQ-API_BASIC-017 @REQ-API_BASIC-019 @error
  Scenario: NewPackage respects context cancellation
    Given a cancelled context
    When NewPackage is called
    Then structured context error is returned
    And error type is context cancellation
    And error follows structured error format

  @REQ-API_BASIC-019 @error
  Scenario: NewPackage handles context timeout
    Given a context with timeout
    And operation exceeds timeout duration
    When NewPackage is called
    Then context timeout error is returned
    And error type is context timeout
