@domain:metadata @m2 @REQ-META-119 @spec(metadata.md#133-tag-inheritance-rules) @spec(api_metadata.md#81-path-metadata-structures)
Feature: Path Display in Metadata Operations

  Background:
    Given an open NovusPack package
    And a valid context

  @REQ-META-117 @happy
  Scenario: Path metadata stored with leading slash
    Given PathMetadataEntry for a path
    When path is stored in metadata
    Then path has leading slash
    And path follows internal storage format

  @REQ-META-119 @happy
  Scenario: Path metadata displayed without leading slash
    Given PathMetadataEntry with stored path "/assets/textures/"
    When path is displayed to user
    Then displayed path is "assets/textures/"
    And leading slash is stripped for display

  @REQ-META-119 @happy
  Scenario: Path inheritance display strips leading slash
    Given path hierarchy with inheritance
    And paths stored as "/assets/", "/assets/textures/", "/assets/textures/ui/"
    When inheritance chain is displayed
    Then paths shown as "assets/", "assets/textures/", "assets/textures/ui/"
    And leading slashes are not shown to user

  @REQ-META-119 @happy
  Scenario: Path metadata file display format
    Given path metadata file with stored paths
    When metadata is exported for user viewing
    Then paths in exported format do not show leading slash
    And format is user-friendly

  @REQ-META-117 @REQ-META-119 @happy
  Scenario: Path metadata API returns display format
    Given path metadata query operations
    When paths are returned via API
    Then returned paths have leading slash stripped
    And paths are in user-facing format

  @REQ-META-119 @error
  Scenario: Path metadata error messages use display format
    Given path metadata operation that fails
    When error message includes path
    Then error shows path without leading slash
    And path is user-friendly in error message
