@domain:basic_ops @m2 @REQ-API_BASIC-114 @spec(api_basic_operations.md#23-package-aliases)
Feature: Package aliases

  @REQ-API_BASIC-114 @happy
  Scenario: Package aliases avoid naming conflicts when importing multiple packages
    Given a Go file that imports multiple packages with colliding names
    When the consumer applies import aliases
    Then a package alias can be used to avoid naming conflicts
    And aliased imports keep call sites unambiguous
    And aliasing follows standard Go import syntax
    And alias guidance is consistent with documented import patterns
    And alias usage supports readable and maintainable consumer code

