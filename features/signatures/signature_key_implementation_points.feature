@domain:signatures @m2 @REQ-SIG-024 @spec(api_signatures.md#123-key-implementation-points)
Feature: Signature Key Implementation Points

  @REQ-SIG-024 @happy
  Scenario: Key implementation points define critical implementation details
    Given a NovusPack package
    And a valid context
    When key implementation points are examined
    Then SignatureOffset in header points directly to first signature
    And additional signatures follow immediately after first signature
    And no separate signature index is needed
    And signatures are read sequentially
    And each signature is self-contained with its own metadata header
    And context supports cancellation

  @REQ-SIG-024 @happy
  Scenario: SignatureOffset points to first signature
    Given a NovusPack package
    And a valid context
    And a signed package
    When SignatureOffset is examined
    Then offset points directly to first signature location
    And offset enables direct signature access
    And offset is set when first signature is added

  @REQ-SIG-024 @happy
  Scenario: Signatures are stored sequentially without index
    Given a NovusPack package
    And a valid context
    And a package with multiple signatures
    When signatures are accessed
    Then signatures are read sequentially from SignatureOffset
    And no separate signature index is required
    And each signature follows immediately after previous signature

  @REQ-SIG-024 @happy
  Scenario: Each signature is self-contained
    Given a NovusPack package
    And a valid context
    And a signature
    When signature structure is examined
    Then signature contains its own metadata header
    And signature contains its own comment
    And signature contains its own data
    And signature is independently accessible

  @REQ-SIG-024 @error
  Scenario: Key implementation points handle invalid offsets
    Given a NovusPack package
    When invalid SignatureOffset is detected
    Then offset validation detects invalid values
    And appropriate errors are returned
