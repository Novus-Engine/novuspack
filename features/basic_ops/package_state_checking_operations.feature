@domain:basic_ops @m1 @REQ-API_BASIC-014 @spec(api_basic_operations.md#75-check-package-state)
Feature: Package state checking operations

  @happy
  Scenario: IsOpen returns true for open packages
    Given an open NovusPack package
    When IsOpen is called
    Then true is returned

  @happy
  Scenario: IsOpen returns false for closed packages
    Given a closed NovusPack package
    When IsOpen is called
    Then false is returned

  @happy
  Scenario: IsReadOnly returns true for read-only packages
    Given a read-only open NovusPack package
    When IsReadOnly is called
    Then true is returned

  @happy
  Scenario: IsReadOnly returns false for writable packages
    Given a writable open NovusPack package
    When IsReadOnly is called
    Then false is returned

  @happy
  Scenario: GetPath returns package file path
    Given an open NovusPack package at a specific path
    When GetPath is called
    Then the package file path is returned
    And path matches the opened file path

  @happy
  Scenario: GetPath returns empty string for new packages
    Given a new Package instance that has not been opened
    When GetPath is called
    Then empty string is returned

  @happy
  Scenario: Package state transitions correctly
    Given a new Package instance
    When package is created
    Then IsOpen is false
    When package is opened
    Then IsOpen is true
    When package is closed
    Then IsOpen is false again
