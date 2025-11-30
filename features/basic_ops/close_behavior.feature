@domain:basic_ops @m2 @REQ-API_BASIC-047 @spec(api_basic_operations.md#611-close-behavior)
Feature: Close Behavior

  @REQ-API_BASIC-047 @happy
  Scenario: Close closes package file handle
    Given an open NovusPack package
    When Close is called
    Then package file handle is closed
    And file handle resources are released

  @REQ-API_BASIC-047 @happy
  Scenario: Close releases memory buffers and caches
    Given an open NovusPack package
    And package has memory buffers and caches
    When Close is called
    Then memory buffers are released
    And caches are cleared
    And memory is freed

  @REQ-API_BASIC-047 @happy
  Scenario: Close clears package state and metadata
    Given an open NovusPack package
    When Close is called
    Then package state is cleared
    And package metadata is cleared
    And package is no longer open

  @REQ-API_BASIC-047 @happy
  Scenario: Close performs necessary cleanup operations
    Given an open NovusPack package
    When Close is called
    Then cleanup operations are performed
    And resources are released
    And package is properly finalized

  @REQ-API_BASIC-047 @happy
  Scenario: Close does not modify package file
    Given an open NovusPack package opened from disk
    When Close is called
    Then package file is not modified
    And original package file remains unchanged
    And only resources are released
