@domain:security @m2 @REQ-SEC-112 @REQ-SEC-113 @REQ-SEC-114 @REQ-SEC-115 @spec(api_security.md#421-isvalid-requirements) @spec(api_security.md#422-isexpired-requirements) @spec(api_security.md#423-expiration-semantics) @spec(api_security.md#424-validation-order)
Feature: EncryptionKey validity and expiration semantics

  @REQ-SEC-112 @happy
  Scenario: IsValid returns true only when all required fields are valid
    Given an EncryptionKey with key value set
    And KeyID is non-empty
    And CreatedAt is non-zero
    And KeyType is valid
    And ExpiresAt is nil or after CreatedAt
    When IsValid is evaluated
    Then IsValid returns true

  @REQ-SEC-112 @error
  Scenario: IsValid returns false when required fields are missing or invalid
    Given an EncryptionKey with missing key value or empty KeyID or zero CreatedAt
    When IsValid is evaluated
    Then IsValid returns false

  @REQ-SEC-113 @happy
  Scenario: IsExpired returns true only when ExpiresAt is set and now is at or after it
    Given an EncryptionKey with ExpiresAt set
    When current time is equal to or after ExpiresAt
    Then IsExpired returns true

  @REQ-SEC-114 @happy
  Scenario: ExpiresAt nil means never expires
    Given an EncryptionKey with ExpiresAt nil
    When IsExpired is evaluated
    Then IsExpired returns false

  @REQ-SEC-115 @happy
  Scenario: Validation order checks IsValid before IsExpired and both must pass
    Given an EncryptionKey
    When validating key usability
    Then IsValid is checked first
    And IsExpired is checked second
    And key is usable only when IsValid is true and IsExpired is false

