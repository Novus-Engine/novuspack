@domain:basic_ops @m2 @REQ-API_BASIC-168 @spec(api_basic_operations.md#2161-metadata-package-metadata)
Feature: metadata package provides metadata structures

  @REQ-API_BASIC-168 @happy
  Scenario: metadata package provides metadata structures
    Given the metadata subpackage
    When metadata structures are needed
    Then the metadata package provides metadata structures
    And metadata structures support package information and path metadata models
    And metadata structures are used consistently across the API surface
    And metadata structures align with tech spec definitions
    And metadata types support serialization and validation workflows

