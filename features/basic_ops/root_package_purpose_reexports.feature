@domain:basic_ops @m2 @REQ-API_BASIC-159 @spec(api_basic_operations.md#212-root-package-purpose)
Feature: Root package purpose

  @REQ-API_BASIC-159 @happy
  Scenario: Root package purpose is to provide main API access and type re-exports
    Given the documented Go API organization
    When the root package purpose is evaluated
    Then the root package exists to provide a main API entry point
    And the root package re-exports key types for consumer convenience
    And consumers can rely on a single import path for common operations
    And the root package purpose aligns with the subpackage architecture
    And the root package avoids leaking internal-only implementation details

