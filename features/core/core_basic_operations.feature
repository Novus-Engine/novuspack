@domain:core @m2 @REQ-CORE-030 @spec(api_core.md#2-basic-operations)
Feature: Core Basic Operations

  @REQ-CORE-030 @happy
  Scenario: Create creates a new package at the specified path
    Given a package path
    When Create is called
    Then new package is created at specified path
    And package is ready for use
    And package file is initialized

  @REQ-CORE-030 @happy
  Scenario: Open opens an existing package from the specified path
    Given an existing package file
    When Open is called with package path
    Then existing package is opened from specified path
    And package content is loaded
    And package is ready for operations

  @REQ-CORE-030 @happy
  Scenario: Close closes the package and releases resources
    Given an open NovusPack package
    When Close is called
    Then package is closed
    And resources are released
    And file handles are closed

  @REQ-CORE-030 @happy
  Scenario: Write handles package writing with compression
    Given a NovusPack package
    When Write is called with path and compression type
    Then package is written using SafeWrite or FastWrite methods
    And compression handling is applied
    And write operation completes

  @REQ-CORE-030 @happy
  Scenario: Defragment optimizes package structure and removes unused space
    Given a NovusPack package with unused space
    When Defragment is called
    Then package structure is optimized
    And unused space is removed
    And package is more efficient

  @REQ-CORE-030 @happy
  Scenario: Validate validates package format, structure, and integrity
    Given a NovusPack package
    When Validate is called
    Then package format is validated
    And package structure is validated
    And package integrity is checked

  @REQ-CORE-030 @happy
  Scenario: GetInfo gets comprehensive package information
    Given a NovusPack package
    When GetInfo is called
    Then comprehensive package information is retrieved
    And package details are available
    And information includes package metadata
