@domain:compression @m2 @REQ-COMPR-138 @spec(api_package_compression.md#8-concurrency-patterns-and-thread-safety)
Feature: Concurrency patterns and thread safety

  @REQ-COMPR-138 @happy
  Scenario: Thread safety guarantees are defined by ThreadSafetyMode
    Given compression operations with concurrency requirements
    When ThreadSafetyMode is configured
    Then thread safety guarantees are defined
    And guarantees match selected mode
    And mode determines concurrency behavior

  @REQ-COMPR-138 @happy
  Scenario: ThreadSafetyNone provides no thread safety guarantees
    Given compression operations with ThreadSafetyNone
    When operations are called concurrently
    Then no thread safety guarantees are provided
    And operations should not be called concurrently
    And concurrent access may cause issues

  @REQ-COMPR-138 @happy
  Scenario: ThreadSafetyReadOnly allows concurrent read operations
    Given compression operations with ThreadSafetyReadOnly
    When read operations are called concurrently
    Then read-only operations are safe for concurrent access
    And multiple goroutines can safely call read methods
    And read operations are protected

  @REQ-COMPR-138 @happy
  Scenario: ThreadSafetyConcurrent supports concurrent read/write operations
    Given compression operations with ThreadSafetyConcurrent
    When read and write operations are called concurrently
    Then concurrent read/write operations are supported
    And read-write mutex is used for optimal read performance
    And concurrent operations are protected

  @REQ-COMPR-138 @happy
  Scenario: ThreadSafetyFull provides complete synchronization
    Given compression operations with ThreadSafetyFull
    When operations are called concurrently
    Then full thread safety with complete synchronization is provided
    And all operations are protected by appropriate locking
    And maximum thread safety is guaranteed
