@domain:file_mgmt @m2 @REQ-FILEMGMT-095 @spec(api_file_management.md#1112-iscompressed-returns)
Feature: IsCompressed Method

  @REQ-FILEMGMT-095 @happy
  Scenario: Iscompressed returns compression status
    Given a file entry
    When IsCompressed is called
    Then boolean compression status is returned
    And status indicates if file is compressed

  @REQ-FILEMGMT-095 @happy
  Scenario: IsCompressed returns true for compressed files
    Given a file entry with compression enabled
    When IsCompressed is called
    Then true is returned
    And compression status is correctly indicated

  @REQ-FILEMGMT-095 @happy
  Scenario: IsCompressed returns false for uncompressed files
    Given a file entry without compression
    When IsCompressed is called
    Then false is returned
    And uncompressed status is correctly indicated

  @REQ-FILEMGMT-095 @happy
  Scenario: IsCompressed checks compression type from file entry metadata
    Given a file entry
    When IsCompressed is called
    Then compression type is checked from file entry metadata
    And status is determined from metadata
    And result reflects actual compression state

  @REQ-FILEMGMT-095 @happy
  Scenario: IsCompressed provides access to file entry properties
    Given a file entry
    When IsCompressed is called
    Then file entry property is accessed
    And compression property is returned
    And access is efficient
