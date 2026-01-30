@domain:basic_ops @m2 @REQ-API_BASIC-138 @spec(api_basic_operations.md#663-packageconfig-purpose)
Feature: PackageConfig purpose during file addition

  @REQ-API_BASIC-138 @happy
  Scenario: PackageConfig exists to control path handling behavior during file addition operations
    Given file addition operations for a package
    When package-level configuration is evaluated
    Then PackageConfig purpose is to control path handling behavior
    And configuration influences how stored paths and base paths are derived
    And configuration supports predictable results across different file sources
    And configuration enables consistent defaults for common workflows
    And configuration is applied uniformly across file addition APIs

