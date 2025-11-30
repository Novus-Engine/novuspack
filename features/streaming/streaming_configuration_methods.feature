@domain:streaming @m2 @REQ-STREAM-012 @spec(api_streaming.md#2-buffer-management-system)
Feature: Streaming configuration methods

  @happy
  Scenario: DefaultBufferConfig returns default configuration
    When DefaultBufferConfig is called
    Then default BufferConfig is returned
    And configuration values are set to defaults
    And configuration is valid

  @happy
  Scenario: NewStreamingConfigBuilder creates configuration builder
    When NewStreamingConfigBuilder is called
    Then StreamingConfigBuilder is created
    And builder is ready for configuration
    And builder supports fluent API

  @happy
  Scenario: CreateStreamingConfig creates streaming configuration
    Given configuration parameters
    When CreateStreamingConfig is called with parameters
    Then StreamingConfig is created
    And configuration matches parameters
    And configuration is valid

  @happy
  Scenario: ValidateStreamingConfig validates configuration
    Given a StreamingConfig
    When ValidateStreamingConfig is called with config
    Then configuration is validated
    And validation result indicates validity
    And invalid configurations are detected

  @happy
  Scenario: GetStreamingConfigDefaults returns default values
    When GetStreamingConfigDefaults is called
    Then default streaming configuration values are returned
    And values are reasonable
    And values are documented

  @error
  Scenario: ValidateStreamingConfig fails with invalid configuration
    Given an invalid StreamingConfig
    When ValidateStreamingConfig is called
    Then structured validation error is returned
    And error indicates configuration issue
