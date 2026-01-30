@skip @domain:basic_ops @m2 @spec(api_basic_operations.md#12-package-loading-process)
Feature: Basic Operations Definitions

# This feature captures high-level expectations for package loading and header inspection.
# Detailed runnable scenarios for basic operations live in dedicated basic_ops feature files.

  @REQ-API_BASIC-022 @architecture
  Scenario: Package loading reads the header and initializes in-memory state
    Given a valid package file exists on disk
    When the package is opened
    Then the header is parsed and validated
    And core metadata is loaded into the in-memory package state

  @REQ-API_BASIC-022 @REQ-API_BASIC-014 @behavior
  Scenario: Header inspection reads only the header without opening a package instance
    Given a package file path
    When the header is read for inspection
    Then header fields are returned without opening the full package
    And callers can use the result to decide whether to open the package
