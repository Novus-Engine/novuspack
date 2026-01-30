@domain:basic_ops @m2 @REQ-API_BASIC-160 @spec(api_basic_operations.md#2121-root-package-key-features)
Feature: Root package key features

  @REQ-API_BASIC-160 @happy
  Scenario: Root package key features include re-exported types and a unified import path
    Given a consumer writing Go code against NovusPack
    When choosing an import strategy
    Then a unified import path is available via the root package
    And key types are re-exported to reduce required imports
    And re-exported types enable ergonomic access to common functionality
    And the root package surfaces remain coherent and discoverable
    And key features support both simple and advanced use cases

