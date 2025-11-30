@domain:writing @m2 @REQ-WRITE-006 @spec(api_writing.md#4-signed-file-write-operations)
Feature: Signed file write operations

  @happy
  Scenario: Write operations check SignatureOffset before proceeding
    Given a package
    When write operation is attempted
    Then SignatureOffset is checked first
    And if SignatureOffset > 0, package is identified as signed

  @error
  Scenario: Write operations are refused on signed packages by default
    Given a signed package with SignatureOffset > 0
    When write operation is attempted
    Then write operation is refused
    And structured immutability error is returned
    And error indicates signed package protection

  @happy
  Scenario: ClearSignatures flag allows write with signature removal
    Given a signed package
    When write operation is called with clearSignatures flag
    Then signatures are removed
    And write operation proceeds
    And new package is written without signatures

  @error
  Scenario: ClearSignatures requires explicit flag
    Given a signed package
    When write operation is attempted without clearSignatures flag
    Then operation fails
    And error indicates clearSignatures flag required
    And signatures are not removed

  @happy
  Scenario: Signed packages support only signature addition
    Given a signed package
    When signature addition is attempted
    Then signature addition succeeds
    And existing signatures remain valid
    And new signature is appended
