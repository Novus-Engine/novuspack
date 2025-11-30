@domain:compression @m2 @REQ-COMPR-054 @spec(api_package_compression.md#11363-execution)
Feature: Compression Execution Flow

  @REQ-COMPR-054 @happy
  Scenario: Execution calls CompressPackageStream with custom configuration
    Given custom memory management configuration
    When CompressPackageStream is called with ZSTD compression type
    Then compression executes with custom configuration
    And custom memory settings are applied
    And compression completes successfully

  @REQ-COMPR-054 @happy
  Scenario: Execution applies custom memory execution settings
    Given custom memory execution requirements
    When execution occurs
    Then chunk size settings are applied
    And temp directory settings are used
    And memory usage limits are enforced
    And custom configuration optimizes execution

  @REQ-COMPR-054 @error
  Scenario: Execution handles compression errors appropriately
    Given compression execution operations
    When compression errors occur
    Then structured compression errors are returned
    And error handling uses structured error system
    And errors provide context for diagnosis
