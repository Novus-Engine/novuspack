@domain:core @m1 @REQ-CORE-003 @spec(api_core.md#0-overview)
Feature: Report Library and Spec Version

  @REQ-CORE-003 @happy
  Scenario: Version returns semantic library version and spec version
    Given the library is initialized
    When Version is requested
    Then semantic library version is reported
    And spec version is reported
    And both versions are available for query

  @REQ-CORE-003 @happy
  Scenario: Semantic library version follows semantic versioning format
    Given library version information
    When Version is queried
    Then semantic library version follows semantic versioning format (MAJOR.MINOR.PATCH)
    And version format is parseable
    And version enables compatibility checking

  @REQ-CORE-003 @happy
  Scenario: Spec version identifies specification version
    Given specification version information
    When Version is queried
    Then spec version identifies specification version
    And spec version enables format compatibility checking
    And spec version indicates supported format features

  @REQ-CORE-003 @happy
  Scenario: Version reporting enables compatibility verification
    Given version information
    When Version is checked
    Then compatibility verification is enabled
    And library and spec versions enable version checking
    And version information aids in troubleshooting
