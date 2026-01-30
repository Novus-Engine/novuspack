@domain:basic_ops @m2 @REQ-API_BASIC-175 @spec(api_basic_operations.md#2110-generics-package-key-types)
Feature: generics package key types

  @REQ-API_BASIC-175 @happy
  Scenario: generics package defines key types used across the API
    Given the generics package API
    When key types are referenced by other packages
    Then PathEntry is defined as a key type
    And Option is defined as a key type
    And Result is defined as a key type
    And Tag Strategy and Validator are defined as key types
    And key types support consistent composition across capability packages

