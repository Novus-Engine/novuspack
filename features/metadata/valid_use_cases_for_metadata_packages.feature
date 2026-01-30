@domain:metadata @m2 @REQ-META-071 @spec(api_metadata.md#62-valid-use-cases)
Feature: Valid Use Cases for Metadata Packages

  @REQ-META-071 @happy
  Scenario: Valid use cases define metadata-only package use cases
    Given a NovusPack package
    And a metadata-only package
    When valid use cases are examined
    Then package catalogs and registries use metadata-only packages
    And configuration and schema packages use metadata-only packages
    And package management operations use metadata-only packages
    And development and build tools use metadata-only packages

  @REQ-META-071 @happy
  Scenario: Package catalogs and registries use metadata-only packages
    Given a NovusPack package
    And a metadata-only package
    When catalog or registry package is created
    Then package listings with metadata are stored
    And dependency resolution trees are stored
    And searchable indexes are stored

  @REQ-META-071 @happy
  Scenario: Configuration and schema packages use metadata-only packages
    Given a NovusPack package
    And a metadata-only package
    When configuration package is created
    Then configuration templates are stored
    And API specifications are stored
    And data structure definitions are stored

  @REQ-META-071 @happy
  Scenario: Package management operations use metadata-only packages
    Given a NovusPack package
    And a metadata-only package
    When package management package is created
    Then update manifests are stored
    And installation scripts are stored
    And package relationships are stored

  @REQ-META-071 @happy
  Scenario: Development and build tools use metadata-only packages
    Given a NovusPack package
    And a metadata-only package
    When development package is created
    Then build configurations are stored
    And development metadata is stored
    And testing configurations are stored

  @REQ-META-071 @happy
  Scenario: Empty and placeholder packages use metadata-only packages
    Given a NovusPack package
    And a metadata-only package with no special metadata files
    When empty or placeholder package is created
    Then package has FileCount 0
    And package has no special metadata files
    And package serves as placeholder or namespace reservation
    And package is valid for future expansion
