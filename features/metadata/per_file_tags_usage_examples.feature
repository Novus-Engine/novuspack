@domain:metadata @m2 @REQ-META-036 @spec(metadata.md#15-per-file-tags-usage-examples)
Feature: Per-File Tags Usage Examples

  @REQ-META-036 @happy
  Scenario: Texture file tagging sets comprehensive tags
    Given an open NovusPack package
    And a texture file entry
    When texture file tags are set
    Then category tag identifies texture
    And type tag specifies texture type
    And format tag specifies image format
    And size tag contains dimensions
    And compression tag specifies compression
    And priority tag sets priority level
    And descriptive tags provide additional metadata

  @REQ-META-036 @happy
  Scenario: Texture file tagging example demonstrates UI button texture
    Given an open NovusPack package
    And a UI button texture file
    When texture file tags are set with example values
    Then format tag is set to PNG
    And size tag indicates 1024x1024
    And compression tag is set to lossless
    And priority tag is set to 5
    And UI/button/interface tags are included
    And example demonstrates comprehensive tagging

  @REQ-META-036 @happy
  Scenario: Audio file tagging sets audio-specific tags
    Given an open NovusPack package
    And an audio file entry
    When audio file tags are set
    Then category tag identifies audio
    And type tag specifies audio type
    And format tag specifies audio format
    And duration tag contains duration
    And loop tag indicates loop settings
    And volume tag contains volume level

  @REQ-META-036 @happy
  Scenario: Audio file tagging example demonstrates ambient sound
    Given an open NovusPack package
    And an ambient forest sound file
    When audio file tags are set with example values
    Then format tag is set to WAV
    And duration tag indicates 120 seconds
    And loop tag is enabled
    And volume tag is set to 0.7
    And example demonstrates audio metadata

  @REQ-META-036 @happy
  Scenario: Path tagging sets tags for inheritance
    Given an open NovusPack package
    And a PathMetadataEntry
    When path tags are set
    Then category tag identifies path type
    And compression tag sets path compression
    And path tags are inherited by child paths via ParentPath
    And inheritance provides default tag values

  @REQ-META-036 @happy
  Scenario: Path tagging example demonstrates textures path
    Given an open NovusPack package
    And a textures PathMetadataEntry
    When path tags are set with example values
    Then category tag is set to texture
    And compression tag is set to lossless
    And mipmaps tag is enabled
    And example demonstrates path inheritance

  @REQ-META-036 @happy
  Scenario: File search by tags finds files by tag values
    Given an open NovusPack package
    And files with various tags
    When GetFilesByTag is called with category tag
    Then all files with matching category are found
    And search results are accurate
    And tag-based search works correctly

  @REQ-META-036 @happy
  Scenario: File search by tags finds files by type
    Given an open NovusPack package
    And files with type tags
    When GetFilesByTag is called with type tag
    Then all UI files are found
    And search results match type criteria
    And type-based search works correctly

  @REQ-META-036 @happy
  Scenario: File search by tags finds files by priority level
    Given an open NovusPack package
    And files with priority tags
    When GetFilesByTag is called with priority tag
    Then all high-priority files are found
    And search results match priority criteria
    And priority-based search works correctly

  @REQ-META-011 @error
  Scenario: File search by tags fails with invalid tag key
    Given an open NovusPack package
    And files with tags
    When GetFilesByTag is called with invalid tag key
    Then structured validation error is returned
    And error indicates invalid tag key
