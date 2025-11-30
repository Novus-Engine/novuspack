@domain:writing @m2 @REQ-WRITE-036 @spec(api_writing.md#5221-fastwrite-behavior-for-compressed-packages)
Feature: FastWrite Behavior for Compressed Packages

  @REQ-WRITE-036 @error
  Scenario: FastWrite cannot be used with compressed packages
    Given an open NovusPack package
    And the package is compressed (compression type in header flags)
    When FastWrite is called with the target path
    Then FastWriteOnCompressed error is returned
    And error indicates FastWrite is not supported for compressed packages
    And error follows structured error format

  @REQ-WRITE-036 @happy
  Scenario: FastWrite automatically falls back to SafeWrite for compressed packages
    Given an open NovusPack package
    And the package is compressed
    When Write is called with the target path
    Then FastWrite is not attempted
    And SafeWrite is automatically selected
    And compressed package is handled with SafeWrite
    And operation completes successfully

  @REQ-WRITE-036 @error
  Scenario: FastWrite returns error immediately when compressed package is detected
    Given an open NovusPack package
    And the package is compressed
    When FastWrite is called directly with the target path
    Then error is returned immediately
    And no in-place update is attempted
    And error prevents FastWrite execution
    And error indicates compressed package limitation

  @REQ-WRITE-036 @happy
  Scenario: Automatic selection prevents FastWrite attempt on compressed packages
    Given an open NovusPack package
    And the package is compressed
    When Write method selects write strategy
    Then compressed package detection identifies compression
    And FastWrite is not considered
    And SafeWrite is directly selected
    And no FastWrite attempt occurs

  @REQ-WRITE-036 @happy
  Scenario: FastWrite fallback to SafeWrite handles compressed packages correctly
    Given an open NovusPack package
    And the package is compressed
    When write operation is performed
    Then automatic fallback to SafeWrite occurs
    And SafeWrite handles compressed package correctly
    And decompression/recompression is performed
    And operation completes successfully
