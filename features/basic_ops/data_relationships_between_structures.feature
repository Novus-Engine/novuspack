@domain:basic_ops @m2 @REQ-API_BASIC-125 @spec(api_basic_operations.md#335-data-relationships)
Feature: Data relationships

  @REQ-API_BASIC-125 @happy
  Scenario: Data relationships define how core package structures relate
    Given a package loaded from disk
    When data relationships are established in memory
    Then relationships between header, index, file entries, and metadata are defined
    And data relationships support efficient lookup and operations
    And relationship invariants are maintained as the package mutates
    And relationships avoid duplication of authoritative state
    And relationships support validation and reporting

