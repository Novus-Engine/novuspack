@domain:compression @m2 @REQ-COMPR-141 @spec(api_package_compression.md#812-threadsafetyreadonly)
Feature: ThreadSafetyReadOnly Mode

  @REQ-COMPR-141 @happy
  Scenario: ThreadSafetyReadOnly allows concurrent read operations
    Given a compression operation
    And ThreadSafetyReadOnly mode is configured
    When read operations are performed concurrently
    Then read-only operations are safe for concurrent access
    And multiple goroutines can call read methods simultaneously
    And read operations are protected

  @REQ-COMPR-141 @happy
  Scenario: ThreadSafetyReadOnly protects read operations only
    Given a compression operation
    And ThreadSafetyReadOnly mode is configured
    When operations are performed
    Then read operations are protected
    And write operations are not protected
    And concurrent writes may cause issues

  @REQ-COMPR-141 @happy
  Scenario: ThreadSafetyReadOnly supports multiple simultaneous readers
    Given a compression operation
    And ThreadSafetyReadOnly mode is configured
    When multiple goroutines perform read operations simultaneously
    Then all read operations complete successfully
    And read operations do not block each other
    And read performance is optimized for concurrent access

  @REQ-COMPR-141 @error
  Scenario: ThreadSafetyReadOnly does not protect concurrent writes
    Given a compression operation
    And ThreadSafetyReadOnly mode is configured
    When write operations are performed concurrently
    Then write operations are not protected
    And race conditions may occur
    And data corruption is possible

  @REQ-COMPR-141 @happy
  Scenario: ThreadSafetyReadOnly is configurable via StreamConfig
    Given a compression StreamConfig
    When ThreadSafetyReadOnly is set as the ThreadSafetyMode
    Then ThreadSafetyReadOnly mode is active
    And read operations are safe for concurrent access
    And write operations are not protected
