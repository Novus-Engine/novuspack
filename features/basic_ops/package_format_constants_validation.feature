@domain:basic_ops @m1 @REQ-API_BASIC-005 @spec(api_basic_operations.md#21-package-format-constants)
Feature: Package format constants validation

  @happy
  Scenario: NPKMagic constant equals expected value
    Given package format constants
    When NPKMagic is examined
    Then NPKMagic equals 0x4E56504B
    And NPKMagic represents "NVPK" in hex

  @happy
  Scenario: NPKVersion constant equals expected value
    Given package format constants
    When NPKVersion is examined
    Then NPKVersion equals 1
    And NPKVersion is the current format version

  @happy
  Scenario: HeaderSize constant equals expected value
    Given package format constants
    When HeaderSize is examined
    Then HeaderSize equals 112 bytes
    And HeaderSize matches the authoritative header definition

  @happy
  Scenario: Constants validate package header correctly
    Given a NovusPack package file
    When header magic number is read
    Then magic number equals NPKMagic
    When header format version is read
    Then format version equals NPKVersion
    When header size is measured
    Then header size equals HeaderSize

  @error
  Scenario: Invalid magic number is detected using NPKMagic
    Given a file with magic number not equal to NPKMagic
    When header validation is performed
    Then a structured invalid format error is returned
    And error indicates magic number mismatch
