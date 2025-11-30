@domain:basic_ops @m2 @REQ-API_BASIC-023 @REQ-API_BASIC-025 @REQ-API_BASIC-051 @spec(api_basic_operations.md#2-package-constants)
Feature: Package Constants

  @REQ-API_BASIC-023 @happy
  Scenario: NPKMagic constant defines package identifier
    Given the NovusPack package format
    When package constants are examined
    Then NPKMagic equals 0x4E56504B
    And magic number represents "NVPK" in hex
    And magic number identifies .npk files

  @REQ-API_BASIC-023 @happy
  Scenario: NPKVersion constant defines current format version
    Given the NovusPack package format
    When package constants are examined
    Then NPKVersion equals 1
    And version represents current format version
    And version is used for compatibility checking

  @REQ-API_BASIC-023 @happy
  Scenario: HeaderSize constant defines fixed header size
    Given the NovusPack package format
    When package constants are examined
    Then HeaderSize equals 112 bytes
    And header size is fixed and constant
    And header size matches package file format specification

  @REQ-API_BASIC-023 @happy
  Scenario: Package constants are used for validation
    Given a package file to be validated
    When package header is validated
    Then NPKMagic is used to identify package format
    And NPKVersion is used to check compatibility
    And HeaderSize is used for header parsing

  @REQ-API_BASIC-025 @happy
  Scenario: Package lifecycle operations follow create-open-operations-close pattern
    Given the NovusPack system
    When lifecycle operations are examined
    Then first phase is Create to create new package
    And second phase is Open to load existing package
    And third phase is Operations for file and metadata operations
    And fourth phase is Close to release resources

  @REQ-API_BASIC-051 @happy
  Scenario: Package operations include validation and information methods
    Given an open NovusPack package
    When package operations are examined
    Then Validate method validates package format and integrity
    And GetInfo method retrieves comprehensive package information
    And Defragment method optimizes package structure
    And ReadHeader method reads header without opening package
