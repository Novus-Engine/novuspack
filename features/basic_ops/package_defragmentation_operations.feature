@domain:basic_ops @m1 @REQ-API_BASIC-011 @spec(api_basic_operations.md#72-package-defragmentation)
Feature: Package defragmentation operations

  @happy
  Scenario: Defragment optimizes package structure
    Given an open NovusPack package with unused space
    When Defragment is called
    Then unused space from deleted files is removed
    And file entries are reorganized for optimal access
    And data sections are compacted
    And file size is reduced

  @happy
  Scenario: Defragment preserves all package metadata
    Given an open NovusPack package with metadata
    When Defragment is called
    Then all package metadata is preserved
    And package comment is preserved
    And VendorID and AppID are preserved
    And file metadata is preserved

  @happy
  Scenario: Defragment preserves all signatures
    Given a signed open NovusPack package
    When Defragment is called
    Then all signatures are preserved
    And signatures remain valid
    And signature integrity is maintained

  @happy
  Scenario: Defragment updates internal indexes
    Given an open NovusPack package
    When Defragment is called
    Then file index is updated
    And internal references are updated
    And all offsets are recalculated

  @happy
  Scenario: Defragment reorganizes file entries for optimal access
    Given an open NovusPack package
    When Defragment is called
    Then file entries are reorganized
    And frequently accessed files are prioritized
    And package access performance is improved

  @error
  Scenario: Defragment fails if package is not open
    Given a closed NovusPack package
    When Defragment is called
    Then a structured validation error is returned

  @error
  Scenario: Defragment fails if package is read-only
    Given a read-only open NovusPack package
    When Defragment is called
    Then a structured validation error is returned

  @error
  Scenario: Defragment fails on I/O errors
    Given an open NovusPack package
    And file system I/O errors occur
    When Defragment is called
    Then a structured I/O error is returned

  @error
  Scenario: Defragment respects context cancellation
    Given an open NovusPack package
    And a cancelled context
    When Defragment is called
    Then a structured context error is returned
    And defragmentation is cancelled
