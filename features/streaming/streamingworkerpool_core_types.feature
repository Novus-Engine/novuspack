@domain:streaming @m2 @REQ-STREAM-055 @spec(api_streaming.md#32-core-types)
Feature: StreamingWorkerPool core types define worker pool structures

  @REQ-STREAM-055 @happy
  Scenario: StreamingWorkerPool core types define structures
    Given StreamingWorkerPool core type definitions
    When worker pool structures are used
    Then core types define worker pool structures as specified
    And structure fields match the specification
    And the behavior matches the core types specification
    And type safety is preserved
