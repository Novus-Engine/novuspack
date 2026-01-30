@domain:basic_ops @m2 @REQ-API_BASIC-162 @spec(api_basic_operations.md#2131-re-exported-types)
Feature: Re-exported types list

  @REQ-API_BASIC-162 @happy
  Scenario: Root package lists re-exported types available to consumers
    Given the root package API surface
    When consumers access types from the root package
    Then a defined list of re-exported types is available from the root package
    And the list reflects types intended for common consumer workflows
    And re-exported types match the behavior and contracts of their source types
    And re-exported types remain stable across API versions
    And documentation clearly identifies the re-exported type set

