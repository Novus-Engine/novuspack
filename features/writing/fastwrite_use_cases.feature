@domain:writing @m2 @REQ-WRITE-019 @spec(api_writing.md#23-fastwrite-use-cases)
Feature: FastWrite Use Cases

  @REQ-WRITE-019 @happy
  Scenario: FastWrite use case: Incremental file updates
    Given an open NovusPack package
    And an existing package file exists
    And individual files are being added or modified
    When FastWrite is used for incremental updates
    Then FastWrite is appropriate for incremental updates
    And minimal I/O overhead is achieved
    And performance is optimized

  @REQ-WRITE-019 @happy
  Scenario: FastWrite use case: Large package modifications
    Given an open NovusPack package
    And an existing package file exists
    And the package is large (>1GB)
    And modifications are needed
    When FastWrite is used for large package modifications
    Then FastWrite is appropriate for large packages
    And memory efficiency is maintained
    And I/O efficiency is optimized

  @REQ-WRITE-019 @happy
  Scenario: FastWrite use case: Frequent updates
    Given an open NovusPack package
    And an existing package file exists
    And package requires frequent updates
    When FastWrite is used for frequent updates
    Then FastWrite is appropriate when performance is critical
    And frequent updates are efficient
    And performance is optimized

  @REQ-WRITE-019 @happy
  Scenario: FastWrite use case: Existing packages
    Given an open NovusPack package
    And an existing package file exists at target path
    When FastWrite is used for existing packages
    Then FastWrite is appropriate for existing packages
    And in-place updates are performed
    And efficiency is maintained

  @REQ-WRITE-019 @error
  Scenario: FastWrite cannot be used with unsigned packages requirement
    Given an open NovusPack package
    And the package has signatures (SignatureOffset > 0)
    When FastWrite is attempted
    Then FastWrite cannot be used with signed packages
    And error is returned
    And SafeWrite must be used instead

  @REQ-WRITE-019 @error
  Scenario: FastWrite cannot be used with uncompressed packages requirement
    Given an open NovusPack package
    And the package is compressed
    When FastWrite is attempted
    Then FastWrite cannot be used with compressed packages
    And error is returned
    And SafeWrite must be used instead
