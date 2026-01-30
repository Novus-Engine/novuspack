@domain:generics @m2 @REQ-GEN-037 @spec(api_generics.md#1337-path-length-limits-and-warnings)
Feature: Path Length Limits

  @REQ-GEN-037 @happy
  Scenario: Accept paths up to 65,535 bytes
    Given a file with path length of 65535 bytes
    When AddFile is called
    Then path is accepted without error
    And path fits in PathEntry.PathLength uint16
    And file is added successfully
    And no format limit error occurs

  @REQ-GEN-037 @happy
  Scenario: Storage policy accepts any path up to format limit
    Given files with various path lengths up to 65535 bytes
    When files are added to package
    Then all paths are accepted
    And no artificial restrictions are imposed
    And format limit is the only hard limit
    And maximum flexibility is provided

  @REQ-GEN-043 @happy
  Scenario: Windows extraction with extended paths > 260 bytes
    Given a package with paths between 261 and 32767 bytes
    When package is extracted on Windows
    Then extended path syntax is automatically used
    And paths are prefixed with \\?\ automatically
    And extraction succeeds seamlessly
    And no user configuration required

  @REQ-GEN-043 @happy
  Scenario: Windows automatic extended path handling
    Given a file with 5000 byte path in package
    When file is extracted on Windows
    Then implementation uses \\?\ prefix automatically
    And regular paths used for <= 260 bytes for compatibility
    And extended paths used for > 260 bytes
    And seamless support up to ~32,767 bytes

  @REQ-GEN-037 @error
  Scenario: Windows extraction fails for paths > 32,767 bytes
    Given a package with path exceeding 32767 bytes
    When package is extracted on Windows
    Then extraction fails with error
    And error indicates Windows extended path limit exceeded
    And error message is clear about path length issue

  @REQ-GEN-037 @error
  Scenario: Linux extraction fails for paths > 4,096 bytes
    Given a package with path exceeding 4096 bytes
    When package is extracted on Linux
    Then extraction fails with error
    And error indicates PATH_MAX limit exceeded
    And error message states 4096 byte limit

  @REQ-GEN-037 @error
  Scenario: macOS extraction fails for paths > 1,024 bytes
    Given a package with path exceeding 1024 bytes
    When package is extracted on macOS
    Then extraction fails with error
    And error indicates PATH_MAX limit exceeded
    And error message states 1024 byte limit

  @REQ-GEN-037 @happy
  Scenario: Validation only at extraction time
    Given paths exceeding platform limits in package
    When package is stored and transferred
    Then no errors occur during storage
    And package remains valid
    And validation occurs only when extracting
    And platform-specific limits checked at extraction

  @REQ-GEN-043 @happy
  Scenario: Comparison with ZIP/TAR - automatic handling
    Given paths over 260 bytes
    When using NovusPack versus ZIP/TAR
    Then NovusPack handles extended paths automatically
    And ZIP/TAR require manual extended path enabling
    And NovusPack provides better user experience
    And no manual configuration needed
