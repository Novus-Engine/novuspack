@domain:core @m2 @REQ-CORE-187 @spec(api_core.md#1112-openpackage-eager-metadata-load)
Feature: OpenPackage eager metadata load requires all package metadata to be loaded into memory

  @REQ-CORE-187 @happy
  Scenario: OpenPackage loads all metadata into memory
    Given a valid package on disk
    When OpenPackage completes successfully
    Then all package metadata is loaded into memory
    And subsequent reader operations do not require additional metadata I/O
    And the behavior matches the OpenPackage eager metadata load specification
