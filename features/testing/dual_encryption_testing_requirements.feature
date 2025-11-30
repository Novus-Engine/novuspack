@domain:testing @m2 @REQ-TEST-004 @REQ-TEST-007 @spec(testing.md#1-dual-encryption-testing-requirements)
Feature: Dual Encryption Testing Requirements

  @REQ-TEST-004 @happy
  Scenario: Dual encryption testing requirements define encryption testing needs
    Given a NovusPack package
    And dual encryption testing configuration
    When dual encryption testing is performed
    Then ML-KEM encryption testing requirements are defined
    And AES-256-GCM encryption testing requirements are defined
    And dual encryption integration testing requirements are defined
    And encryption testing needs are comprehensive

  @REQ-TEST-007 @happy
  Scenario: Dual encryption integration testing validates combined encryption
    Given a NovusPack package
    And dual encryption integration testing configuration
    When dual encryption integration testing is performed
    Then mixed encryption packages are tested (different encryption types)
    And default behavior is tested (ML-KEM when no type specified)
    And user selection is tested (encryption type per file)
    And package operations work with mixed encryption types
    And appropriate keys are used for each encryption type

  @REQ-TEST-007 @happy
  Scenario: Dual encryption integration testing validates mixed encryption packages
    Given a NovusPack package
    And dual encryption integration testing configuration
    When mixed encryption package testing is performed
    Then packages containing ML-KEM encrypted files are tested
    And packages containing AES-256-GCM encrypted files are tested
    And packages with both encryption types are tested
    And mixed encryption packages operate correctly

  @REQ-TEST-007 @happy
  Scenario: Dual encryption integration testing validates default behavior
    Given a NovusPack package
    And dual encryption integration testing configuration
    When default behavior testing is performed
    Then ML-KEM is used when no encryption type is specified
    And default encryption type behavior is verified
    And default behavior is consistent across operations

  @REQ-TEST-007 @happy
  Scenario: Dual encryption integration testing validates key management
    Given a NovusPack package
    And dual encryption integration testing configuration
    When key management testing is performed
    Then ML-KEM keys are used for ML-KEM encrypted files
    And AES keys are used for AES-256-GCM encrypted files
    And appropriate keys are used for each encryption type
    And key management handles mixed encryption correctly
