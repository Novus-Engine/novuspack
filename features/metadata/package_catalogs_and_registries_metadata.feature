@domain:metadata @m2 @REQ-META-072 @spec(api_metadata.md#621-package-catalogs-and-registries)
Feature: Package Catalogs and Registries

  @REQ-META-072 @happy
  Scenario: Package catalogs and registries use metadata-only packages
    Given a NovusPack package
    And a metadata-only package
    When package catalogs and registries are created
    Then package listings with metadata are stored
    And dependency resolution trees are stored
    And searchable indexes are stored
    And package contains only special metadata files

  @REQ-META-072 @happy
  Scenario: Package listings store catalogs of available packages
    Given a NovusPack package
    And a metadata-only package
    When package catalog is created
    Then catalog contains listings of available packages with metadata
    And catalog serves as package discovery mechanism
    And catalog enables package browsing

  @REQ-META-072 @happy
  Scenario: Dependency resolution stores dependency trees
    Given a NovusPack package
    And a metadata-only package
    When dependency resolution package is created
    Then package defines dependency trees
    And dependencies can be resolved from metadata
    And dependency relationships are stored

  @REQ-META-072 @happy
  Scenario: Searchable indexes store package repository indexes
    Given a NovusPack package
    And a metadata-only package
    When searchable index package is created
    Then index enables package discovery
    And index supports searching package repositories
    And index provides navigation data

  @REQ-META-072 @error
  Scenario: Package catalogs and registries must be valid metadata-only packages
    Given a NovusPack package
    When catalog or registry package is validated
    Then package must have FileCount of 0
    And package must have special metadata files
    And package must be valid metadata-only package
