@domain:metadata @m2 @REQ-META-026 @spec(metadata.md#129-reserved)
Feature: Reserved Metadata Value Types

  @REQ-META-026 @happy
  Scenario: Reserved value types are reserved for future use
    Given a NovusPack package
    When reserved tag value types are examined
    Then range 0x11-0xFF is reserved for future value types
    And reserved types cannot be used currently
    And reserved types are available for future expansion

  @REQ-META-026 @happy
  Scenario: Reserved types span from 0x11 to 0xFF
    Given a NovusPack package
    When reserved type range is examined
    Then 0x11 is the first reserved value type
    And 0xFF is the last reserved value type
    And range provides 239 reserved value types

  @REQ-META-026 @error
  Scenario: Reserved types cannot be used
    Given a NovusPack package
    When reserved value type is used
    Then reserved type is rejected
    And appropriate error indicates type is reserved
    And error follows structured error format

  @REQ-META-026 @happy
  Scenario: Reserved types enable future extensibility
    Given a NovusPack package
    When tag value type system is extended
    Then reserved types provide room for expansion
    And new types can be assigned from reserved range
