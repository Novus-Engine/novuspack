@domain:metadata @m2 @REQ-META-075 @spec(api_metadata.md#624-development-and-build-tools)
Feature: Development and Build Tools

  @REQ-META-075 @happy
  Scenario: Development and build tools use metadata-only packages
    Given a NovusPack package
    And a metadata-only package
    When development and build tools use packages
    Then build configurations are stored in metadata files
    And development metadata is stored in metadata files
    And testing configurations are stored in metadata files
    And package contains only special metadata files

  @REQ-META-075 @happy
  Scenario: Build configurations use metadata-only packages
    Given a NovusPack package
    And a metadata-only package
    When build configuration package is created
    Then package contains build system configurations
    And package has FileCount of 0
    And package serves as build configuration reference

  @REQ-META-075 @happy
  Scenario: Development metadata uses metadata-only packages
    Given a NovusPack package
    And a metadata-only package
    When development metadata package is created
    Then package contains development environment specifications
    And package contains development tool configurations
    And package serves as development reference

  @REQ-META-075 @happy
  Scenario: Testing configurations use metadata-only packages
    Given a NovusPack package
    And a metadata-only package
    When testing configuration package is created
    Then package contains test specifications
    And package contains test configurations
    And package serves as testing reference

  @REQ-META-075 @error
  Scenario: Development packages must be valid metadata-only packages
    Given a NovusPack package
    When development package is validated
    Then package must have FileCount of 0
    And package must have special metadata files
    And package must be valid metadata-only package
