@domain:basic_ops @m2 @REQ-API_BASIC-096 @spec(api_basic_operations.md#44-settargetpath-method)
Feature: SetTargetPath Method

  @REQ-API_BASIC-096 @REQ-API_BASIC-097 @happy
  Scenario: SetTargetPath changes package target write path
    Given a package is created or opened
    When SetTargetPath is called with a new path
    Then the package target write path is updated
    And the new path will be used for Write, SafeWrite, or FastWrite operations

  @REQ-API_BASIC-098 @happy
  Scenario: SetTargetPath validates path immediately
    Given a package needs a new target path
    When SetTargetPath is called with a valid path
    Then the provided path is validated for correctness
    And path format is checked
    And no package file is created on disk

  @REQ-API_BASIC-098 @happy
  Scenario: SetTargetPath validates target directory immediately
    Given a package needs a new target path
    When SetTargetPath is called
    Then target directory existence is verified
    And directory writability is checked
    And minimal filesystem I/O is performed for validation

  @REQ-API_BASIC-098 @happy
  Scenario: SetTargetPath enables early error detection
    Given a package with invalid target directory
    When SetTargetPath is called before write operations
    Then path validation error is detected early
    And error is returned before any write attempt
    And package state remains unchanged

  @REQ-API_BASIC-103 @happy
  Scenario: SetTargetPath differs from Create in usage
    Given an existing package opened from disk
    When SetTargetPath is called to change write location
    Then package target path is changed
    And package content is not reinitialized
    And only path configuration is updated

  @REQ-API_BASIC-101 @error
  Scenario: SetTargetPath returns error for invalid path
    Given a package needs a new target path
    When SetTargetPath is called with invalid or malformed path
    Then validation error is returned
    And error indicates invalid path format
    And package target path remains unchanged

  @REQ-API_BASIC-101 @error
  Scenario: SetTargetPath returns error when directory does not exist
    Given a package needs a new target path
    When SetTargetPath is called with non-existent directory
    Then validation error is returned
    And error indicates directory does not exist
    And package target path remains unchanged

  @REQ-API_BASIC-101 @error
  Scenario: SetTargetPath returns error when directory is not writable
    Given a package needs a new target path
    When SetTargetPath is called with non-writable directory
    Then validation error is returned
    And error indicates insufficient permissions
    And package target path remains unchanged
