@domain:generics @m2 @REQ-GEN-020 @spec(api_generics.md#23-factory-functions)
Feature: Generics Factory Functions

  @REQ-GEN-020 @happy
  Scenario: NewTypedTag creates type-safe typed tag
    Given a tag key, value, and tag type
    When NewTypedTag is called with key, value, and tag type
    Then TypedTag instance is created
    And TypedTag is type-safe
    And TypedTag value type is enforced

  @REQ-GEN-020 @happy
  Scenario: NewBufferPool creates type-safe buffer pool
    Given a BufferConfig
    When NewBufferPool is called with config
    Then BufferPool instance is created
    And BufferPool is type-safe
    And BufferPool manages buffers efficiently

  @REQ-GEN-020 @happy
  Scenario: NewConfigBuilder creates type-safe config builder
    Given a generic type parameter
    When NewConfigBuilder is called
    Then ConfigBuilder instance is created
    And ConfigBuilder is type-safe
    And ConfigBuilder supports fluent configuration

  @REQ-GEN-020 @happy
  Scenario: Factory functions provide type-safe instance creation
    Given generic factory functions
    When factory functions are used with different types
    Then type safety is enforced at compile time
    And factory functions create type-safe instances
    And generic factory patterns are reusable

  @REQ-GEN-020 @error
  Scenario: Factory functions validate configuration before creation
    Given invalid factory configuration
    When factory is called
    Then appropriate error is returned
    And error indicates invalid configuration
    And instance is not created with invalid config
