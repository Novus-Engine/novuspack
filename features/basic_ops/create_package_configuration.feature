@domain:basic_ops @m1 @REQ-API_BASIC-006 @spec(api_basic_operations.md#42-create-method)
Feature: Create package configuration

  @happy
  Scenario: Create configures package for writing at valid path
    Given a new Package instance
    And an existing writable directory
    When Create is called with a valid path in that directory
    Then the package is configured for writing at the specified path
    And the package remains in memory
    And no file is created on disk

  @REQ-API_BASIC-020 @error
  Scenario: Create fails if target directory does not exist
    Given a new Package instance
    And a non-existent directory path
    When Create is called with path in that directory
    Then a structured validation error is returned
    And error indicates directory does not exist
    And package is not configured

  @REQ-API_BASIC-020 @error
  Scenario: Create fails if target directory is not writable
    Given a new Package instance
    And an existing read-only directory
    When Create is called with path in that directory
    Then a structured validation error is returned
    And error indicates directory is not writable
    And package is not configured

  @REQ-API_BASIC-018 @error
  Scenario: Create fails with invalid file path
    Given a new Package instance
    When Create is called with an invalid or malformed path
    Then a structured validation error is returned
    And package is not configured

  @REQ-API_BASIC-017 @REQ-API_BASIC-019 @error
  Scenario: Create respects context cancellation
    Given a new Package instance
    And a cancelled context
    When Create is called
    Then a structured context error is returned
    And error type is context cancellation
    And package is not configured

  @REQ-API_BASIC-020 @error
  Scenario: Create does not create parent directories
    Given a new Package instance
    And a path with non-existent parent directories
    When Create is called with that path
    Then a structured validation error is returned
    And error indicates parent directory does not exist
    And no directories are created

  @REQ-API_BASIC-006 @happy
  Scenario: Create initializes package structure in memory
    Given a new Package instance
    And an existing writable directory
    When Create is called with a valid path
    Then package header structure is configured
    And package metadata is initialized to default values
    And basic package structure is set up
    And target path is stored for later writing
    And package remains in memory
