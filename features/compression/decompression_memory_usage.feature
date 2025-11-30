@domain:compression @m2 @REQ-COMPR-084 @spec(api_package_compression.md#1342-decompression)
Feature: Decompression Memory Usage

  @REQ-COMPR-084 @happy
  Scenario: Decompression requires memory for decompressed content
    Given a decompression operation
    And a compressed package
    When decompression is performed
    Then memory is required for decompressed content
    And memory needs scale with decompressed size
    And memory allocation occurs during decompression

  @REQ-COMPR-084 @happy
  Scenario: Decompression may need temporary storage for large packages
    Given a decompression operation
    And a large compressed package
    When decompression is performed
    Then temporary storage may be needed
    And temporary files are used if memory is insufficient
    And decompression handles large packages successfully

  @REQ-COMPR-084 @happy
  Scenario: Decompression uses streaming for memory-constrained environments
    Given a decompression operation
    And a large compressed package
    And memory constraints exist
    When decompression is performed
    Then streaming decompression is used
    And memory constraints are respected
    And decompression succeeds despite memory limitations

  @REQ-COMPR-084 @happy
  Scenario: Large files use chunked decompression with temp file management
    Given a decompression operation
    And a large compressed package
    And memory constraints exist
    When decompression is performed
    Then chunked decompression is used
    And temporary file management handles memory constraints
    And decompression processes chunks sequentially

  @REQ-COMPR-084 @happy
  Scenario: Memory limits enforce MaxMemoryUsage to prevent system OOM during decompression
    Given a decompression operation
    And MaxMemoryUsage is configured
    And a compressed package
    When decompression is performed
    Then MaxMemoryUsage limit is enforced
    And system out-of-memory errors are prevented
    And decompression memory usage stays within configured limits
