@domain:signatures @security @m2 @v2 @REQ-SIG-003 @spec(api_signatures.md#13-immutability-check) @spec(api_signatures.md#immutability-check) @spec(api_writing.md#4-signed-file-write-operations)
Feature: Enforce immutability after signing

  @security
  Scenario: Post-sign write operations are blocked
    Given a package that has been signed
    When I try to add a new file
    Then the operation should be rejected due to immutability

  @security
  Scenario: Immutability check verifies SignatureOffset
    Given a package
    When write operation is attempted
    Then SignatureOffset is checked
    And if SignatureOffset > 0, immutability is enforced
    And write operation is prevented

  @security
  Scenario: All write operations check immutability
    Given a signed package
    When any write operation is attempted
    Then immutability check occurs first
    And write operation is refused
    And structured immutability error is returned

  @security
  Scenario: Only signature addition is allowed on signed packages
    Given a signed package
    When signature addition is attempted
    Then operation succeeds
    When content modification is attempted
    Then operation fails
    And immutability error is returned

  @security
  Scenario: Immutability prevents accidental signature invalidation
    Given a signed package
    When content modification is attempted
    Then immutability protection prevents modification
    And all signatures remain valid
    And package integrity is maintained
