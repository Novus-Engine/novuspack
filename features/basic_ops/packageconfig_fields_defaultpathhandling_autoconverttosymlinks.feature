@domain:basic_ops @m2 @REQ-API_BASIC-139 @spec(api_basic_operations.md#664-packageconfig-fields)
Feature: PackageConfig fields for default path handling and auto-convert-to-symlinks

  @REQ-API_BASIC-139 @happy
  Scenario: PackageConfig exposes DefaultPathHandling and AutoConvertToSymlinks configuration options
    Given a package configured for file addition operations
    When PackageConfig fields are set
    Then DefaultPathHandling controls default multi-path behavior
    And AutoConvertToSymlinks enables or disables automatic conversion to symlinks
    And field values are applied during file addition workflows
    And field defaults are applied when options are not explicitly set
    And configuration fields are reflected in in-memory package behavior

