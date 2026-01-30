@domain:core @m2 @REQ-CORE-112 @spec(api_core.md#1185-packagemetadata-contents)
Feature: PackageMetadata contents reference PackageInfo structure definition

  @REQ-CORE-112 @happy
  Scenario: PackageMetadata includes package information based on PackageInfo
    Given an opened package
    When GetMetadata is called
    Then PackageMetadata includes package information based on PackageInfo
    And PackageMetadata contains the expected top-level package fields
    And package-level fields are consistent with header-derived information
