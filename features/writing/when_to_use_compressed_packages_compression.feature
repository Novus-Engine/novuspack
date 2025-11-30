@domain:writing @m2 @REQ-WRITE-053 @spec(api_writing.md#581-when-to-use-compressed-packages)
Feature: When to Use Compressed Packages

  @REQ-WRITE-053 @happy
  Scenario: When to use compressed packages guides compression usage
    Given a NovusPack package
    When compression usage guidance is needed
    Then archival storage scenarios favor compressed packages
    And network distribution scenarios favor compressed packages
    And small file collections favor compressed packages
    And text-heavy content favors compressed packages

  @REQ-WRITE-053 @happy
  Scenario: Archival storage benefits from compressed packages
    Given a NovusPack package
    And archival storage requirements
    When compression decision is made
    Then long-term storage benefits from space savings
    And space is more important than speed for archival
    And compressed packages reduce storage costs
    And compression is appropriate for archival scenarios

  @REQ-WRITE-053 @happy
  Scenario: Network distribution benefits from compressed packages
    Given a NovusPack package
    And network distribution requirements
    When compression decision is made
    Then packages distributed over networks benefit from compression
    And smaller packages transfer faster
    And bandwidth savings improve distribution efficiency
    And compression is appropriate for network scenarios

  @REQ-WRITE-053 @happy
  Scenario: Small file collections benefit from compressed packages
    Given a NovusPack package
    And many small files
    When compression decision is made
    Then packages with many small files benefit from compression
    And package-level compression combines files efficiently
    And compression improves small file storage efficiency
    And compression is appropriate for small file collections

  @REQ-WRITE-053 @happy
  Scenario: Text-heavy content benefits from compressed packages
    Given a NovusPack package
    And primarily text or structured data
    When compression decision is made
    Then text and structured data compress well
    And compression ratio is high for text content
    And space savings are significant
    And compression is appropriate for text-heavy content

  @REQ-WRITE-053 @happy
  Scenario: Compression guidance considers multiple factors
    Given a NovusPack package
    And multiple usage factors
    When compression guidance is consulted
    Then guidance considers archival vs access patterns
    And guidance considers content type and compressibility
    And guidance considers network transfer requirements
    And guidance helps make informed compression decisions
