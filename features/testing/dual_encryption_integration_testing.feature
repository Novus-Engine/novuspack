@domain:testing @m2 @REQ-TEST-007 @spec(testing.md#13-dual-encryption-integration-testing)
Feature: Dual Encryption Integration Testing

  @REQ-TEST-007 @happy
  Scenario: Dual encryption integration testing validates combined encryption
    Given a NovusPack package
    And dual encryption integration testing configuration
    When dual encryption integration testing is performed
    Then mixed encryption packages are tested
    And default behavior is tested (ML-KEM when no type specified)
    And user selection is tested (encryption type per file)
    And package operations work with mixed encryption types
    And key management is tested for each encryption type

  @REQ-TEST-007 @happy
  Scenario: Dual encryption integration testing validates mixed encryption packages
    Given a NovusPack package
    And dual encryption integration testing configuration
    When mixed encryption package testing is performed
    Then packages containing files with different encryption types are tested
    And ML-KEM encrypted files are tested
    And AES-256-GCM encrypted files are tested
    And packages with both encryption types operate correctly

  @REQ-TEST-007 @happy
  Scenario: Dual encryption integration testing validates default behavior
    Given a NovusPack package
    And dual encryption integration testing configuration
    When default behavior testing is performed
    Then ML-KEM is used when no encryption type is specified
    And default encryption behavior is verified
    And default behavior is consistent

  @REQ-TEST-007 @happy
  Scenario: Dual encryption integration testing validates user selection
    Given a NovusPack package
    And dual encryption integration testing configuration
    When user selection testing is performed
    Then users can choose encryption type per file
    And user-selected encryption types are honored
    And per-file encryption selection works correctly

  @REQ-TEST-007 @error
  Scenario: Dual encryption integration testing handles errors correctly
    Given a NovusPack package
    When dual encryption integration testing encounters errors
    Then structured error is returned
    And error details are provided
    And error handling validates error conditions
