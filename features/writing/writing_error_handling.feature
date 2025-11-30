@domain:writing @m2 @REQ-WRITE-027 @REQ-WRITE-050 @REQ-WRITE-051 @spec(api_writing.md#44-error-conditions)
Feature: Writing Error Handling

  @REQ-WRITE-027 @error
  Scenario: Error conditions define signed file write errors
    Given a NovusPack package
    And a signed package
    When write operation errors are examined
    Then error conditions define signed file write errors
    And errors indicate signed package write restrictions
    And errors follow structured error format

  @REQ-WRITE-027 @error
  Scenario: Signed file write errors prevent invalid operations
    Given a NovusPack package
    And a signed package
    When write operation is attempted
    Then error indicates signed package write restrictions
    And error prevents signature invalidation
    And error guides proper workflow

  @REQ-WRITE-050 @error
  Scenario: Compression errors define compression-specific errors
    Given a NovusPack package
    When compression errors are examined
    Then CompressionFailure indicates compression operation failure
    And DecompressionFailure indicates decompression operation failure
    And UnsupportedCompression indicates unsupported compression type
    And CorruptedCompressedData indicates corrupted compressed data
    And CompressSignedPackageError indicates signed package compression attempt

  @REQ-WRITE-050 @error
  Scenario: Compression errors provide diagnostic information
    Given a NovusPack package
    And compression error conditions
    When compression errors are returned
    Then errors provide specific diagnostic information
    And information helps identify compression issues
    And diagnostics enable error resolution

  @REQ-WRITE-051 @error
  Scenario: Write strategy errors define write-specific errors
    Given a NovusPack package
    When write strategy errors are examined
    Then FastWriteOnCompressed indicates FastWrite on compressed package attempt
    And CompressionMismatch indicates compression type mismatch
    And MemoryInsufficient indicates insufficient memory
    And errors indicate write strategy issues

  @REQ-WRITE-051 @error
  Scenario: Write strategy errors provide diagnostic information
    Given a NovusPack package
    And write strategy error conditions
    When write strategy errors are returned
    Then errors provide specific diagnostic information
    And information helps identify write strategy issues
    And diagnostics enable error resolution

  @REQ-WRITE-027 @REQ-WRITE-050 @REQ-WRITE-051 @error
  Scenario: Error handling follows structured error format
    Given a NovusPack package
    And error conditions
    When errors are returned
    Then all errors follow structured error format
    And errors provide consistent diagnostic information
    And error format enables systematic error handling
