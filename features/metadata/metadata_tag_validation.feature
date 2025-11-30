@domain:metadata @m2 @REQ-META-035 @spec(metadata.md#14-tag-validation)
Feature: Metadata Tag Validation

  @REQ-META-035 @happy
  Scenario: Tag validation validates key format
    Given an open NovusPack package
    And a file entry
    When tag with key is validated
    Then key must be valid UTF-8
    And key must be non-empty
    And key must be max 255 bytes
    And key validation passes for valid keys

  @REQ-META-035 @happy
  Scenario: Tag validation validates value type
    Given an open NovusPack package
    And a file entry
    When tag with value is validated
    Then value must match declared type
    And type-specific validation is performed
    And validation ensures type correctness
    And value validation passes for correct types

  @REQ-META-035 @happy
  Scenario: Tag validation validates JSON values
    Given an open NovusPack package
    And a file entry
    When JSON tag value is validated
    Then JSON syntax must be valid
    And JSON parsing must succeed
    And invalid JSON is rejected
    And valid JSON passes validation

  @REQ-META-035 @happy
  Scenario: Tag validation validates integer values
    Given an open NovusPack package
    And a file entry
    When integer tag value is validated
    Then value must be valid 64-bit signed integer
    And integer parsing must succeed
    And invalid integers are rejected
    And valid integers pass validation

  @REQ-META-035 @happy
  Scenario: Tag validation validates all supported types
    Given an open NovusPack package
    And a file entry
    When tags with various types are validated
    Then string type validation passes
    And integer type validation passes
    And float type validation passes
    And boolean type validation passes
    And all supported types are validated

  @REQ-META-035 @error
  Scenario: Tag validation fails with invalid key format
    Given an open NovusPack package
    And a file entry
    When tag with invalid key format is set
    Then structured validation error is returned
    And error indicates invalid key format
    And error specifies validation failure reason

  @REQ-META-035 @error
  Scenario: Tag validation fails with key exceeding max length
    Given an open NovusPack package
    And a file entry
    When tag with key exceeding 255 bytes is set
    Then structured validation error is returned
    And error indicates key length exceeded
    And error specifies maximum key length

  @REQ-META-035 @error
  Scenario: Tag validation fails with type mismatch
    Given an open NovusPack package
    And a file entry
    When tag value does not match declared type
    Then structured validation error is returned
    And error indicates type mismatch
    And error specifies expected and actual types

  @REQ-META-035 @error
  Scenario: Tag validation fails with invalid JSON
    Given an open NovusPack package
    And a file entry
    When JSON tag with invalid syntax is set
    Then structured validation error is returned
    And error indicates invalid JSON syntax
    And error specifies JSON parsing failure
