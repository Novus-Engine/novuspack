@domain:file_format @m1 @REQ-FILEFMT-016 @spec(package_file_format.md#28-header-initialization)
Feature: Header initialization on package creation

  @happy
  Scenario: Header is initialized with correct default values
    Given a new NovusPack package creation
    When the header is initialized
    Then Magic equals 0x4E56504B
    And FormatVersion equals 1
    And Flags are set based on package configuration
    And PackageDataVersion equals 1
    And MetadataVersion equals 1
    And PackageCRC equals 0 or calculated CRC32
    And CreatedTime equals current timestamp
    And ModifiedTime equals current timestamp
    And LocaleID equals 0 or specified locale
    And Reserved equals 0
    And AppID equals 0 or specified application identifier
    And VendorID equals 0 or specified vendor identifier
    And CreatorID equals 0
    And ArchiveChainID equals 0 or archive chain identifier
    And ArchivePartInfo equals 0x00010000 for single archive
    And CommentSize equals 0 if no comment
    And CommentStart equals 0 if no comment
    And SignatureOffset equals 0 if no signatures

  @happy
  Scenario: CreatedTime is immutable after creation
    Given a NovusPack package
    When the package is created
    Then CreatedTime is set to creation timestamp
    When the package is modified
    Then CreatedTime remains unchanged
    And ModifiedTime is updated

  @happy
  Scenario: ModifiedTime updates on package changes
    Given a NovusPack package
    When the package is created
    Then ModifiedTime equals CreatedTime
    When a file is added
    Then ModifiedTime is updated
    When package metadata is changed
    Then ModifiedTime is updated again

  @happy
  Scenario: IndexStart and IndexSize are set when index is written
    Given a NovusPack package
    When the file index is written
    Then IndexStart equals offset to file index
    And IndexSize equals size of file index in bytes
    And IndexStart and IndexSize are consistent

  @happy
  Scenario: CommentSize and CommentStart are set when comment is added
    Given a NovusPack package
    When a package comment is added
    Then CommentSize equals size of comment including null terminator
    And CommentStart equals offset to comment
    And CommentSize matches actual comment length

  @happy
  Scenario: CommentSize and CommentStart are zero when no comment
    Given a NovusPack package without a comment
    When the header is initialized
    Then CommentSize equals 0
    And CommentStart equals 0

  @happy
  Scenario: SignatureOffset is zero when no signatures
    Given a NovusPack package without signatures
    When the header is initialized
    Then SignatureOffset equals 0

  @happy
  Scenario: SignatureOffset is set when first signature is added
    Given a NovusPack package
    When the first signature is added
    Then SignatureOffset equals offset to signature index
    And SignatureOffset is non-zero

  @happy
  Scenario: ArchivePartInfo defaults to single archive
    Given a new NovusPack package for a single archive
    When the header is initialized
    Then ArchivePartInfo equals 0x00010000
    And part number equals 1
    And total parts equals 1

  @happy
  Scenario: ArchivePartInfo is set for split archives
    Given a NovusPack package that is part of a split archive
    When the header is initialized for part N of M
    Then ArchivePartInfo encodes part N in bits 31-16
    And ArchivePartInfo encodes total M in bits 15-0
    And ArchiveChainID links related archive parts
