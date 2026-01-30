@domain:basic_ops @m2 @REQ-API_BASIC-111 @spec(api_basic_operations.md#21-module-path)
Feature: Module path standard import

  @REQ-API_BASIC-111 @happy
  Scenario: Module path defines the standard import path for all packages
    Given a consumer integrating NovusPack in Go
    When the module is imported
    Then the documented module path is used as the standard import path
    And subpackages are imported under that module path
    And example import statements follow the module path convention
    And the standard import path is consistent across documentation and code organization
    And tools and builds resolve packages using the documented module path

