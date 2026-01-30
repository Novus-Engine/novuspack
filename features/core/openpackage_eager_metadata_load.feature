@domain:core @m2 @REQ-CORE-048 @spec(api_core.md#111-opened-package-reader-contract)
Feature: OpenPackage eager metadata load loads all required metadata

  @REQ-CORE-048 @happy
  Scenario: OpenPackage eagerly loads all required metadata into memory
    Given a valid package on disk
    When OpenPackage completes successfully
    Then all required package metadata is loaded into memory
    And subsequent PackageReader calls do not require additional metadata loading
    And the package is ready for pure in-memory queries
