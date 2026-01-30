@domain:signatures @m2 @v2 @REQ-SIG-073 @REQ-SIG-074 @REQ-SIG-075 @REQ-SIG-076 @spec(api_signatures.md#421-isvalid-requirements) @spec(api_signatures.md#422-isexpired-requirements) @spec(api_signatures.md#423-expiration-semantics) @spec(api_signatures.md#424-validation-order)
Feature: SigningKey validity and expiration semantics

  @REQ-SIG-073 @happy
  Scenario: IsValid returns true only when required fields are valid
    Given a SigningKey with key set
    And KeyID is non-empty
    And CreatedAt is non-zero
    And KeyType is valid
    And ExpiresAt is nil or after CreatedAt
    When IsValid is evaluated
    Then IsValid returns true

  @REQ-SIG-074 @happy
  Scenario: IsExpired returns true only when ExpiresAt is set and now is at or after it
    Given a SigningKey with ExpiresAt set
    When current time is equal to or after ExpiresAt
    Then IsExpired returns true

  @REQ-SIG-075 @happy
  Scenario: ExpiresAt nil means never expires and CreatedAt zero is invalid
    Given a SigningKey with ExpiresAt nil
    When IsExpired is evaluated
    Then IsExpired returns false
    Given a SigningKey with CreatedAt equal to the zero time
    When IsValid is evaluated
    Then IsValid returns false

  @REQ-SIG-076 @happy
  Scenario: Validation order checks IsValid then IsExpired and both must pass
    Given a SigningKey
    When validating signing key usability
    Then IsValid is checked first
    And IsExpired is checked second
    And the key is usable only when IsValid is true and IsExpired is false

