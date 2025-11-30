@domain:compression @m2 @REQ-COMPR-139 @spec(api_package_compression.md#81-thread-safety-guarantees)
Feature: Thread Safety Guarantees

  @REQ-COMPR-139 @happy
  Scenario: ThreadSafetyNone provides no thread safety guarantees
    Given a compression operation
    And ThreadSafetyNone mode is configured
    When operations are performed
    Then no thread safety guarantees are provided
    And operations should not be called concurrently
    And concurrent access may cause race conditions

  @REQ-COMPR-139 @happy
  Scenario: ThreadSafetyReadOnly allows concurrent read operations
    Given a compression operation
    And ThreadSafetyReadOnly mode is configured
    When read operations are performed concurrently
    Then read-only operations are safe for concurrent access
    And multiple goroutines can call read methods simultaneously
    And read operations are protected

  @REQ-COMPR-139 @happy
  Scenario: ThreadSafetyConcurrent supports concurrent read and write operations
    Given a compression operation
    And ThreadSafetyConcurrent mode is configured
    When read and write operations are performed concurrently
    Then concurrent read/write operations are supported
    And read-write mutex provides optimal read performance
    And write operations are synchronized

  @REQ-COMPR-139 @happy
  Scenario: ThreadSafetyFull provides complete synchronization
    Given a compression operation
    And ThreadSafetyFull mode is configured
    When operations are performed concurrently
    Then full thread safety is provided
    And all operations are protected by locking mechanisms
    And complete synchronization is ensured

  @REQ-COMPR-139 @happy
  Scenario: Thread safety mode is configurable via StreamConfig
    Given a compression operation
    And a StreamConfig with ThreadSafetyMode setting
    When compression is performed
    Then configured thread safety mode is applied
    And mode determines level of synchronization
    And mode affects concurrent access behavior
