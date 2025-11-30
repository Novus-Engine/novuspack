@domain:compression @m2 @REQ-COMPR-078 @spec(api_package_compression.md#13222-chunk-size-calculation)
Feature: Chunk Size Calculation

  @REQ-COMPR-078 @happy
  Scenario: Chunk size calculation determines optimal chunk sizes automatically
    Given a compression operation
    And chunk size is set to 0 for auto-calculation
    When chunk size calculation runs
    Then optimal chunk size is calculated
    And chunk size is 25% of allocated memory limit
    And chunk size fits within memory constraints

  @REQ-COMPR-078 @happy
  Scenario: Chunk size calculation ensures chunks fit in memory
    Given a compression operation
    And allocated memory limit is known
    When chunk size is calculated
    Then chunk size is smaller than memory limit
    And multiple chunks can fit in memory
    And concurrent operations are possible

  @REQ-COMPR-078 @happy
  Scenario: Chunk size calculation allows for multiple concurrent chunks
    Given a compression operation
    And parallel processing is enabled
    When chunk size is calculated
    Then chunk size allows multiple concurrent chunks
    And worker threads can process chunks simultaneously
    And memory is efficiently utilized

  @REQ-COMPR-078 @happy
  Scenario: Chunk size can be explicitly specified
    Given a compression operation
    And explicit chunk size is specified
    When compression runs
    Then specified chunk size is used
    And auto-calculation is bypassed
    And explicit value takes precedence

  @REQ-COMPR-078 @happy
  Scenario: Chunk size calculation adapts to memory strategy
    Given a compression operation
    And memory strategy is selected
    When chunk size is calculated
    Then chunk size is based on memory strategy allocation
    And Conservative strategy results in smaller chunks
    And Aggressive strategy results in larger chunks
