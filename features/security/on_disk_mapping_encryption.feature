@domain:security @m2 @REQ-SEC-110 @spec(api_security.md#33-on-disk-mapping)
Feature: On-disk mapping of encryption maps in-memory types to on-disk values

  @REQ-SEC-110 @happy
  Scenario: On-disk mapping of encryption is defined
    Given package file format and encryption types
    When encryption is persisted or read from disk
    Then on-disk mapping is defined by the package file format
    And in-memory encryption types map to on-disk values as specified
    And ML-KEM variant derivation from key material is applied
    And the behavior matches the on-disk mapping specification
