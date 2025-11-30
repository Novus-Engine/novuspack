@domain:core @m1 @REQ-CORE-006 @spec(api_core.md#13-package-interface)
Feature: Package interface combines reader and writer

  @happy
  Scenario: Package interface implements PackageReader
    Given a Package instance
    When interface is checked
    Then Package implements PackageReader interface
    And ReadFile method is available
    And ListFiles method is available
    And GetMetadata method is available
    And Validate method is available
    And GetInfo method is available

  @happy
  Scenario: Package interface implements PackageWriter
    Given a Package instance
    When interface is checked
    Then Package implements PackageWriter interface
    And WriteFile method is available
    And RemoveFile method is available
    And Write method is available
    And SafeWrite method is available
    And FastWrite method is available

  @happy
  Scenario: Package interface provides Close method
    Given a Package instance
    When interface is checked
    Then Package provides Close method
    And Close releases resources
    And Close sets IsOpen to false

  @happy
  Scenario: Package interface provides IsOpen method
    Given a Package instance
    When interface is checked
    Then Package provides IsOpen method
    And IsOpen returns package state

  @happy
  Scenario: Package interface provides Defragment method
    Given a Package instance
    When interface is checked
    Then Package provides Defragment method
    And Defragment optimizes package structure

  @happy
  Scenario: Package interface methods work together
    Given a Package instance
    When package operations are performed
    Then reader methods work correctly
    And writer methods work correctly
    And state methods work correctly
    And all interface methods are accessible

  @REQ-CORE-015 @REQ-CORE-016 @error
  Scenario: Package interface methods accept context and respect cancellation
    Given a Package instance
    And a cancelled context
    When package operation with context is called
    Then structured context error is returned
    And error type is context cancellation

  @REQ-CORE-015 @REQ-CORE-018 @error
  Scenario: Package interface methods validate input parameters
    Given a Package instance
    When package operation is called with invalid parameters
    Then structured validation error is returned
    And error indicates invalid input
