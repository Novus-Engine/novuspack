@domain:file_format @m2 @REQ-FILEFMT-035 @spec(package_file_format.md#251-flags-field-encoding)
Feature: Flags Field Encoding

  @REQ-FILEFMT-035 @happy
  Scenario: Flags field encoding defines flag representation
    Given a NovusPack package header
    When flags field encoding is examined
    Then flag representation is defined
    And flags are 32-bit unsigned integer
    And flag bits encode package features and compression type

  @REQ-FILEFMT-035 @happy
  Scenario: Bits 31-24 are reserved and must be 0
    Given a NovusPack package header
    When flags field encoding is examined
    Then bits 31-24 are reserved for future use
    And reserved bits must be 0
    And reserved bits enable future extensibility

  @REQ-FILEFMT-035 @happy
  Scenario: Bits 23-16 are reserved and must be 0
    Given a NovusPack package header
    When flags field encoding is examined
    Then bits 23-16 are reserved for future use
    And reserved bits must be 0
    And reserved bits enable future extensibility

  @REQ-FILEFMT-035 @happy
  Scenario: Bits 15-8 encode package compression type
    Given a NovusPack package header
    When flags field encoding is examined
    Then bits 15-8 encode package compression type
    And compression type 0 indicates no compression
    And compression type 1 indicates Zstd
    And compression type 2 indicates LZ4
    And compression type 3 indicates LZMA

  @REQ-FILEFMT-035 @happy
  Scenario: Bits 7-0 encode package features
    Given a NovusPack package header
    When flags field encoding is examined
    Then bits 7-0 encode package features
    And feature bits indicate package-level features
    And feature bits include metadata and content flags

  @REQ-FILEFMT-035 @error
  Scenario: Flags with non-zero reserved bits are invalid
    Given a NovusPack package header
    And flags have non-zero reserved bits (31-16)
    When flags are validated
    Then validation fails
    And structured validation error is returned
    And reserved bits violation is detected
