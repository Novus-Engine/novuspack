@domain:writing @m2 @REQ-WRITE-046 @spec(api_writing.md#56-performance-considerations)
Feature: Writing Performance Considerations

  @REQ-WRITE-046 @happy
  Scenario: Performance considerations: Compressed packages have slower read speed
    Given a compressed NovusPack package
    When performance characteristics are examined
    Then read speed is slower than uncompressed packages
    And decompression overhead affects performance
    And direct access is not possible

  @REQ-WRITE-046 @happy
  Scenario: Performance considerations: Compressed packages have slower write speed
    Given a compressed NovusPack package
    When performance characteristics are examined
    Then write speed is slower than uncompressed packages
    And compression overhead affects performance
    And recompression is required after modifications

  @REQ-WRITE-046 @happy
  Scenario: Performance considerations: Compressed packages have lower disk usage
    Given a compressed NovusPack package
    When performance characteristics are examined
    Then disk usage is lower than uncompressed packages
    And compression reduces file size
    And storage efficiency is improved

  @REQ-WRITE-046 @happy
  Scenario: Performance considerations: Compressed packages have higher memory usage
    Given a compressed NovusPack package
    When performance characteristics are examined
    Then memory usage is higher than uncompressed packages
    And compression buffers are required
    And decompression buffers are required
    And memory overhead is present

  @REQ-WRITE-046 @happy
  Scenario: Performance considerations: Compressed packages have faster network transfer
    Given a compressed NovusPack package
    When performance characteristics are examined
    Then network transfer is faster than uncompressed packages
    And smaller file size reduces transfer time
    And bandwidth efficiency is improved

  @REQ-WRITE-046 @happy
  Scenario: Performance considerations: Uncompressed packages have faster read speed
    Given an uncompressed NovusPack package
    When performance characteristics are examined
    Then read speed is faster than compressed packages
    And no decompression overhead exists
    And direct file access is possible

  @REQ-WRITE-046 @happy
  Scenario: Performance considerations: Uncompressed packages have faster write speed
    Given an uncompressed NovusPack package
    When performance characteristics are examined
    Then write speed is faster than compressed packages
    And no compression overhead exists
    And direct file writing is possible
