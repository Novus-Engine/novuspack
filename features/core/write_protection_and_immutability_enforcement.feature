@domain:core @m1 @REQ-CORE-009 @spec(api_core.md#74-write-protection-and-immutability-enforcement)
Feature: Write protection and immutability enforcement

  @happy
  Scenario: Signed file detection checks SignatureOffset
    Given a NovusPack package
    When write operation is attempted
    Then SignatureOffset is checked first
    And if SignatureOffset > 0, package is identified as signed
    And write protection is enforced

  @error
  Scenario: Write operations are refused on signed packages
    Given a signed NovusPack package with SignatureOffset > 0
    When AddFileFromMemory is called
    Then write operation is refused
    And a structured immutability error is returned
    And error type is ErrTypeValidation
    And package remains unchanged

  @error
  Scenario: Header modifications are prohibited on signed packages
    Given a signed NovusPack package
    When header modification is attempted
    Then modification is prohibited
    And a structured immutability error is returned
    And header remains unchanged

  @error
  Scenario: Content changes are prohibited on signed packages
    Given a signed NovusPack package
    When file entry modification is attempted
    Then modification is prohibited
    And a structured immutability error is returned
    And file entries remain unchanged

  @happy
  Scenario: Read operations are allowed on signed packages
    Given a signed NovusPack package
    When ReadFile is called
    Then read operation succeeds
    And file content is returned
    And no immutability error is returned

  @happy
  Scenario: Signature addition is allowed on signed packages
    Given a signed NovusPack package
    When additional signature is added
    Then signature addition succeeds
    And signature is appended
    And existing signatures remain valid

  @error
  Scenario: Write protection prevents accidental signature invalidation
    Given a signed NovusPack package
    When content modification is attempted
    Then write protection prevents modification
    And all signatures remain valid
    And package integrity is maintained

  @happy
  Scenario: Unsigned packages allow all write operations
    Given an unsigned NovusPack package with SignatureOffset = 0
    When write operations are performed
    Then all write operations succeed
    And no immutability error is returned
    And package can be modified freely
