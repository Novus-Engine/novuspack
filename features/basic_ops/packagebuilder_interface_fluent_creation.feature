@domain:basic_ops @m2 @REQ-API_BASIC-134 @spec(api_basic_operations.md#651-packagebuilder-interface)
Feature: PackageBuilder interface supports fluent creation

  @REQ-API_BASIC-134 @happy
  Scenario: PackageBuilder provides a fluent interface for complex configurations
    Given a package creation workflow requiring multiple configuration steps
    When a PackageBuilder is used
    Then the builder exposes a fluent interface for configuration
    And configuration steps can be chained for readability
    And builder operations map to supported configuration options
    And the builder produces a package ready for create or open operations
    And builder-based creation supports complex configurations consistently

