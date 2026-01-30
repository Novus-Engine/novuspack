@domain:basic_ops @m2 @REQ-API_BASIC-163 @spec(api_basic_operations.md#214-file-format-package-fileformat)
Feature: fileformat package provides binary file format structures

  @REQ-API_BASIC-163 @happy
  Scenario: fileformat package provides binary file format structures
    Given the fileformat subpackage
    When file format structures are needed by consumers or internal packages
    Then fileformat provides binary file format structures
    And structures represent on-disk layout concepts like headers and indexes
    And structures support parsing and writing workflows
    And fileformat is organized as a dedicated package for binary format concerns
    And fileformat types are referenced consistently across the API

