@domain:core @m2 @REQ-CORE-041 @spec(api_core.md#71-core-integration-points)
Feature: Core Integration Points

  @REQ-CORE-041 @happy
  Scenario: Core package interface integrates with signature system through immutability enforcement
    Given an open NovusPack package
    And package has signatures
    When write operations are attempted
    Then all write operations check SignatureOffset before proceeding
    And immutability is enforced for signed packages
    And signature system integration is maintained

  @REQ-CORE-041 @happy
  Scenario: Core package interface integrates with signature system through write protection
    Given an open NovusPack package
    And package has signatures
    When write operations are attempted
    Then signed packages are protected from write operations by default
    And write protection prevents modification
    And signature integrity is maintained

  @REQ-CORE-041 @happy
  Scenario: Core package interface integrates with signature system through context integration
    Given signature operations
    When signature operations are performed
    Then all signature operations accept ctx context.Context as first parameter
    And context integration follows standard patterns
    And context supports cancellation and timeout

  @REQ-CORE-041 @happy
  Scenario: Core package interface integrates with signature system through error handling
    Given signature operations
    When signature operations encounter errors
    Then signature operations use structured error system
    And error handling follows core error system patterns
    And errors are consistent with API standards
