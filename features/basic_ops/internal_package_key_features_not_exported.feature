@domain:basic_ops @m2 @REQ-API_BASIC-190 @spec(api_basic_operations.md#21171-internal-package-key-features)
Feature: internal package key features

  @REQ-API_BASIC-190 @happy
  Scenario: internal helpers are not exported for external use
    Given helper functions implemented under an internal package
    When external consumers attempt to use those helpers
    Then helpers are not exported or importable externally
    And internal helpers remain accessible to implementation packages within the module
    And internal helpers are designed to support internal workflows
    And key features focus on implementation convenience over API stability
    And internal helpers respect the internal-only boundary consistently

