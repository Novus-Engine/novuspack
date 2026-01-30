@domain:metadata @m2 @REQ-META-032 @spec(metadata.md#1341-example-1-basic-path-inheritance)
Feature: Example 1 Basic Path Inheritance

  @REQ-META-032 @happy
  Scenario: Path metadata file defines path with properties
    Given an open NovusPack package
    And path metadata file with "/assets/" path
    When path metadata is examined
    Then path is "/assets/"
    And path has category property set to "texture"
    And path has compression property set to "lossless"

  @REQ-META-032 @happy
  Scenario: Path inheritance is enabled with priority
    Given an open NovusPack package
    And path metadata file with "/assets/" path
    When path inheritance settings are examined
    Then inheritance is enabled
    And priority is set to 1
    And path provides tag inheritance

  @REQ-META-032 @happy
  Scenario: File inherits tags from path via PathMetadataEntry
    Given an open NovusPack package
    And path metadata file with "/assets/" path
    And path has category="texture" and compression="lossless"
    And file at "/assets/texture.png" with associated PathMetadataEntry
    When GetEffectiveTags is called on PathMetadataEntry
    Then effective tags include category tag from path
    And effective tags include compression tag from path
    And inherited tags are applied via PathMetadataEntry

  @REQ-META-032 @happy
  Scenario: Multiple path levels provide inheritance hierarchy
    Given an open NovusPack package
    And path metadata with "/assets/" path
    And path metadata with "/assets/textures/" path
    And file at "/assets/textures/ui/button.png" with associated PathMetadataEntry
    When GetEffectiveTags is called on PathMetadataEntry
    Then effective tags include tags from "/assets/" path via ParentPath
    And effective tags include tags from "/assets/textures/" path
    And inheritance follows path hierarchy via ParentPath

  @REQ-META-032 @happy
  Scenario: Child path tags override parent path tags
    Given an open NovusPack package
    And "/assets/" path with category="texture"
    And "/assets/textures/" path with category="image"
    And file at "/assets/textures/file.png" with associated PathMetadataEntry
    When GetEffectiveTags is called on PathMetadataEntry
    Then effective tags include category="image" from child path
    And child path tag overrides parent path tag
    And override priority rules are followed

  @REQ-META-032 @happy
  Scenario: Basic inheritance example demonstrates simple tag inheritance
    Given an open NovusPack package
    And path metadata file with basic inheritance setup
    And files with associated PathMetadataEntry instances
    When GetEffectiveTags is called on PathMetadataEntry instances
    Then effective tags include tags from parent paths via ParentPath
    And inheritance follows basic rules
    And inheritance behavior matches example

  @REQ-META-011 @error
  Scenario: Inheritance fails if path metadata is invalid
    Given an open NovusPack package
    And invalid path metadata file
    When inheritance is applied
    Then structured validation error is returned
    And error indicates invalid path metadata

  @REQ-META-011 @error
  Scenario: Inheritance fails if path does not match
    Given an open NovusPack package
    And path metadata with path "/assets/"
    And file at path not matching path metadata
    When inheritance is applied
    Then structured validation error is returned
    And error indicates path mismatch
