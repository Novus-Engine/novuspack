@domain:compression @m2 @REQ-COMPR-075 @spec(api_package_compression.md#13214-custom-strategy)
Feature: Custom strategy allows custom memory configuration

  @REQ-COMPR-075 @happy
  Scenario: Custom strategy uses explicit MaxMemoryUsage value
    Given compression operations with MemoryStrategyCustom
    When custom strategy is applied
    Then explicit MaxMemoryUsage value is used
    And automatic detection is overridden
    And custom memory limit is enforced

  @REQ-COMPR-075 @happy
  Scenario: Custom strategy allows specific memory constraints
    Given compression operations requiring specific memory limits
    When MemoryStrategyCustom is used
    Then specific memory constraints can be set
    And MaxMemoryUsage value is explicitly specified
    And memory allocation matches requirements

  @REQ-COMPR-075 @happy
  Scenario: Custom strategy overrides automatic detection
    Given compression operations
    When MemoryStrategyCustom is selected
    Then automatic memory detection is bypassed
    And explicit memory configuration is used
    And custom memory management is applied

  @REQ-COMPR-075 @happy
  Scenario: Custom strategy is useful for specific memory constraints
    Given compression operations with specific memory requirements
    When custom strategy is used
    Then custom strategy is useful for specific constraints
    And memory configuration matches exact needs
    And fine-tuned memory control is available
