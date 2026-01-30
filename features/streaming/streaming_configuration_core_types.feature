@domain:streaming @m2 @REQ-STREAM-061 @spec(api_streaming.md#42-core-types)
Feature: Streaming configuration core types define configuration structures

  @REQ-STREAM-061 @happy
  Scenario: Streaming configuration core types define structures
    Given streaming configuration core type definitions
    When configuration structures are used
    Then core types define configuration structures as specified
    And structure fields match the specification
    And the behavior matches the core types specification
    And configuration is type-safe
