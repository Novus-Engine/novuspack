@domain:metadata @m2 @REQ-META-037 @spec(metadata.md#151-texture-file-tagging)
Feature: Texture File Tagging

  @REQ-META-037 @happy
  Scenario: Texture file tagging demonstrates texture metadata
    Given a NovusPack package
    And a texture file
    When texture file tagging is used
    Then comprehensive tags are set on texture files
    And category tag identifies texture type
    And type tag identifies texture classification
    And format tag specifies image format
    And size tag specifies texture dimensions
    And compression tag specifies compression method
    And priority tag specifies priority level
    And descriptive tags provide additional information

  @REQ-META-037 @happy
  Scenario: Texture file tagging example demonstrates usage
    Given a NovusPack package
    And a UI button texture file
    When texture file is tagged
    Then category tag is set to texture category
    And format tag is set to "PNG"
    And size tag specifies dimensions like 1024x1024
    And compression tag is set to "lossless"
    And priority tag is set to priority level like 5
    And UI/button/interface tags provide descriptive information

  @REQ-META-037 @happy
  Scenario: Texture tags support texture management
    Given a NovusPack package
    And multiple texture files with tags
    When texture files are searched by tags
    Then files can be found by category
    And files can be found by format
    And files can be found by priority level
    And tags enable texture organization

  @REQ-META-037 @error
  Scenario: Texture tags validate tag values
    Given a NovusPack package
    When invalid texture tag values are provided
    Then tag validation detects invalid values
    And appropriate errors are returned
