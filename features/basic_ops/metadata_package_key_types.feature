@domain:basic_ops @m2 @REQ-API_BASIC-170 @spec(api_basic_operations.md#2171-metadata-package-key-types)
Feature: metadata package key types

  @REQ-API_BASIC-170 @happy
  Scenario: metadata package defines key types for package and path metadata
    Given the metadata package API
    When key types are used
    Then PackageInfo is defined as a key type
    And PathMetadataEntry is defined as a key type
    And FileEntry is defined as a key type
    And related supporting structures are defined consistently
    And key types align with the metadata system requirements and tech specs

