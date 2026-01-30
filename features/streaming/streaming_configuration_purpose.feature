@domain:streaming @m2 @REQ-STREAM-060 @spec(api_streaming.md#41-purpose)
Feature: Streaming configuration purpose defines configuration interface

  @REQ-STREAM-060 @happy
  Scenario: Streaming configuration purpose defines interface
    Given streaming configuration for FileStream or BufferPool
    When configuration is applied
    Then the purpose defines configuration interface behavior
    And configuration options are applied as specified
    And the behavior matches the streaming configuration purpose specification
    And configuration is validated before use
