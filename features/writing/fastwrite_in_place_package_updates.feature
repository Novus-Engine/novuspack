@domain:writing @m2 @REQ-WRITE-002 @spec(api_writing.md#2-fastwrite---in-place-package-updates)
Feature: FastWrite in-place package updates

  @happy
  Scenario: FastWrite allowed only when safe criteria met
    Given a package meeting fast write safety criteria
    When I perform a fast write
    Then the write should complete without violating safety guarantees

  @happy
  Scenario: FastWrite performs entry comparison
    Given an existing package with files
    When FastWrite is called
    Then existing entries are compared with new entries
    And changes are detected
    And only changed entries are updated

  @happy
  Scenario: FastWrite updates only changed entries
    Given an existing package with multiple files
    When FastWrite is called with file modifications
    Then only modified entries are updated in place
    And unchanged entries remain unchanged
    And new entries are appended

  @happy
  Scenario: FastWrite appends new entries to end
    Given an existing package
    When FastWrite is called with new files
    Then new entries are appended to end of file
    And existing entries are not moved
    And package structure is updated

  @happy
  Scenario: FastWrite updates file index and offsets
    Given an existing package
    When FastWrite is called
    Then file index is updated
    And entry offsets are recalculated
    And index consistency is maintained

  @happy
  Scenario: FastWrite is efficient for incremental updates
    Given an existing package with many files
    When FastWrite is called with single file change
    Then only affected data is written
    And I/O operations are minimized
    And write performance is optimized

  @happy
  Scenario: FastWrite is efficient for multiple file changes
    Given an existing package
    When FastWrite is called with multiple file changes
    Then changes are batched efficiently
    And write operations are optimized
    And performance is maintained

  @error
  Scenario: FastWrite fails on signed packages
    Given a signed package with SignatureOffset > 0
    When FastWrite is called
    Then structured validation error is returned
    And error indicates signed package cannot use FastWrite
    And SafeWrite must be used instead

  @error
  Scenario: FastWrite fails on compressed packages
    Given a compressed package
    When FastWrite is called
    Then structured validation error is returned
    And error indicates compressed package cannot use FastWrite
    And SafeWrite must be used instead

  @error
  Scenario: FastWrite validates existing package
    Given an invalid or corrupted existing package
    When FastWrite is called
    Then structured validation error is returned
    And error indicates package validation failure

  @error
  Scenario: FastWrite falls back to SafeWrite on failure
    Given an existing package
    And FastWrite operation fails
    When FastWrite encounters an error
    Then operation falls back to SafeWrite
    And write operation completes using SafeWrite
    And fallback is transparent

  @error
  Scenario: FastWrite handles partial update failures
    Given an existing package
    When FastWrite encounters error during update
    Then partial recovery is attempted
    And successfully updated entries are preserved
    And error is reported

  @error
  Scenario: FastWrite respects context cancellation
    Given an existing package
    And a cancelled context
    When FastWrite is called
    Then structured context error is returned
    And operation is cancelled

  @REQ-WRITE-008 @REQ-WRITE-009 @error
  Scenario: FastWrite validates path parameter
    Given an existing package
    When FastWrite is called with empty path
    Then structured validation error is returned
    And error indicates invalid path

  @REQ-WRITE-008 @REQ-WRITE-011 @error
  Scenario: FastWrite stops operation on context cancellation
    Given an existing package
    And a cancelled context
    When FastWrite is called
    Then structured context error is returned
    And error type is context cancellation
    And write operation stops immediately
