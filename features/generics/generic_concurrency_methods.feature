@domain:generics @m2 @REQ-GEN-017 @spec(api_generics.md#18-generic-concurrency-methods)
Feature: Generic Concurrency Methods

  @REQ-GEN-017 @happy
  Scenario: WorkerPool Start initializes and starts worker pool
    Given a WorkerPool instance
    And a valid context
    When Start is called with context
    Then worker pool is initialized
    And worker pool is started
    And workers are ready to process jobs
    And error is nil on success

  @REQ-GEN-017 @happy
  Scenario: WorkerPool Stop gracefully shuts down worker pool
    Given a running WorkerPool instance
    And a valid context
    When Stop is called with context
    Then worker pool is gracefully shut down
    And workers complete current jobs
    And resources are cleaned up
    And error is nil on success

  @REQ-GEN-017 @happy
  Scenario: WorkerPool SubmitJob submits job for processing
    Given a running WorkerPool instance
    And a Job with data
    And a valid context
    When SubmitJob is called with context and job
    Then job is submitted to worker pool
    And job is processed by available worker
    And error is nil on success

  @REQ-GEN-017 @happy
  Scenario: WorkerPool GetWorkerStats returns worker statistics
    Given a running WorkerPool instance
    When GetWorkerStats is called
    Then worker pool statistics are returned
    And statistics include worker count and status
    And statistics provide operational insights

  @REQ-GEN-017 @happy
  Scenario: ProcessConcurrently processes items concurrently
    Given multiple items to process
    And a Strategy implementation
    And a ConcurrencyConfig
    And a valid context
    When ProcessConcurrently is called with context, items, processor, and config
    Then items are processed concurrently
    And Result slice is returned
    And concurrent processing is type-safe
    And error is nil on success
