@domain:metadata @m2 @REQ-META-045 @spec(metadata.md#asset-metadata-example)
Feature: Asset Metadata Example

  @REQ-META-045 @happy
  Scenario: Asset metadata example tracks texture file count
    Given an open NovusPack package
    And texture files in package
    And asset metadata
    When asset metadata example is examined
    Then textures field contains number of texture files
    And texture count matches actual files
    And example demonstrates texture tracking

  @REQ-META-045 @happy
  Scenario: Asset metadata example tracks sound file count
    Given an open NovusPack package
    And sound files in package
    And asset metadata
    When asset metadata example is examined
    Then sounds field contains number of sound files
    And sound count matches actual files
    And example demonstrates sound tracking

  @REQ-META-045 @happy
  Scenario: Asset metadata example tracks model file count
    Given an open NovusPack package
    And model files in package
    And asset metadata
    When asset metadata example is examined
    Then models field contains number of model files
    And model count matches actual files
    And example demonstrates model tracking

  @REQ-META-045 @happy
  Scenario: Asset metadata example tracks script file count
    Given an open NovusPack package
    And script files in package
    And asset metadata
    When asset metadata example is examined
    Then scripts field contains number of script files
    And script count matches actual files
    And example demonstrates script tracking

  @REQ-META-045 @happy
  Scenario: Asset metadata example tracks total asset size
    Given an open NovusPack package
    And asset files in package
    And asset metadata
    When asset metadata example is examined
    Then total_size field contains total asset size in bytes
    And total size matches sum of asset files
    And example demonstrates size tracking

  @REQ-META-045 @happy
  Scenario: Asset metadata example demonstrates complete asset tracking
    Given an open NovusPack package
    And asset files of all types
    And asset metadata with example values
    When asset metadata is retrieved
    Then all asset counts are populated
    And total size is calculated
    And example demonstrates comprehensive asset tracking

  @REQ-META-011 @error
  Scenario: Asset metadata example validation fails with invalid counts
    Given an open NovusPack package
    And asset metadata with invalid count values
    When asset metadata is validated
    Then structured validation error is returned
    And error indicates invalid asset counts
