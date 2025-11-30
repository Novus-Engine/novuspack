@domain:writing @m2 @REQ-WRITE-054 @spec(api_writing.md#582-when-to-use-uncompressed-packages)
Feature: When to Use Uncompressed Packages

  @REQ-WRITE-054 @happy
  Scenario: When to use uncompressed packages guides uncompressed usage
    Given a NovusPack package
    When uncompressed usage guidance is needed
    Then frequent access scenarios favor uncompressed packages
    And large binary files favor uncompressed packages
    And development workflow scenarios favor uncompressed packages
    And memory-constrained systems favor uncompressed packages

  @REQ-WRITE-054 @happy
  Scenario: Frequent access benefits from uncompressed packages
    Given a NovusPack package
    And frequent access requirements
    When compression decision is made
    Then packages accessed frequently benefit from uncompressed format
    And speed is critical for frequent access
    And decompression overhead is avoided
    And uncompressed packages are appropriate for frequent access

  @REQ-WRITE-054 @happy
  Scenario: Large binary files benefit from uncompressed packages
    Given a NovusPack package
    And large binary files
    When compression decision is made
    Then large binary files may not compress well
    And compression overhead may outweigh benefits
    And uncompressed format avoids unnecessary processing
    And uncompressed packages are appropriate for large binaries

  @REQ-WRITE-054 @happy
  Scenario: Development workflow benefits from uncompressed packages
    Given a NovusPack package
    And development workflow requirements
    When compression decision is made
    Then packages modified frequently during development benefit from uncompressed
    And faster write operations improve development experience
    And immediate access supports iterative development
    And uncompressed packages are appropriate for development

  @REQ-WRITE-054 @happy
  Scenario: Memory constraints favor uncompressed packages
    Given a NovusPack package
    And limited memory for compression operations
    When compression decision is made
    Then systems with limited memory benefit from uncompressed format
    And compression buffers increase memory requirements
    And uncompressed format reduces memory pressure
    And uncompressed packages are appropriate for memory-constrained systems

  @REQ-WRITE-054 @happy
  Scenario: Uncompressed guidance considers speed requirements
    Given a NovusPack package
    And speed-critical scenarios
    When uncompressed guidance is consulted
    Then guidance considers access frequency
    And guidance considers write frequency
    And guidance considers performance priorities
    And guidance helps make informed uncompressed decisions
