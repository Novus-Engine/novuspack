@domain:basic_ops @m2 @REQ-API_BASIC-058 @spec(api_basic_operations.md#732-getinfo-example-usage)
Feature: GetInfo example usage

  @REQ-API_BASIC-058 @happy
  Scenario: GetInfo example demonstrates package information retrieval
    Given an open NovusPack package
    When GetInfo is called
    Then PackageInfo structure is returned
    And file count information is available
    And package version information is available
    And comprehensive package information is retrieved

  @REQ-API_BASIC-058 @happy
  Scenario: GetInfo example shows accessing package information
    Given GetInfo returns PackageInfo
    When package information is accessed
    Then file count can be retrieved
    And package version can be retrieved
    And VendorID and AppID are accessible
    And package metadata is available
    And signature information is accessible
    And security and compression status is available

  @REQ-API_BASIC-058 @happy
  Scenario: GetInfo example demonstrates information display
    Given PackageInfo from GetInfo
    When information is formatted for display
    Then file count is displayed
    And package version is displayed
    And package information is presented clearly
    And example demonstrates usage pattern
