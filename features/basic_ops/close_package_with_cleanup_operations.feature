@domain:basic_ops @m1 @REQ-API_BASIC-009 @spec(api_basic_operations.md#62-close-with-cleanup)
Feature: Close package with cleanup operations

  @happy
  Scenario: CloseWithCleanup closes package and performs cleanup
    Given an open NovusPack package
    When CloseWithCleanup is called
    Then package file is closed
    And cleanup operations are performed
    And resources are released
    And IsOpen is set to false

  @happy
  Scenario: CloseWithCleanup performs defragmentation if needed
    Given an open NovusPack package with unused space
    When CloseWithCleanup is called
    Then package structure is optimized
    And unused space is removed
    And data sections are compacted
    And indexes are updated

  @happy
  Scenario: CloseWithCleanup preserves metadata and signatures
    Given an open NovusPack package with metadata and signatures
    When CloseWithCleanup is called
    Then all package metadata is preserved
    And all signatures remain valid
    And package integrity is maintained

  @happy
  Scenario: CloseWithCleanup releases all resources
    Given an open NovusPack package
    When CloseWithCleanup is called
    Then file handle is closed
    And memory buffers are released
    And caches are cleared
    And package state is reset

  @error
  Scenario: CloseWithCleanup fails if package is not open
    Given a closed NovusPack package
    When CloseWithCleanup is called
    Then a structured validation error is returned

  @error
  Scenario: CloseWithCleanup respects context cancellation
    Given an open NovusPack package
    And a cancelled context
    When CloseWithCleanup is called
    Then a structured context error is returned
    And cleanup operations are cancelled gracefully
