@domain:basic_ops @m2 @REQ-API_BASIC-158 @spec(api_basic_operations.md#2111-root-package-novuspack)
Feature: Root package provides main public API

  @REQ-API_BASIC-158 @happy
  Scenario: Root package novuspack provides the main public API and re-exports types
    Given the root package import
    When consumers use the root package for typical workflows
    Then the root package exposes the main public API surface
    And commonly used types are re-exported from the root package
    And consumers can avoid importing many subpackages for standard usage
    And the root package remains stable across minor versions
    And re-exported types remain consistent with their source packages

