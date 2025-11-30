@domain:security @m2 @REQ-SEC-045 @spec(security.md#62-novuspack-security-advantages)
Feature: NovusPack Security Implementation

  @REQ-SEC-045 @happy
  Scenario: NovusPack provides quantum-safe signatures
    Given an open NovusPack package
    And package with quantum-safe signatures
    When security advantages are examined
    Then NovusPack is first package format with ML-DSA support
    And NovusPack is first package format with SLH-DSA support
    And quantum-safe signatures provide future-proof security

  @REQ-SEC-045 @happy
  Scenario: NovusPack provides unified format with multiple signature types
    Given an open NovusPack package
    And package with multiple signature types
    When security advantages are examined
    Then single format supports multiple signature types
    And format supports traditional signatures (PGP, X.509)
    And format supports quantum-safe signatures (ML-DSA, SLH-DSA)
    And format enables flexible signature selection

  @REQ-SEC-045 @happy
  Scenario: NovusPack provides cross-platform compatibility
    Given an open NovusPack package
    And package with cross-platform features
    When security advantages are examined
    Then package works on Windows platforms
    And package works on macOS platforms
    And package works on Linux platforms
    And package provides consistent security across platforms

  @REQ-SEC-045 @happy
  Scenario: NovusPack provides flexible per-file encryption
    Given an open NovusPack package
    And package with per-file encryption
    When security advantages are examined
    Then per-file encryption selection enables selective encryption
    And quantum-safe encryption options are available
    And traditional encryption options are available
    And encryption selection is file-specific

  @REQ-SEC-045 @happy
  Scenario: NovusPack provides transparent antivirus-friendly design
    Given an open NovusPack package
    And package with transparent structure
    When security advantages are examined
    Then package structure is inspectable
    And package headers are designed for easy scanning
    And package format is transparent to antivirus software
    And package uses standard file system operations

  @REQ-SEC-045 @happy
  Scenario: NovusPack provides industry compliance with unique advantages
    Given an open NovusPack package
    And package with industry standard features
    When security advantages are examined
    Then package is compatible with existing security infrastructure
    And package provides quantum-safe advantages over traditional formats
    And package maintains interoperability with industry standards
