@domain:basic_ops @m2 @REQ-API_BASIC-019 @REQ-API_BASIC-021 @spec(api_basic_operations.md#442-example-usage)
Feature: PackageBuilder example usage

  @REQ-API_BASIC-019 @REQ-API_BASIC-021 @happy
  Scenario: PackageBuilder example demonstrates fluent interface usage
    Given package needs to be created with complex configuration
    When PackageBuilder is used with fluent methods
    Then WithCompression method configures compression
    And WithEncryption method configures encryption
    And WithMetadata method configures metadata
    And WithComment method sets comment
    And WithVendorID method sets VendorID
    And WithAppID method sets AppID
    And Build method creates configured package

  @REQ-API_BASIC-019 @REQ-API_BASIC-021 @happy
  Scenario: PackageBuilder example shows builder pattern workflow
    Given PackageBuilder pattern
    When builder methods are chained
    Then methods return builder for chaining
    And configuration is built incrementally
    And final Build creates package

  @REQ-API_BASIC-019 @REQ-API_BASIC-021 @happy
  Scenario: PackageBuilder example improves code readability
    Given complex package configuration needs
    When PackageBuilder is used
    Then code is more readable than parameter-heavy methods
    And configuration options are clear
    And builder pattern simplifies complex setups
