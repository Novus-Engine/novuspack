@domain:compression @m2 @REQ-COMPR-077 @spec(api_package_compression.md#13221-memory-detection-process)
Feature: Memory Detection Process

  @REQ-COMPR-077 @happy
  Scenario: Memory detection process queries available system RAM
    Given a compression operation
    When memory detection process runs
    Then system queries available system RAM
    And RAM information is retrieved
    And memory limits are calculated based on available RAM

  @REQ-COMPR-077 @happy
  Scenario: Memory detection process selects Conservative strategy for systems with less than 4GB RAM
    Given a system with less than 4GB RAM
    When memory detection process runs
    Then Conservative strategy is automatically selected
    And 25% of total RAM is allocated for compression
    And system stability is prioritized

  @REQ-COMPR-077 @happy
  Scenario: Memory detection process selects Balanced strategy for systems with 4-16GB RAM
    Given a system with 4-16GB RAM
    When memory detection process runs
    Then Balanced strategy is automatically selected
    And 50% of available RAM is used for compression
    And optimal performance is achieved
    And system responsiveness is maintained

  @REQ-COMPR-077 @happy
  Scenario: Memory detection process selects Aggressive strategy for systems with more than 16GB RAM
    Given a system with more than 16GB RAM
    When memory detection process runs
    Then Aggressive strategy is automatically selected
    And 75% of available RAM is used for compression
    And maximum compression performance is achieved

  @REQ-COMPR-077 @happy
  Scenario: Memory detection process calculates appropriate memory limits
    Given a compression operation
    And available system RAM is known
    When memory detection process runs
    Then appropriate memory limits are calculated
    And limits are based on selected strategy
    And limits are based on available RAM
    And limits prevent out of memory errors

  @REQ-COMPR-077 @happy
  Scenario: Memory detection process adapts to system capabilities
    Given a compression operation
    And system capabilities change
    When memory detection process runs
    Then detection adapts to current system state
    And memory limits are recalculated
    And optimal settings are selected
