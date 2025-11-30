@domain:metadata @m2 @REQ-META-074 @spec(api_metadata.md#623-package-management-operations)
Feature: Package Management Operations

  @REQ-META-074 @happy
  Scenario: Package management operations use metadata-only packages
    Given a NovusPack package
    And a metadata-only package
    When package management operations use packages
    Then update manifests are stored in metadata files
    And installation scripts are stored in metadata files
    And package relationships are stored in metadata files
    And package contains only special metadata files

  @REQ-META-074 @happy
  Scenario: Update manifests describe updates without actual files
    Given a NovusPack package
    And a metadata-only package
    When update manifest package is created
    Then manifest describes updates without actual files
    And update information is stored in metadata
    And manifest enables update management

  @REQ-META-074 @happy
  Scenario: Installation scripts contain installation instructions
    Given a NovusPack package
    And a metadata-only package
    When installation script package is created
    Then package contains installation instructions
    And instructions are stored in metadata files
    And scripts enable automated installation

  @REQ-META-074 @happy
  Scenario: Package relationships define inter-package relationships
    Given a NovusPack package
    And a metadata-only package
    When package relationship package is created
    Then package defines inter-package relationships
    And relationships are stored in metadata
    And relationships enable package dependency management

  @REQ-META-074 @error
  Scenario: Package management operations must be valid metadata-only packages
    Given a NovusPack package
    When package management package is validated
    Then package must have FileCount of 0
    And package must have special metadata files
    And package must be valid metadata-only package
