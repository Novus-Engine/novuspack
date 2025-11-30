@domain:basic_ops @m2 @REQ-API_BASIC-026 @spec(api_basic_operations.md#4-package-creation)
Feature: Package Creation

  @REQ-API_BASIC-026 @happy
  Scenario: Package creation uses NewPackage constructor
    Given a valid context
    When NewPackage is called
    Then new empty package is created
    And package has default header values
    And package has empty file index
    And package has empty comment
    And package exists only in memory

  @REQ-API_BASIC-026 @happy
  Scenario: Package creation uses Create method
    Given a package created with NewPackage
    And a valid context
    And a valid file path
    When Create is called with path
    Then package is configured for creation
    And target path is stored
    And package structure is prepared
    And package remains in memory until Write is called

  @REQ-API_BASIC-026 @happy
  Scenario: Package creation uses CreateWithOptions method
    Given a package created with NewPackage
    And a valid context
    And a valid file path
    And package creation options
    When CreateWithOptions is called
    Then package is configured with options
    And options are applied to package
    And package structure is prepared with options

  @REQ-API_BASIC-026 @happy
  Scenario: Package creation supports builder pattern
    Given a valid context
    When package builder is used
    Then builder provides fluent interface
    And builder allows complex configuration
    And builder creates package when Build is called

  @REQ-API_BASIC-026 @happy
  Scenario: Package creation requires Write to persist
    Given a package created and configured
    When package operations are performed
    Then package remains in memory
    And package is not written to disk
    And Write, SafeWrite, or FastWrite must be called to persist

  @REQ-API_BASIC-026 @error
  Scenario: Package creation validates target directory exists
    Given a package created with NewPackage
    And a valid context
    And a path with non-existent directory
    When Create is called
    Then validation error is returned
    And error indicates directory does not exist
    And package is not configured

  @REQ-API_BASIC-026 @error
  Scenario: Package creation validates target directory is writable
    Given a package created with NewPackage
    And a valid context
    And a path with read-only directory
    When Create is called
    Then validation error is returned
    And error indicates directory is not writable
    And package is not configured
