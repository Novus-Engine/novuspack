@domain:file_format @m2 @REQ-FILEFMT-039 @spec(package_file_format.md#281-initial-package-creation)
Feature: Initial Package Creation

  @REQ-FILEFMT-039 @happy
  Scenario: Initial package creation defines package initialization
    Given a new NovusPack package is being created
    When initial package creation is performed
    Then package initialization is defined
    And all header fields are set to default values
    And package is ready for use

  @REQ-FILEFMT-039 @happy
  Scenario: Magic field is set to 0x4E56504B on creation
    Given a new NovusPack package is being created
    When initial package creation is performed
    Then Magic field is set to 0x4E56504B
    And Magic identifies package as NovusPack format
    And Magic value is correct

  @REQ-FILEFMT-039 @happy
  Scenario: FormatVersion is set to 1 on creation
    Given a new NovusPack package is being created
    When initial package creation is performed
    Then FormatVersion is set to 1
    And FormatVersion indicates current format version
    And FormatVersion enables format compatibility checking

  @REQ-FILEFMT-039 @happy
  Scenario: Flags are set based on package configuration
    Given a new NovusPack package is being created
    And package has specific configuration (encryption, signing, compression)
    When initial package creation is performed
    Then Flags are set based on package configuration
    And Flags encode package features
    And Flags encode compression type

  @REQ-FILEFMT-039 @happy
  Scenario: PackageDataVersion and MetadataVersion are set to 1 on creation
    Given a new NovusPack package is being created
    When initial package creation is performed
    Then PackageDataVersion is set to 1
    And MetadataVersion is set to 1
    And initial versions indicate new package

  @REQ-FILEFMT-039 @happy
  Scenario: PackageCRC is set to 0 if skipped on creation
    Given a new NovusPack package is being created
    When initial package creation is performed with CRC calculation skipped
    Then PackageCRC is set to 0 if skipped
    And PackageCRC enables integrity validation when calculated

  @REQ-FILEFMT-039 @happy
  Scenario: PackageCRC is set to calculated CRC32 on creation
    Given a new NovusPack package is being created
    When initial package creation is performed with CRC calculation enabled
    Then PackageCRC is set to calculated CRC32
    And PackageCRC enables integrity validation when calculated

  @REQ-FILEFMT-039 @happy
  Scenario: CreatedTime and ModifiedTime are set to current timestamp on creation
    Given a new NovusPack package is being created
    When initial package creation is performed
    Then CreatedTime is set to current timestamp
    And ModifiedTime is set to current timestamp
    And CreatedTime is immutable after creation
    And ModifiedTime updates on package changes

  @REQ-FILEFMT-039 @happy
  Scenario: LocaleID, AppID, and VendorID have default values on creation
    Given a new NovusPack package is being created
    When initial package creation is performed
    Then LocaleID is set to 0 (system default) or specified locale
    And AppID is set to 0 (no association) or specified application identifier
    And VendorID is set to 0 (no association) or specified vendor identifier

  @REQ-FILEFMT-039 @happy
  Scenario: Reserved and CreatorID are set to 0 on creation
    Given a new NovusPack package is being created
    When initial package creation is performed
    Then Reserved is set to 0
    And CreatorID is set to 0
    And reserved fields are initialized for future use

  @REQ-FILEFMT-039 @happy
  Scenario: ArchivePartInfo is set to 0x00010001 for single archive on creation
    Given a new NovusPack package is being created
    And package is single archive
    When initial package creation is performed
    Then ArchivePartInfo is set to 0x00010001
    And ArchivePartInfo indicates part 1 of 1
    And ArchivePartInfo encodes single archive format
