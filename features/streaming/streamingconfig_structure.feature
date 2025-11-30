@domain:streaming @m2 @REQ-STREAM-062 @spec(api_streaming.md#421-streamingconfig-struct)
Feature: StreamingConfig Structure

  @REQ-STREAM-062 @happy
  Scenario: StreamingConfig struct provides streaming configuration structure
    Given a NovusPack package
    When StreamingConfig struct is used
    Then struct extends Config for streaming-specific settings
    And struct contains StreamBufferSize option for buffer size
    And struct contains ChunkProcessingMode option for chunk processing
    And struct contains MaxStreamsPerWorker option for stream limits
    And struct contains StreamTimeout option for operation timeout

  @REQ-STREAM-062 @happy
  Scenario: StreamingConfig extends generic Config
    Given a NovusPack package
    And a StreamingConfig
    When configuration structure is examined
    Then StreamingConfig embeds Config with FileStream type parameter
    And extension provides streaming-specific configuration
    And generic base enables type-safe configuration patterns

  @REQ-STREAM-062 @happy
  Scenario: StreamBufferSize option configures stream buffer size
    Given a NovusPack package
    And a StreamingConfig
    When StreamBufferSize option is set
    Then buffer size for stream operations is configured
    And option uses Option[int] for optional configuration
    And buffer size affects stream performance

  @REQ-STREAM-062 @happy
  Scenario: ChunkProcessingMode option configures chunk processing
    Given a NovusPack package
    And a StreamingConfig
    When ChunkProcessingMode option is set
    Then how chunks are processed concurrently is configured
    And option uses Option[ChunkMode] for processing mode
    And mode determines concurrent chunk handling

  @REQ-STREAM-062 @happy
  Scenario: MaxStreamsPerWorker option configures stream limits
    Given a NovusPack package
    And a StreamingConfig
    When MaxStreamsPerWorker option is set
    Then maximum streams per worker is configured
    And option uses Option[int] for limit configuration
    And limit controls concurrent stream processing

  @REQ-STREAM-062 @happy
  Scenario: StreamTimeout option configures operation timeout
    Given a NovusPack package
    And a StreamingConfig
    When StreamTimeout option is set
    Then timeout for stream operations is configured
    And option uses Option[time.Duration] for timeout
    And timeout prevents indefinite stream operations

  @REQ-STREAM-062 @happy
  Scenario: StreamingConfigBuilder enables fluent configuration
    Given a NovusPack package
    When StreamingConfigBuilder is used
    Then builder provides fluent API for configuration
    And builder methods return builder for chaining
    And Build method creates final StreamingConfig

  @REQ-STREAM-062 @error
  Scenario: StreamingConfig struct validates configuration values
    Given a NovusPack package
    And a StreamingConfig with invalid values
    When StreamingConfig is validated
    Then validation error is returned
    And error indicates invalid configuration field
    And error follows structured error format
