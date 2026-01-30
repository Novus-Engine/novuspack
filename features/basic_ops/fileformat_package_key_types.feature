@domain:basic_ops @m2 @REQ-API_BASIC-165 @spec(api_basic_operations.md#215-file-format-package-key-types)
Feature: fileformat package key types

  @REQ-API_BASIC-165 @happy
  Scenario: fileformat defines key types for representing package file format structures
    Given the fileformat package API
    When key types are used to represent the package file format
    Then PackageHeader is defined as a key type
    And FileIndex is defined as a key type
    And IndexEntry is defined as a key type
    And key types map to the documented binary layout and semantics
    And key types support higher-level operations through stable contracts

