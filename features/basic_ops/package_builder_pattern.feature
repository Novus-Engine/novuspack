@domain:basic_ops @m1 @REQ-API_BASIC-008 @spec(api_basic_operations.md#44-package-builder-pattern)
Feature: Package builder pattern

  @happy
  Scenario: NewBuilder creates a new builder instance
    Given package builder functionality
    When NewBuilder is called
    Then a PackageBuilder instance is returned
    And builder is ready for configuration

  @happy
  Scenario: WithCompression sets compression type
    Given a PackageBuilder instance
    When WithCompression is called with a compression type
    Then builder returns itself for chaining
    And compression type is stored for build

  @happy
  Scenario: WithEncryption sets encryption type
    Given a PackageBuilder instance
    When WithEncryption is called with an encryption type
    Then builder returns itself for chaining
    And encryption type is stored for build

  @happy
  Scenario: WithMetadata sets package metadata
    Given a PackageBuilder instance
    When WithMetadata is called with metadata map
    Then builder returns itself for chaining
    And metadata is stored for build

  @happy
  Scenario: WithComment sets package comment
    Given a PackageBuilder instance
    When WithComment is called with a comment string
    Then builder returns itself for chaining
    And comment is stored for build

  @happy
  Scenario: WithVendorID sets vendor identifier
    Given a PackageBuilder instance
    When WithVendorID is called with a vendor ID
    Then builder returns itself for chaining
    And vendor ID is stored for build

  @happy
  Scenario: WithAppID sets application identifier
    Given a PackageBuilder instance
    When WithAppID is called with an app ID
    Then builder returns itself for chaining
    And app ID is stored for build

  @happy
  Scenario: Builder methods can be chained
    Given a PackageBuilder instance
    When multiple configuration methods are called in sequence
    Then each method returns the builder
    And all configurations are accumulated

  @happy
  Scenario: Build creates package with all configurations
    Given a PackageBuilder with multiple configurations
    When Build is called with context
    Then a Package instance is returned
    And package is created with all specified configurations
    And package is ready for use

  @error
  Scenario: Build fails with invalid configuration
    Given a PackageBuilder with invalid configuration
    When Build is called
    Then an appropriate error is returned
    And package is not created

  @error
  Scenario: Build respects context cancellation
    Given a PackageBuilder with configurations
    And a cancelled context
    When Build is called
    Then a structured context error is returned
    And package creation is cancelled
