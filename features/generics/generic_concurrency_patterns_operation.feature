@skip @domain:generics @m2 @REQ-GEN-016 @spec(api_generics.md#18-generic-concurrency-patterns)
Feature: Generic Concurrency Patterns

# This feature captures generic concurrency patterns from the generics API specification.
# More detailed runnable scenarios live in dedicated generics feature files.

  @REQ-GEN-016 @architecture
  Scenario: WorkerPool provides a reusable pattern for concurrent processing
    Given a set of independent jobs that can run in parallel
    When a WorkerPool is configured with a maximum number of workers
    Then jobs are processed concurrently up to the configured parallelism
    And the pool supports coordinated shutdown

  @REQ-GEN-016 @constraint
  Scenario: Jobs carry context for cancellation and timeouts
    Given a Job submitted to a WorkerPool
    When the job Context is cancelled
    Then the worker stops processing the job as soon as practical
    And the job result reports cancellation to the caller
