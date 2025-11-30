@domain:compression @m2 @REQ-COMPR-076 @spec(api_package_compression.md#1322-auto-detection-logic)
Feature: Auto-Detection Logic

  @REQ-COMPR-076 @happy
  Scenario: Auto-detection logic detects system RAM and configures memory
    Given a compression operation
    When auto-detection logic runs
    Then system RAM is queried
    And appropriate memory limits are calculated
    And memory strategy is selected based on available RAM

  @REQ-COMPR-076 @happy
  Scenario: Auto-detection selects Conservative strategy for systems with less than 4GB RAM
    Given a system with less than 4GB RAM
    When auto-detection logic runs
    Then Conservative strategy is automatically selected
    And 25% of total RAM is allocated for compression
    And system stability is maintained

  @REQ-COMPR-076 @happy
  Scenario: Auto-detection selects Balanced strategy for systems with 4-16GB RAM
    Given a system with 4-16GB RAM
    When auto-detection logic runs
    Then Balanced strategy is automatically selected
    And 50% of available RAM is used for compression
    And optimal performance is achieved

  @REQ-COMPR-076 @happy
  Scenario: Auto-detection selects Aggressive strategy for systems with more than 16GB RAM
    Given a system with more than 16GB RAM
    When auto-detection logic runs
    Then Aggressive strategy is automatically selected
    And 75% of available RAM is used for compression
    And maximum performance is achieved

  @REQ-COMPR-076 @happy
  Scenario: Auto-detection calculates optimal chunk size automatically
    Given a compression operation
    And chunk size is set to 0 for auto-calculation
    When auto-detection logic runs
    Then chunk size is calculated as 25% of allocated memory limit
    And chunk size fits within memory constraints
    And chunk size allows for concurrent operations

  @REQ-COMPR-076 @happy
  Scenario: Auto-detection detects CPU cores for worker count
    Given a compression operation
    And MaxWorkers is set to 0 for auto-detection
    When auto-detection logic runs
    Then number of available CPU cores is detected
    And worker count is set accordingly
    And optimal parallel processing is enabled
