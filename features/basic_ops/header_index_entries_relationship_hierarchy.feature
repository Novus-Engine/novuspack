@domain:basic_ops @m2 @REQ-API_BASIC-208 @spec(api_basic_operations.md#3353-header--index--entries)
Feature: Header to Index to Entries relationship

  @REQ-API_BASIC-208 @happy
  Scenario: Package structure hierarchy relates header, index, and entries
    Given a package loaded from disk
    When internal structure is represented in memory
    Then the header relates to the index structure
    And the index relates to the set of file entries
    And relationships support efficient file lookup and iteration
    And relationships remain consistent across mutations and writes
    And hierarchy representation aligns with the documented package format

