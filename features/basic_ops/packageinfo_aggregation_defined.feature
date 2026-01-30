@domain:basic_ops @m2 @REQ-API_BASIC-209 @spec(api_basic_operations.md#3354-packageinfo-aggregation)
Feature: PackageInfo aggregation

  @REQ-API_BASIC-209 @happy
  Scenario: PackageInfo aggregates header-derived and computed information
    Given a package loaded from disk
    When PackageInfo is constructed in memory
    Then PackageInfo aggregation defines how information is combined
    And header-derived values are included in PackageInfo
    And computed statistics are included when applicable
    And aggregation avoids duplicating authoritative state
    And aggregation aligns with documented PackageInfo semantics

