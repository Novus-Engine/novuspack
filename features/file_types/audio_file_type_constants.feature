@domain:file_types @m2 @REQ-FILETYPES-021 @spec(file_type_system.md#31126-audio-file-types-7000-7999)
Feature: Audio File Type Constants

  @REQ-FILETYPES-021 @happy
  Scenario: Audio file type constants define audio file type range
    Given a NovusPack package
    When audio file type constants are examined
    Then FileTypeAudioStart is 7000
    And FileTypeAudioEnd is 7999
    And FileTypeAudio is 7000
    And audio file types are within range 7000-7999

  @REQ-FILETYPES-021 @happy
  Scenario: Specific audio file type constants are defined
    Given a NovusPack package
    When audio file type constants are examined
    Then FileTypeMP3 is 7001
    And FileTypeWAV is 7002
    And FileTypeOGG is 7003
    And FileTypeFLAC is 7004
    And FileTypeAAC is 7005
    And FileTypeWMA is 7006
    And FileTypeAIFF is 7007
    And FileTypeALAC is 7008
    And FileTypeAPE is 7009
    And FileTypeOpus is 7010
    And FileTypeM4A is 7011

  @REQ-FILETYPES-021 @happy
  Scenario: Audio file types are recognized by IsAudioFile
    Given a NovusPack package
    When FileTypeMP3 is checked with IsAudioFile
    Then IsAudioFile returns true
    When FileTypeWAV is checked with IsAudioFile
    Then IsAudioFile returns true
    When FileTypeFLAC is checked with IsAudioFile
    Then IsAudioFile returns true
    When FileTypeOGG is checked with IsAudioFile
    Then IsAudioFile returns true

  @REQ-FILETYPES-021 @error
  Scenario: Non-audio file types are not recognized by IsAudioFile
    Given a NovusPack package
    When FileTypeText is checked with IsAudioFile
    Then IsAudioFile returns false
    When FileTypeImage is checked with IsAudioFile
    Then IsAudioFile returns false
    When FileTypeBinary is checked with IsAudioFile
    Then IsAudioFile returns false

  @REQ-FILETYPES-021 @happy
  Scenario: Audio file types support format validation
    Given a NovusPack package
    And an audio file with type FileTypeMP3
    When the audio file is validated
    Then the file format validation is performed
    And format validation is appropriate for audio files
