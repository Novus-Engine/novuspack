@domain:streaming @m2 @REQ-STREAM-063 @spec(api_streaming.md#43-key-methods)
Feature: Streaming: Streaming Key Methods (Configuration)

  @REQ-STREAM-063 @happy
  Scenario: Streaming configuration key methods provide configuration management
    Given an open NovusPack package
    And a valid context
    And streaming configuration system
    When streaming configuration key methods are examined
    Then configuration methods provide streaming configuration management
    And configuration enables stream customization
    And configuration supports streaming operations

  @REQ-STREAM-063 @happy
  Scenario: Streaming configuration key methods support context integration
    Given an open NovusPack package
    And a valid context
    And streaming configuration
    When streaming configuration methods are used
    Then all methods accept context.Context
    And context supports cancellation
    And context supports timeout handling

  @REQ-STREAM-063 @happy
  Scenario: Streaming configuration key methods provide type-safe configuration
    Given an open NovusPack package
    And a valid context
    And streaming configuration system
    When streaming configuration methods are used
    Then type-safe configuration is provided
    And configuration ensures streaming correctness
    And configuration enables flexible streaming setup
