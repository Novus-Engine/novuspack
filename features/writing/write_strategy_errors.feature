@domain:writing @m2 @REQ-WRITE-051 @spec(api_writing.md#572-write-strategy-errors)
Feature: Write Strategy Errors

  @REQ-WRITE-051 @error
  Scenario: Write strategy errors define write-specific errors
    Given a NovusPack package
    When write strategy errors are examined
    Then FastWriteOnCompressed is returned when attempting FastWrite on compressed package
    And CompressionMismatch is returned when compression type doesn't match expectations
    And MemoryInsufficient is returned when insufficient memory for compression operations
    And errors indicate write strategy issues

  @REQ-WRITE-051 @error
  Scenario: FastWriteOnCompressed indicates FastWrite on compressed package attempt
    Given a NovusPack package
    And a compressed package
    When FastWrite is attempted on compressed package
    Then FastWriteOnCompressed error is returned
    And error indicates FastWrite cannot be used with compressed packages
    And error follows structured error format

  @REQ-WRITE-051 @error
  Scenario: CompressionMismatch indicates compression type mismatch
    Given a NovusPack package
    And compression type that doesn't match expectations
    When write operation encounters compression mismatch
    Then CompressionMismatch error is returned
    And error indicates compression type doesn't match
    And error follows structured error format

  @REQ-WRITE-051 @error
  Scenario: MemoryInsufficient indicates insufficient memory
    Given a NovusPack package
    And insufficient memory for compression operations
    When write operation requires compression
    Then MemoryInsufficient error is returned
    And error indicates insufficient memory for compression
    And error follows structured error format

  @REQ-WRITE-051 @error
  Scenario: Write strategy errors provide diagnostic information
    Given a NovusPack package
    And write strategy error conditions
    When write strategy errors are returned
    Then errors provide diagnostic information
    And information helps identify write strategy issues
    And diagnostics enable error resolution
