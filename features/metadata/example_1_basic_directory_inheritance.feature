@domain:metadata @m2 @REQ-META-032 @spec(metadata.md#1341-example-1-basic-directory-inheritance)
Feature: Example 1 Basic Directory Inheritance

  @REQ-META-032 @happy
  Scenario: Directory metadata file defines directory with properties
    Given an open NovusPack package
    And directory metadata file with "/assets/" directory
    When directory metadata is examined
    Then directory path is "/assets/"
    And directory has category property set to "texture"
    And directory has compression property set to "lossless"

  @REQ-META-032 @happy
  Scenario: Directory inheritance is enabled with priority
    Given an open NovusPack package
    And directory metadata file with "/assets/" directory
    When directory inheritance settings are examined
    Then inheritance is enabled
    And priority is set to 1
    And directory provides tag inheritance

  @REQ-META-032 @happy
  Scenario: File inherits tags from directory
    Given an open NovusPack package
    And directory metadata file with "/assets/" directory
    And directory has category="texture" and compression="lossless"
    And file at "/assets/texture.png"
    When file tags are examined
    Then file inherits category tag from directory
    And file inherits compression tag from directory
    And inherited tags are applied to file

  @REQ-META-032 @happy
  Scenario: Multiple directory levels provide inheritance hierarchy
    Given an open NovusPack package
    And directory metadata with "/assets/" directory
    And directory metadata with "/assets/textures/" directory
    And file at "/assets/textures/ui/button.png"
    When file tags are examined
    Then file inherits from "/assets/" directory
    And file inherits from "/assets/textures/" directory
    And inheritance follows directory path

  @REQ-META-032 @happy
  Scenario: Child directory tags override parent directory tags
    Given an open NovusPack package
    And "/assets/" directory with category="texture"
    And "/assets/textures/" directory with category="image"
    And file at "/assets/textures/file.png"
    When file tags are examined
    Then file inherits category="image" from child directory
    And child directory tag overrides parent directory tag
    And override priority rules are followed

  @REQ-META-032 @happy
  Scenario: Basic inheritance example demonstrates simple tag inheritance
    Given an open NovusPack package
    And directory metadata file with basic inheritance setup
    And files in directory hierarchy
    When inheritance is applied
    Then files inherit tags from parent directories
    And inheritance follows basic rules
    And inheritance behavior matches example

  @REQ-META-011 @error
  Scenario: Inheritance fails if directory metadata is invalid
    Given an open NovusPack package
    And invalid directory metadata file
    When inheritance is applied
    Then structured validation error is returned
    And error indicates invalid directory metadata

  @REQ-META-011 @error
  Scenario: Inheritance fails if path does not match
    Given an open NovusPack package
    And directory metadata with path "/assets/"
    And file at path not matching directory
    When inheritance is applied
    Then structured validation error is returned
    And error indicates path mismatch
