@domain:basic_ops @m1 @REQ-API_BASIC-005 @spec(api_basic_operations.md#21-package-format-constants)
Feature: Package format constants validation

  @happy
  Scenario: NVPKMagic constant equals expected value
    Given package format constants
    When NVPKMagic is examined
    Then NVPKMagic equals 0x4E56504B
    And NVPKMagic represents "NVPK" in hex

  @happy
  Scenario: NVPKVersion constant equals expected value
    Given package format constants
    When NVPKVersion is examined
    Then NVPKVersion equals 1
    And NVPKVersion is the current format version

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
    Then magic number equals NVPKMagic
    When header format version is read
    Then format version equals NVPKVersion
    When header size is measured
    Then header size equals HeaderSize

  @error
  Scenario: Invalid magic number is detected using NVPKMagic
    Given a file with magic number not equal to NVPKMagic
    When header validation is performed
    Then a structured invalid format error is returned
    And error indicates magic number mismatch
