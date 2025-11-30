@domain:basic_ops @m2 @REQ-API_BASIC-030 @spec(api_basic_operations.md#422-create-behavior)
Feature: Create Behavior

  @REQ-API_BASIC-030 @happy
  Scenario: Create validates path and directory
    Given a package needs to be created
    When Create is called with path
    Then provided path is validated
    And path is checked for validity
    And target directory existence is verified
    And directory writability is checked

  @REQ-API_BASIC-030 @happy
  Scenario: Create configures package in memory
    Given a valid path for package creation
    When Create is called
    Then package structure is configured in memory
    And package header is initialized
    And package metadata is set to defaults
    And target path is stored for later writing

  @REQ-API_BASIC-030 @happy
  Scenario: Create does not write to disk
    Given a package creation operation
    When Create is called
    Then no file I/O operations are performed
    And package remains in memory only
    And package file is not created on disk
    And Write method must be called to save

  @REQ-API_BASIC-030 @happy
  Scenario: Create sets up basic package structure
    Given Create is called
    When package is created
    Then file entries structure is initialized
    And data sections are set up
    And package remains unsigned and uncompressed

  @REQ-API_BASIC-030 @error
  Scenario: Create returns error for invalid path
    Given an invalid or malformed file path
    When Create is called
    Then validation error is returned
    And error indicates invalid path
    And package is not created

  @REQ-API_BASIC-030 @error
  Scenario: Create returns error when directory does not exist
    Given a path with non-existent directory
    When Create is called
    Then validation error is returned
    And error indicates directory does not exist
    And package is not created

  @REQ-API_BASIC-030 @error
  Scenario: Create returns error when directory is not writable
    Given a path with non-writable directory
    When Create is called
    Then validation error is returned
    And error indicates directory is not writable
    And package is not created
