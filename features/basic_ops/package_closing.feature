@domain:basic_ops @m2 @REQ-API_BASIC-046 @spec(api_basic_operations.md#6-package-closing)
Feature: Package Closing

  @REQ-API_BASIC-046 @happy
  Scenario: Close closes package file handle
    Given an open NovusPack package
    When Close is called
    Then package file handle is closed
    And file handle resources are released
    And file is no longer accessible through package

  @REQ-API_BASIC-046 @happy
  Scenario: Close releases memory buffers and caches
    Given an open NovusPack package
    And package has loaded memory buffers
    When Close is called
    Then memory buffers are released
    And cached data is cleared
    And memory resources are freed

  @REQ-API_BASIC-046 @happy
  Scenario: Close clears package state and metadata
    Given an open NovusPack package
    And package has loaded state and metadata
    When Close is called
    Then package state is cleared
    And package metadata is cleared
    And package IsOpen state is set to false

  @REQ-API_BASIC-046 @happy
  Scenario: Close performs necessary cleanup operations
    Given an open NovusPack package
    And package has active operations
    When Close is called
    Then cleanup operations are performed
    And all resources are properly released
    And package is in clean state

  @REQ-API_BASIC-046 @happy
  Scenario: Close does not modify package file
    Given an open NovusPack package
    And package has unsaved changes
    When Close is called
    Then package file is not modified
    And unsaved changes are not persisted
    And Write methods must be used to save changes

  @REQ-API_BASIC-046 @error
  Scenario: Close returns I/O error on file system errors
    Given an open NovusPack package
    And file system error occurs during closing
    When Close is called
    Then I/O error is returned
    And error indicates file system issue
    And error follows structured error format

  @REQ-API_BASIC-046 @error
  Scenario: Close returns validation error when package not open
    Given a NovusPack package that is not open
    When Close is called
    Then validation error is returned
    And error indicates package is not currently open
    And error follows structured error format

  @REQ-API_BASIC-003 @happy
  Scenario: CloseWithCleanup closes package and performs cleanup
    Given an open NovusPack package
    When CloseWithCleanup is called with context
    Then package file is closed
    And cleanup operations are performed
    And all resources including metadata caches are released
    And package transitions to cleaned state

  @REQ-API_BASIC-003 @happy
  Scenario: CloseWithCleanup clears in-memory metadata caches
    Given an open NovusPack package
    And package has loaded metadata caches
    When CloseWithCleanup is called with context
    Then all in-memory metadata caches are cleared
    And file entry caches are cleared
    And path metadata caches are cleared
    And pure in-memory read operations fail after cleanup

  @REQ-API_BASIC-003 @happy
  Scenario: CloseWithCleanup performs defragmentation and optimization
    Given an open NovusPack package
    When CloseWithCleanup is called with context
    Then defragmentation may be performed
    And package optimization may be performed
    And cleanup operations complete successfully
    And operation may take longer than standard Close

  @REQ-API_BASIC-003 @happy
  Scenario: CloseWithCleanup vs Close behavior differences
    Given an open NovusPack package
    When CloseWithCleanup is called instead of Close
    Then CloseWithCleanup clears metadata caches
    But Close may preserve metadata caches
    And CloseWithCleanup prevents cached read operations
    But Close may allow cached GetInfo and ListFiles

  @REQ-API_BASIC-003 @error
  Scenario: CloseWithCleanup respects context cancellation
    Given an open NovusPack package
    And a cancelled context
    When CloseWithCleanup is called
    Then context cancellation is respected
    And cleanup operations are cancelled
    And structured context error is returned
    And partial cleanup may have occurred

  @REQ-API_BASIC-003 @error
  Scenario: CloseWithCleanup respects context timeout
    Given an open NovusPack package
    And a context with short timeout
    And cleanup operations are time-consuming
    When CloseWithCleanup is called
    Then timeout is respected
    And structured context error is returned
    And error indicates timeout occurred

  @REQ-API_BASIC-003 @happy
  Scenario: CloseWithCleanup is idempotent
    Given a package that has been closed with CloseWithCleanup
    When CloseWithCleanup is called again
    Then operation succeeds without error
    And no additional cleanup is performed
    And idempotent behavior is maintained
