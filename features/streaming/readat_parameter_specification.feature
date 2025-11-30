@domain:streaming @m2 @REQ-STREAM-037 @spec(api_streaming.md#1534-readat-parameters)
Feature: ReadAt Parameter Specification

  @REQ-STREAM-037 @happy
  Scenario: ReadAt parameters define p buffer parameter
    Given an open NovusPack package
    And a valid context
    And file stream
    When ReadAt is called
    Then p parameter is buffer to read data into
    And buffer receives read data
    And buffer size determines read amount

  @REQ-STREAM-037 @happy
  Scenario: ReadAt parameters define off offset parameter
    Given an open NovusPack package
    And a valid context
    And file stream
    When ReadAt is called
    Then off parameter is offset to read from
    And offset specifies starting position
    And random access is supported

  @REQ-STREAM-037 @happy
  Scenario: ReadAt parameters support random access reading
    Given an open NovusPack package
    And a valid context
    And file stream with data
    When ReadAt is called with different offsets
    Then reading at offset 0 reads from start
    And reading at offset 1024 reads from position 1024
    And reading at various offsets is supported
    And random access enables flexible reading

  @REQ-STREAM-037 @happy
  Scenario: ReadAt parameters support context integration
    Given an open NovusPack package
    And a valid context
    And file stream
    When ReadAt is called
    Then operation accepts context.Context
    And context supports cancellation
    And context supports timeout handling

  @REQ-STREAM-037 @error
  Scenario: ReadAt parameters handle invalid offsets
    Given an open NovusPack package
    And a valid context
    And file stream
    And offset beyond file size
    When ReadAt is called with invalid offset
    Then appropriate error is returned
    And error indicates offset issue
    And error follows structured error format
