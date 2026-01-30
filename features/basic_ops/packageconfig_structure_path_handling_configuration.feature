@domain:basic_ops @m2 @REQ-API_BASIC-136 @spec(api_basic_operations.md#661-packageconfig-structure)
Feature: PackageConfig structure for path handling

  @REQ-API_BASIC-136 @happy
  Scenario: PackageConfig defines package-level configuration for path handling behavior
    Given a package configured for file addition operations
    When PackageConfig is provided as part of configuration
    Then PackageConfig structure defines package-level path handling configuration
    And configuration affects how paths are interpreted and stored
    And configuration supports consistent behavior across file addition operations
    And configuration is applied through creation options or builder patterns
    And the configuration is part of the in-memory package state

