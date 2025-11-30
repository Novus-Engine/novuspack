@domain:metadata @m2 @REQ-META-046 @spec(metadata.md#audio-file-tagging)
Feature: Audio Metadata Tags

  @REQ-META-046 @happy
  Scenario: Audio file tagging demonstrates audio metadata
    Given a NovusPack package
    And an audio file
    When audio file tagging is used
    Then category tag is set for audio files
    And type tag identifies audio format
    And format tag specifies audio codec
    And duration tag specifies audio length
    And loop tag indicates looping behavior
    And volume tag specifies audio level

  @REQ-META-046 @happy
  Scenario: Audio file tagging sets audio-specific tags
    Given a NovusPack package
    And an audio file with WAV format
    When audio file is tagged
    Then category tag is set to "audio"
    And format tag is set to "WAV"
    And audio-specific properties are tagged

  @REQ-META-046 @happy
  Scenario: Audio file tagging includes duration and loop settings
    Given a NovusPack package
    And an audio file with 120-second duration
    When audio file is tagged
    Then duration tag is set to 120 seconds
    And loop tag can be enabled or disabled
    And volume tag can specify audio level

  @REQ-META-046 @error
  Scenario: Audio file tagging validates tag values
    Given a NovusPack package
    And an audio file
    When invalid tag values are provided
    Then tag validation detects invalid values
    And appropriate error is returned
