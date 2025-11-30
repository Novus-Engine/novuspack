@domain:basic_ops @m1 @REQ-API_BASIC-003 @spec(api_basic_operations.md#61-close-method)
Feature: Close and Release Resources

  @REQ-API_BASIC-003 @happy
  Scenario: Close flushes buffers and finalizes package
    Given an open NovusPack package
    And the package has unsaved changes in buffers
    When Close is called
    Then buffers are flushed
    And package file handle is closed
    And memory buffers and caches are released
    And package state and metadata are cleared
    And cleanup operations are performed

  @REQ-API_BASIC-003 @happy
  Scenario: Close does not modify package file
    Given an open NovusPack package
    And the package has been opened from disk
    When Close is called
    Then package file is not modified
    And original package file remains unchanged
    And resources are released

  @REQ-API_BASIC-003 @error
  Scenario: Close returns error when package is not open
    Given a NovusPack package that is not open
    When Close is called
    Then validation error is returned
    And error indicates package is not currently open

  @REQ-API_BASIC-003 @error
  Scenario: Close handles file system errors gracefully
    Given an open NovusPack package
    And a file system error occurs during closing
    When Close is called
    Then I/O error is returned
    And error indicates file system issue
    And partial cleanup is attempted
