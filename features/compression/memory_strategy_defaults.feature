@domain:compression @m2 @REQ-COMPR-071 @spec(api_package_compression.md#1321-memory-strategy-defaults)
Feature: Memory Strategy Defaults

  @REQ-COMPR-071 @happy
  Scenario: MemoryStrategyBalanced is the default memory strategy
    Given a compression operation
    And no memory strategy is explicitly specified
    When compression is performed
    Then MemoryStrategyBalanced is used by default
    And 50% of available RAM is used
    And default provides optimal performance

  @REQ-COMPR-071 @happy
  Scenario: Memory strategy defaults provide sensible defaults for most use cases
    Given a compression operation
    And default memory strategy is used
    When compression is performed
    Then defaults provide sensible configuration
    And defaults work for most use cases
    And defaults balance performance and resource usage

  @REQ-COMPR-071 @happy
  Scenario: Memory strategy defaults can be overridden with explicit configuration
    Given a compression operation
    And explicit memory strategy is specified
    When compression is performed
    Then explicit configuration overrides defaults
    And specified strategy is used
    And defaults are bypassed

  @REQ-COMPR-071 @happy
  Scenario: Memory strategy defaults adapt to system capabilities
    Given a compression operation
    And system capabilities are detected
    When default memory strategy is selected
    Then defaults adapt to available RAM
    And defaults adapt to system resources
    And defaults provide appropriate configuration
