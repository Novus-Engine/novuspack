@domain:file_format @m2 @REQ-FILEFMT-051 @spec(package_file_format.md#4124-hash-based-content-identification)
Feature: Hash-Based Content Identification

  @REQ-FILEFMT-051 @happy
  Scenario: Hash-based content identification provides content addressing
    Given a NovusPack package
    And file entry has hash data
    When hash-based content identification is used
    Then content addressing is provided
    And files can be identified by content hash
    And content-based lookup is enabled

  @REQ-FILEFMT-051 @happy
  Scenario: Hash-based content identification supports deduplication
    Given a NovusPack package
    And file entries have content hashes
    When hash-based content identification is used
    Then deduplication can identify duplicate content
    And hash-based matching enables content reuse
    And storage optimization is supported

  @REQ-FILEFMT-051 @happy
  Scenario: Hash-based content identification supports integrity verification
    Given a NovusPack package
    And file entry has content hash
    When hash-based content identification is used
    Then content integrity can be verified
    And hash verification detects corruption
    And content authenticity is validated
