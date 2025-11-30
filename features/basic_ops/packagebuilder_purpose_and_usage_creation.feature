@domain:basic_ops @m2 @REQ-API_BASIC-037 @spec(api_basic_operations.md#441-purpose)
Feature: PackageBuilder Purpose and Usage

  @REQ-API_BASIC-037 @happy
  Scenario: PackageBuilder provides fluent interface for package creation
    Given a valid context
    When PackageBuilder is used
    Then builder provides fluent interface
    And builder methods can be chained
    And builder improves code readability

  @REQ-API_BASIC-037 @happy
  Scenario: PackageBuilder supports complex package configurations
    Given a valid context
    And complex package configuration requirements
    When PackageBuilder is used
    Then builder supports compression configuration
    And builder supports encryption configuration
    And builder supports metadata configuration
    And builder supports comment, VendorID, and AppID configuration

  @REQ-API_BASIC-037 @happy
  Scenario: PackageBuilder reduces parameter complexity
    Given a valid context
    And package with many configuration options
    When PackageBuilder is used instead of direct method calls
    Then parameter complexity is reduced
    And code is more readable
    And configuration is easier to manage

  @REQ-API_BASIC-037 @happy
  Scenario: PackageBuilder creates package when Build is called
    Given a PackageBuilder instance
    And builder is configured with options
    And a valid context
    When Build is called
    Then package is created with specified configuration
    And package follows builder configuration
    And package is ready for use

  @REQ-API_BASIC-037 @happy
  Scenario: PackageBuilder example demonstrates usage pattern
    Given a code example using PackageBuilder
    When example is examined
    Then builder methods are chained together
    And WithCompression, WithEncryption, WithComment are used
    And Build is called to create package
    And pattern improves code readability
