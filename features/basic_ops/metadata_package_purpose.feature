@domain:basic_ops @m2 @REQ-API_BASIC-169 @spec(api_basic_operations.md#217-metadata-package-purpose)
Feature: metadata package purpose

  @REQ-API_BASIC-169 @happy
  Scenario: metadata package purpose includes package information and path metadata
    Given the metadata subpackage
    When its purpose is evaluated
    Then the metadata package purpose includes package information structures
    And the metadata package purpose includes path metadata structures
    And the metadata package supports representing special metadata file content
    And the metadata package integrates with other packages via stable types
    And the package purpose aligns with the documented API organization

