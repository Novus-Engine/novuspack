@domain:generics @m2 @REQ-GEN-018 @spec(api_generics.md#19-generic-configuration-patterns)
Feature: Generic Configuration Patterns

  @REQ-GEN-018 @happy
  Scenario: Config provides type-safe configuration for any data type
    Given a generic Config type
    When Config is used for configuration
    Then type-safe configuration is provided
    And Config supports Option fields for optional values
    And Config is flexible and type-safe

  @REQ-GEN-018 @happy
  Scenario: Config supports ChunkSize, MaxMemoryUsage, CompressionLevel options
    Given a Config instance
    When Option fields are set
    Then ChunkSize option can be configured
    And MaxMemoryUsage option can be configured
    And CompressionLevel option can be configured
    And type-safe configuration is provided

  @REQ-GEN-018 @happy
  Scenario: ConfigBuilder provides fluent configuration building
    Given a ConfigBuilder instance
    When builder methods are called
    Then fluent configuration interface is provided
    And WithChunkSize configures chunk size
    And WithMemoryUsage configures memory usage
    And WithCompressionLevel configures compression level
    And WithStrategy configures strategy

  @REQ-GEN-018 @happy
  Scenario: ConfigBuilder Build creates Config instance
    Given a ConfigBuilder with configured options
    When Build is called
    Then Config instance is created
    And Config contains configured options
    And type-safe configuration is provided

  @REQ-GEN-018 @happy
  Scenario: Generic configuration patterns support type-safe configuration management
    Given generic configuration patterns
    When configuration patterns are used with different types
    Then type safety is enforced at compile time
    And configuration patterns work with any type
    And generic patterns are reusable
