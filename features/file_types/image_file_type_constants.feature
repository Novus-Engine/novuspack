@domain:file_types @m2 @REQ-FILETYPES-020 @spec(file_type_system.md#31125-image-file-types-5000-6999)
Feature: Image File Type Constants

  @REQ-FILETYPES-020 @happy
  Scenario: Image file type constants define image file type range
    Given a NovusPack package
    When image file type constants are examined
    Then FileTypeImageStart is 5000
    And FileTypeImageEnd is 6999
    And FileTypeImage is 5000
    And image file types are within range 5000-6999

  @REQ-FILETYPES-020 @happy
  Scenario: Specific image file type constants are defined
    Given a NovusPack package
    When image file type constants are examined
    Then FileTypePNG is 5001
    And FileTypeJPEG is 5002
    And FileTypeGIF is 5003
    And FileTypeBMP is 5004
    And FileTypeWebP is 5005
    And FileTypeTIFF is 5006
    And FileTypeSVG is 5007
    And FileTypeRAW is 5008
    And FileTypeHEIC is 5009
    And FileTypeAVIF is 5010

  @REQ-FILETYPES-020 @happy
  Scenario: Image file types are recognized by IsImageFile
    Given a NovusPack package
    When FileTypeImage is checked with IsImageFile
    Then IsImageFile returns true
    When FileTypePNG is checked with IsImageFile
    Then IsImageFile returns true
    When FileTypeJPEG is checked with IsImageFile
    Then IsImageFile returns true
    When FileTypeGIF is checked with IsImageFile
    Then IsImageFile returns true

  @REQ-FILETYPES-020 @error
  Scenario: Non-image file types are not recognized by IsImageFile
    Given a NovusPack package
    When FileTypeText is checked with IsImageFile
    Then IsImageFile returns false
    When FileTypeAudio is checked with IsImageFile
    Then IsImageFile returns false
    When FileTypeBinary is checked with IsImageFile
    Then IsImageFile returns false

  @REQ-FILETYPES-020 @happy
  Scenario: Image file types support format validation
    Given a NovusPack package
    And an image file with type FileTypePNG
    When the image file is validated
    Then format validation is performed
    And image processing is appropriate for image files
