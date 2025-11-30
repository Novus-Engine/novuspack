@domain:file_format @m1 @REQ-FILEFMT-017 @spec(package_file_format.md#11-file-layout-order)
Feature: File layout order validation

  @happy
  Scenario: File layout follows correct order
    Given a valid NovusPack package file
    When file structure is examined
    Then package header comes first (fixed-size 112 bytes)
    Then file entries and data follow (variable length)
    Then file index follows file entries and data
    Then package comment follows file index (optional, variable length)
    Then digital signatures follow comment (optional, variable length)

  @happy
  Scenario: Package header is at offset 0
    Given a NovusPack package file
    When file structure is examined
    Then package header starts at offset 0
    And header is exactly HeaderSize bytes
    And header is immediately followed by file entries

  @happy
  Scenario: File entries and data are interleaved
    Given a NovusPack package with multiple files
    When file structure is examined
    Then file entry 1 is immediately followed by file data 1
    And file entry 2 is immediately followed by file data 2
    And pattern continues for all files
    And entries and data alternate correctly

  @happy
  Scenario: File index follows all file entries and data
    Given a NovusPack package
    When file structure is examined
    Then file index starts after all file entries and data
    And IndexStart points to file index location
    And IndexSize matches file index size

  @happy
  Scenario: Package comment follows file index
    Given a NovusPack package with comment
    When file structure is examined
    Then package comment starts after file index
    And CommentStart points to comment location
    And CommentSize matches comment size

  @happy
  Scenario: Digital signatures follow package comment
    Given a signed NovusPack package
    When file structure is examined
    Then digital signatures start after package comment
    And SignatureOffset points to signature location
    And signatures are at end of file

  @happy
  Scenario: File layout order is enforced during creation
    Given a new NovusPack package
    When package is created
    Then sections are written in correct order
    And offsets are calculated correctly
    And layout matches specification

  @error
  Scenario: Invalid layout order is detected
    Given a corrupted NovusPack package with invalid layout
    When package structure is validated
    Then layout validation fails
    And structured validation error is returned
    And error indicates layout violation

  @error
  Scenario: Overlapping sections are invalid
    Given a NovusPack package with overlapping sections
    When package structure is validated
    Then validation fails
    And structured corruption error is returned
    And error indicates section overlap
