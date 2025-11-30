@domain:writing @m2 @REQ-WRITE-007 @spec(api_writing.md#1-safewrite---atomic-package-writing)
Feature: General write method with compression handling

  @happy
  Scenario: Write method writes package with compression handling
    Given an open writable package
    When Write is called
    Then package is written to disk
    And compression is handled according to package state
    And write operation completes successfully

  @happy
  Scenario: Write handles compressed packages correctly
    Given an open writable compressed package
    When Write is called
    Then package is written with compression preserved
    And compressed data is written correctly

  @happy
  Scenario: Write handles uncompressed packages correctly
    Given an open writable uncompressed package
    When Write is called
    Then package is written without compression
    And data is written correctly

  @error
  Scenario: Write fails if package is read-only
    Given a read-only package
    When Write is called
    Then structured validation error is returned

  @error
  Scenario: Write fails if package is signed
    Given a signed package
    When Write is called
    Then structured validation error is returned
    And error indicates immutability violation

  @REQ-WRITE-008 @REQ-WRITE-009 @error
  Scenario: Write validates path parameter
    Given an open writable package
    When Write is called with empty path
    Then structured validation error is returned
    And error indicates invalid path

  @REQ-WRITE-008 @REQ-WRITE-010 @error
  Scenario: Write validates compression type parameter
    Given an open writable package
    When Write is called with invalid compression type
    Then structured validation error is returned
    And error indicates unsupported compression type

  @REQ-WRITE-008 @REQ-WRITE-011 @error
  Scenario: Write respects context cancellation
    Given an open writable package
    And a cancelled context
    When Write is called
    Then structured context error is returned
    And error type is context cancellation
    And write operation stops
