@domain:metadata @m2 @REQ-META-044 @spec(metadata.md#asset-metadata)
Feature: Asset Metadata Tags

  @REQ-META-044 @happy
  Scenario: Asset metadata provides asset tagging support
    Given a NovusPack package
    When asset metadata is used
    Then textures count is tracked
    And sounds count is tracked
    And models count is tracked
    And scripts count is tracked
    And total asset size is tracked

  @REQ-META-044 @happy
  Scenario: Asset metadata tracks texture file counts
    Given a NovusPack package
    And texture files in the package
    When asset metadata is examined
    Then textures field contains number of texture files
    And texture count is accurate

  @REQ-META-044 @happy
  Scenario: Asset metadata tracks audio file counts
    Given a NovusPack package
    And sound files in the package
    When asset metadata is examined
    Then sounds field contains number of sound files
    And sound count is accurate

  @REQ-META-044 @happy
  Scenario: Asset metadata tracks model and script counts
    Given a NovusPack package
    And model files in the package
    And script files in the package
    When asset metadata is examined
    Then models field contains number of model files
    And scripts field contains number of script files
    And counts are accurate

  @REQ-META-044 @happy
  Scenario: Asset metadata tracks total asset size
    Given a NovusPack package
    And asset files in the package
    When asset metadata is examined
    Then total_size field contains total asset size in bytes
    And total size is calculated correctly
