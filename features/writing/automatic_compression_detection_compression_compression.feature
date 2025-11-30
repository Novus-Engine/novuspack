@domain:writing @m2 @REQ-WRITE-044 @spec(api_writing.md#551-automatic-compression-detection)
Feature: Automatic Compression Detection

  @REQ-WRITE-044 @happy
  Scenario: Automatic compression detection identifies compression needs
    Given a NovusPack package
    When Write method is called without compression parameter
    Then automatic compression detection determines compression state
    And detection preserves current compression state by default
    And detection maintains compression consistency

  @REQ-WRITE-044 @happy
  Scenario: Automatic detection preserves compression state
    Given a NovusPack package
    And a compressed package
    When Write method is called automatically
    Then compressed input results in compressed output
    And compression state is preserved from input to output
    And state preservation maintains package consistency

  @REQ-WRITE-044 @happy
  Scenario: Automatic detection preserves uncompressed state
    Given a NovusPack package
    And an uncompressed package
    When Write method is called automatically
    Then uncompressed input results in uncompressed output
    And compression state is preserved from input to output
    And state preservation maintains package consistency

  @REQ-WRITE-044 @happy
  Scenario: Automatic detection handles new packages
    Given a new NovusPack package
    When Write method is called automatically
    Then new packages are uncompressed by default
    And default state enables direct write operations
    And default provides predictable behavior

  @REQ-WRITE-044 @happy
  Scenario: Automatic detection checks for signed packages
    Given a NovusPack package
    And a signed package
    When Write method is called automatically
    Then signed package check prevents compression
    And compression is refused if package is signed
    And check protects signature integrity

  @REQ-WRITE-044 @error
  Scenario: Automatic detection handles error conditions
    Given a NovusPack package
    And error conditions during detection
    When automatic compression detection encounters errors
    Then structured error is returned
    And error indicates detection failure
    And error follows structured error format
