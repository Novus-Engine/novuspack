@domain:basic_ops @m2 @REQ-API_BASIC-095 @spec(api_metadata.md#713-packageinfofromheader-helper-method)
Feature: Package loading synchronizes PackageInfo from header

  @REQ-API_BASIC-095 @happy
  Scenario: Package loading uses PackageInfo.FromHeader to synchronize in-memory PackageInfo
    Given a package file opened from disk
    And a header has been read from the on-disk package
    When package loading populates package metadata
    Then PackageInfo is synchronized from the header
    And the synchronization uses PackageInfo.FromHeader
    And PackageInfo values reflect the on-disk header values
    And PackageInfo is ready for API reads after open completes
    And PackageInfo synchronization occurs as part of package loading

