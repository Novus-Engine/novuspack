@domain:basic_ops @m2 @REQ-API_BASIC-131 @spec(api_basic_operations.md#411-format-constants)
Feature: Format constants define identifier, version, and header size

  @REQ-API_BASIC-131 @happy
  Scenario: Format constants provide package identifier, version, and header size values
    Given NovusPack package file format constants
    When a package header is created or validated
    Then constants define the package identifier value
    And constants define the package version values
    And constants define the fixed header size value
    And constant values are used consistently across read and write paths
    And constant values enable compatibility checks during open

