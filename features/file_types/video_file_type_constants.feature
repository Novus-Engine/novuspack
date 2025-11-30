@domain:file_types @m2 @REQ-FILETYPES-022 @spec(file_type_system.md#31127-video-file-types-8000-9999)
Feature: Video File Type Constants

  @REQ-FILETYPES-022 @happy
  Scenario: Video file type constants define video file type range
    Given a NovusPack package
    When video file type constants are examined
    Then FileTypeVideoStart is 8000
    And FileTypeVideoEnd is 9999
    And FileTypeVideo is 8000
    And video file types are within range 8000-9999

  @REQ-FILETYPES-022 @happy
  Scenario: Specific video file type constants are defined
    Given a NovusPack package
    When video file type constants are examined
    Then FileTypeMP4 is 8001
    And FileTypeMKV is 8002
    And FileTypeAVI is 8003
    And FileTypeMOV is 8004
    And FileTypeWebM is 8005
    And FileTypeWMV is 8006
    And FileTypeFLV is 8007
    And FileTypeM4V is 8008
    And FileTypeMPEG is 8009

  @REQ-FILETYPES-022 @happy
  Scenario: Video file types are recognized by IsVideoFile
    Given a NovusPack package
    When FileTypeVideo is checked with IsVideoFile
    Then IsVideoFile returns true
    When FileTypeMP4 is checked with IsVideoFile
    Then IsVideoFile returns true
    When FileTypeMKV is checked with IsVideoFile
    Then IsVideoFile returns true
    When FileTypeAVI is checked with IsVideoFile
    Then IsVideoFile returns true

  @REQ-FILETYPES-022 @error
  Scenario: Non-video file types are not recognized by IsVideoFile
    Given a NovusPack package
    When FileTypeText is checked with IsVideoFile
    Then IsVideoFile returns false
    When FileTypeImage is checked with IsVideoFile
    Then IsVideoFile returns false
    When FileTypeAudio is checked with IsVideoFile
    Then IsVideoFile returns false

  @REQ-FILETYPES-022 @happy
  Scenario: Video file types support format validation
    Given a NovusPack package
    And a video file with type FileTypeMP4
    When the video file is validated
    Then format validation is performed
    And video processing is appropriate for video files
