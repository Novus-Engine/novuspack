@domain:basic_ops @m2 @REQ-API_BASIC-046 @spec(api_basic_operations.md#6-package-closing)
Feature: Package Closing

  @REQ-API_BASIC-046 @happy
  Scenario: Close closes package file handle
    Given an open NovusPack package
    When Close is called
    Then package file handle is closed
    And file handle resources are released
    And file is no longer accessible through package

  @REQ-API_BASIC-046 @happy
  Scenario: Close releases memory buffers and caches
    Given an open NovusPack package
    And package has loaded memory buffers
    When Close is called
    Then memory buffers are released
    And cached data is cleared
    And memory resources are freed

  @REQ-API_BASIC-046 @happy
  Scenario: Close clears package state and metadata
    Given an open NovusPack package
    And package has loaded state and metadata
    When Close is called
    Then package state is cleared
    And package metadata is cleared
    And package IsOpen state is set to false

  @REQ-API_BASIC-046 @happy
  Scenario: Close performs necessary cleanup operations
    Given an open NovusPack package
    And package has active operations
    When Close is called
    Then cleanup operations are performed
    And all resources are properly released
    And package is in clean state

  @REQ-API_BASIC-046 @happy
  Scenario: Close does not modify package file
    Given an open NovusPack package
    And package has unsaved changes
    When Close is called
    Then package file is not modified
    And unsaved changes are not persisted
    And Write methods must be used to save changes

  @REQ-API_BASIC-046 @error
  Scenario: Close returns I/O error on file system errors
    Given an open NovusPack package
    And file system error occurs during closing
    When Close is called
    Then I/O error is returned
    And error indicates file system issue
    And error follows structured error format

  @REQ-API_BASIC-046 @error
  Scenario: Close returns validation error when package not open
    Given a NovusPack package that is not open
    When Close is called
    Then validation error is returned
    And error indicates package is not currently open
    And error follows structured error format
