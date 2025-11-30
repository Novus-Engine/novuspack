@domain:compression @m2 @REQ-COMPR-140 @spec(api_package_compression.md#811-threadsafetynone)
Feature: ThreadSafetyNone Mode

  @REQ-COMPR-140 @happy
  Scenario: ThreadSafetyNone provides no thread safety guarantees
    Given a compression operation
    And ThreadSafetyNone mode is configured
    When operations are performed
    Then no thread safety guarantees are provided
    And operations should not be called concurrently
    And concurrent access may cause race conditions

  @REQ-COMPR-140 @happy
  Scenario: ThreadSafetyNone is appropriate for single-threaded usage
    Given a compression operation
    And single-threaded usage is required
    When ThreadSafetyNone is used
    Then no synchronization overhead is incurred
    And performance is optimal for single-threaded access
    And operations are not protected by locking mechanisms

  @REQ-COMPR-140 @happy
  Scenario: ThreadSafetyNone has no locking overhead
    Given a compression operation
    And ThreadSafetyNone mode is configured
    When operations are performed
    Then no locking mechanisms are used
    And no mutex overhead is incurred
    And operations execute without synchronization delay

  @REQ-COMPR-140 @error
  Scenario: ThreadSafetyNone causes race conditions with concurrent access
    Given a compression operation
    And ThreadSafetyNone mode is configured
    When multiple goroutines access the operation concurrently
    Then race conditions may occur
    And data corruption is possible
    And results are undefined

  @REQ-COMPR-140 @happy
  Scenario: ThreadSafetyNone is configurable via StreamConfig
    Given a compression StreamConfig
    When ThreadSafetyNone is set as the ThreadSafetyMode
    Then ThreadSafetyNone mode is active
    And no thread safety guarantees are provided
    And operations should not be called concurrently
