@domain:streaming @m2 @REQ-STREAM-065 @spec(api_streaming.md#223-bufferpool-struct)
Feature: BufferPool struct duplicate provides buffer pool structure alternative

  @REQ-STREAM-065 @happy
  Scenario: BufferPool struct duplicate provides structure alternative
    Given a BufferPool struct duplicate reference
    When buffer pool structure is used
    Then the struct provides buffer pool structure as specified
    And the alternative structure matches the specification
    And the behavior matches the BufferPool struct specification
    And compatibility is maintained where required
