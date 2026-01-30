@domain:basic_ops @m2 @REQ-API_BASIC-188 @spec(api_basic_operations.md#21161-internal-package-internal)
Feature: internal package provides internal helpers

  @REQ-API_BASIC-188 @happy
  Scenario: internal package contains helper functions not exposed as public API
    Given internal helper logic needed by implementation packages
    When helpers are organized
    Then helper functions are placed in an internal-only package
    And helpers are not importable by external consumers
    And helpers support implementation without expanding the public API surface
    And internal helpers follow the documented internal package boundaries
    And internal helpers reduce duplication across implementation packages

