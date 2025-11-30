@domain:file_format @m2 @REQ-FILEFMT-037 @spec(package_file_format.md#253-content-related-flags)
Feature: Content-Related Flags

  @REQ-FILEFMT-037 @happy
  Scenario: Content-related flags define content flag bits for package features
    Given a NovusPack package
    And package header is examined
    When content-related flags are inspected
    Then Bit 2 indicates Has encrypted files
    And Bit 1 indicates Has compressed files
    And Bit 0 indicates Has signatures

  @REQ-FILEFMT-037 @happy
  Scenario: Bit 2 indicates package contains encrypted files
    Given a NovusPack package
    And package contains encrypted files
    When package header flags are examined
    Then Bit 2 is set to 1
    And flag indicates per-file encryption is present
    And flag corresponds to HasEncryptedData in PackageInfo

  @REQ-FILEFMT-037 @happy
  Scenario: Bit 1 indicates package contains compressed files
    Given a NovusPack package
    And package contains compressed files
    When package header flags are examined
    Then Bit 1 is set to 1
    And flag indicates compressed files are present
    And flag corresponds to HasCompressedData in PackageInfo

  @REQ-FILEFMT-037 @happy
  Scenario: Bit 0 indicates package has digital signatures
    Given a NovusPack package
    And package has digital signatures
    When package header flags are examined
    Then Bit 0 is set to 1
    And flag indicates signatures are present
    And flag corresponds to HasSignatures in PackageInfo
    And flag must be set before adding first signature
