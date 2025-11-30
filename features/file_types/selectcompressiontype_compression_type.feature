@domain:file_types @m2 @REQ-FILETYPES-011 @spec(file_type_system.md#221-selectcompressiontype)
Feature: SelectCompressionType

  @REQ-FILETYPES-011 @happy
  Scenario: SelectCompressionType determines compression type from file type
    Given a NovusPack package
    And file data
    And a file type
    When SelectCompressionType is called with data and file type
    Then compression type is determined based on file type
    And compression type is returned as uint8

  @REQ-FILETYPES-011 @happy
  Scenario: SelectCompressionType skips compression for already compressed formats
    Given a NovusPack package
    And file data for JPEG format
    And FileTypeJPEG file type
    When SelectCompressionType is called
    Then CompressionNone is returned
    When file data for PNG format with FileTypePNG is processed
    Then CompressionNone is returned
    When file data for GIF format with FileTypeGIF is processed
    Then CompressionNone is returned
    When file data for MP3 format with FileTypeMP3 is processed
    Then CompressionNone is returned
    When file data for MP4 format with FileTypeMP4 is processed
    Then CompressionNone is returned
    When file data for OGG format with FileTypeOGG is processed
    Then CompressionNone is returned
    When file data for FLAC format with FileTypeFLAC is processed
    Then CompressionNone is returned

  @REQ-FILETYPES-011 @happy
  Scenario: SelectCompressionType handles special files correctly
    Given a NovusPack package
    And file data
    When FileTypeSignature is processed
    Then CompressionNone is returned for signature files
    When FileTypeMetadata is processed
    Then CompressionZstd is returned for YAML special files
    When FileTypeManifest is processed
    Then CompressionZstd is returned for YAML special files
    When FileTypeIndex is processed
    Then CompressionZstd is returned for YAML special files
    When other special files are processed
    Then CompressionZstd is returned as default for special files

  @REQ-FILETYPES-011 @happy
  Scenario: SelectCompressionType selects CompressionZstd for text-based files
    Given a NovusPack package
    And file data
    When FileTypeText is processed
    Then CompressionZstd is returned
    When FileTypeScript is processed
    Then CompressionZstd is returned
    When FileTypeConfig is processed
    Then CompressionZstd is returned

  @REQ-FILETYPES-011 @happy
  Scenario: SelectCompressionType selects CompressionLZ4 for binary media files
    Given a NovusPack package
    And file data
    When FileTypeImage is processed
    Then CompressionLZ4 is returned
    When FileTypeAudio is processed
    Then CompressionLZ4 is returned
    When FileTypeVideo is processed
    Then CompressionLZ4 is returned

  @REQ-FILETYPES-011 @happy
  Scenario: SelectCompressionType defaults to CompressionZstd
    Given a NovusPack package
    And file data
    And file type that does not match specific rules
    When SelectCompressionType is called
    Then CompressionZstd is returned as default compression method
