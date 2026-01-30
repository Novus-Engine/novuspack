@domain:basic_ops @m2 @REQ-API_BASIC-140 @spec(api_basic_operations.md#665-packageconfig-backward-compatibility)
Feature: PackageConfig backward compatibility defaults

  @REQ-API_BASIC-140 @happy
  Scenario: PackageConfig defaults preserve backward compatible behavior
    Given a package created without explicitly setting PackageConfig options
    When PackageConfig is applied with default values
    Then default values maintain existing behavior for older consumers
    And default path handling behavior matches the documented backward compatible mode
    And default AutoConvertToSymlinks behavior preserves prior expectations
    And upgrades do not silently change path handling semantics
    And explicit configuration is required to opt into new behaviors

