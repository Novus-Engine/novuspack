@domain:metadata @m2 @REQ-META-073 @spec(api_metadata.md#622-configuration-and-schema-packages)
Feature: Configuration and Schema Packages

  @REQ-META-073 @happy
  Scenario: Configuration and schema packages use metadata-only packages
    Given a NovusPack package
    And a metadata-only package
    When configuration and schema packages are created
    Then configuration templates are stored in metadata files
    And API specifications are stored in metadata files
    And data structure definitions are stored in metadata files
    And package contains only special metadata files

  @REQ-META-073 @happy
  Scenario: Configuration templates use metadata-only packages
    Given a NovusPack package
    And a metadata-only package
    When configuration template package is created
    Then package contains configuration schemas
    And package has FileCount of 0
    And package contains special metadata files
    And package serves as configuration template

  @REQ-META-073 @happy
  Scenario: API specifications use metadata-only packages
    Given a NovusPack package
    And a metadata-only package
    When API specification package is created
    Then package contains API definitions
    And package contains API schemas
    And package serves as API documentation

  @REQ-META-073 @happy
  Scenario: Data structure definitions use metadata-only packages
    Given a NovusPack package
    And a metadata-only package
    When data structure definition package is created
    Then package contains data models
    And package contains structure definitions
    And package serves as schema reference

  @REQ-META-073 @error
  Scenario: Configuration packages must be valid metadata-only packages
    Given a NovusPack package
    When configuration package is validated
    Then package must have FileCount of 0
    And package must have special metadata files
    And package must be valid metadata-only package
