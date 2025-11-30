@domain:metadata @m2 @REQ-META-017 @REQ-META-035 @spec(metadata.md#12-tag-value-types)
Feature: Tag Value Types

  @REQ-META-017 @happy
  Scenario: Tag value types support basic types
    Given an open NovusPack package
    And a file entry
    When tags with basic types are set
    Then string type (0x00) is supported
    And integer type (0x01) is supported
    And float type (0x02) is supported
    And boolean type (0x03) is supported

  @REQ-META-017 @happy
  Scenario: Tag value types support structured data
    Given an open NovusPack package
    And a file entry
    When tags with structured data types are set
    Then JSON type (0x04) is supported
    And YAML type (0x05) is supported
    And StringList type (0x06) is supported

  @REQ-META-017 @happy
  Scenario: Tag value types support identifiers
    Given an open NovusPack package
    And a file entry
    When tags with identifier types are set
    Then UUID type (0x07) is supported
    And Hash type (0x08) is supported
    And Version type (0x09) is supported

  @REQ-META-017 @happy
  Scenario: Tag value types support time values
    Given an open NovusPack package
    And a file entry
    When tags with time types are set
    Then Timestamp type (0x0A) is supported
    And ISO8601 timestamp format is used

  @REQ-META-017 @happy
  Scenario: Tag value types support network communication
    Given an open NovusPack package
    And a file entry
    When tags with network types are set
    Then URL type (0x0B) is supported
    And Email type (0x0C) is supported

  @REQ-META-017 @happy
  Scenario: Tag value types support file system
    Given an open NovusPack package
    And a file entry
    When tags with file system types are set
    Then Path type (0x0D) is supported
    And MimeType type (0x0E) is supported

  @REQ-META-017 @happy
  Scenario: Tag value types support localization
    Given an open NovusPack package
    And a file entry
    When tags with localization types are set
    Then Language type (0x0F) is supported
    And ISO 639-1 language codes are used

  @REQ-META-017 @happy
  Scenario: Tag value types support NovusPack special files
    Given an open NovusPack package
    And a file entry
    When tags with NovusPack metadata types are set
    Then NovusPackMetadata type (0x10) is supported
    And special file references are supported

  @REQ-META-035 @happy
  Scenario: Tag validation validates key format
    Given an open NovusPack package
    And a file entry
    When tag with key is validated
    Then key must be valid UTF-8
    And key must be non-empty
    And key must be max 255 bytes

  @REQ-META-035 @happy
  Scenario: Tag validation validates value type
    Given an open NovusPack package
    And a file entry
    When tag with value is validated
    Then value must match declared type
    And type-specific validation is performed
    And validation ensures type correctness

  @REQ-META-035 @happy
  Scenario: Tag validation validates JSON values
    Given an open NovusPack package
    And a file entry
    When JSON tag value is validated
    Then JSON syntax must be valid
    And JSON parsing must succeed
    And invalid JSON is rejected

  @REQ-META-035 @happy
  Scenario: Tag validation validates integer values
    Given an open NovusPack package
    And a file entry
    When integer tag value is validated
    Then value must be valid 64-bit signed integer
    And integer parsing must succeed
    And invalid integers are rejected

  @REQ-META-035 @error
  Scenario: Tag validation fails with invalid key
    Given an open NovusPack package
    And a file entry
    When tag with invalid key is set
    Then structured validation error is returned
    And error indicates invalid key format

  @REQ-META-035 @error
  Scenario: Tag validation fails with type mismatch
    Given an open NovusPack package
    And a file entry
    When tag value does not match declared type
    Then structured validation error is returned
    And error indicates type mismatch
