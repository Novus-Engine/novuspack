@domain:basic_ops @m2 @REQ-API_BASIC-135 @spec(api_basic_operations.md#652-newbuilder-function)
Feature: NewBuilder creates a new builder

  @REQ-API_BASIC-135 @happy
  Scenario: NewBuilder returns a new PackageBuilder instance
    Given a need to create a package with a builder
    When NewBuilder is called
    Then a new PackageBuilder instance is returned
    And the builder starts from default configuration state
    And configuration methods can be applied before creating the package
    And builder instances are independent across calls
    And the builder can be used to create packages with options

