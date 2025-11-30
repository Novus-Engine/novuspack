@domain:security @m2 @REQ-SEC-063 @spec(security.md#831-comment-storage-security)
Feature: Comment Storage Security

  @REQ-SEC-063 @happy
  Scenario: Comment storage security provides immutable storage
    Given an open NovusPack package
    And a valid context
    And package with comments
    When comment storage security is examined
    Then comments are stored in immutable sections after signing
    And immutable storage prevents comment modification
    And comments cannot be changed after package signing

  @REQ-SEC-063 @happy
  Scenario: Comment storage security provides integrity protection
    Given an open NovusPack package
    And a valid context
    And package with comments
    When comment storage security is examined
    Then comments are protected by digital signatures
    And signature validation ensures comment integrity
    And comment tampering is detected by signatures

  @REQ-SEC-063 @happy
  Scenario: Comment storage security provides access control
    Given an open NovusPack package
    And a valid context
    And package with comments
    When comment storage security is examined
    Then comments are read-only after package creation
    And write access is restricted after creation
    And access control prevents unauthorized modifications

  @REQ-SEC-063 @happy
  Scenario: Comment storage security provides audit logging
    Given an open NovusPack package
    And a valid context
    And package with comments
    When comment storage security is examined
    Then all comment modifications are logged for security auditing
    And audit trail tracks comment changes
    And security events are recorded

  @REQ-SEC-063 @happy
  Scenario: Comment storage security maintains security throughout package lifecycle
    Given an open NovusPack package
    And a valid context
    And package with comments through lifecycle
    When comment storage security is maintained
    Then security protections remain active
    And immutable storage is maintained
    And integrity protection is maintained
    And access control is maintained
