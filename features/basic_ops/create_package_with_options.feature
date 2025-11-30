@domain:basic_ops @m1 @REQ-API_BASIC-007 @spec(api_basic_operations.md#43-create-with-options)
Feature: Create package with options

  @happy
  Scenario: CreateWithOptions sets initial package comment
    Given a new Package instance
    And an existing writable directory
    When CreateWithOptions is called with a comment and valid path
    Then the package is configured with the comment set
    And CommentSize and CommentStart are updated in header
    And package comment is accessible via GetInfo
    And package remains in memory

  @happy
  Scenario: CreateWithOptions sets initial VendorID
    Given a new Package instance
    And an existing writable directory
    When CreateWithOptions is called with VendorID and valid path
    Then the package is configured with VendorID set in header
    And VendorID is accessible via GetInfo
    And package remains in memory

  @happy
  Scenario: CreateWithOptions sets initial AppID
    Given a new Package instance
    And an existing writable directory
    When CreateWithOptions is called with AppID and valid path
    Then the package is configured with AppID set in header
    And AppID is accessible via GetInfo
    And package remains in memory

  @happy
  Scenario: CreateWithOptions stores file permissions
    Given a new Package instance
    And an existing writable directory
    When CreateWithOptions is called with Permissions and valid path
    Then file permissions are stored for later use
    And permissions match CreateOptions.Permissions
    And package remains in memory

  @happy
  Scenario: CreateWithOptions combines multiple options
    Given a new Package instance
    And an existing writable directory
    When CreateWithOptions is called with comment, VendorID, AppID, permissions, and valid path
    Then all options are applied correctly
    And package header reflects all specified options
    And package is configured successfully
    And package remains in memory

  @REQ-API_BASIC-020 @error
  Scenario: CreateWithOptions fails with invalid file path
    Given a new Package instance
    When CreateWithOptions is called with an invalid path
    Then a structured validation error is returned
    And package configuration fails

  @REQ-API_BASIC-020 @error
  Scenario: CreateWithOptions fails if target directory does not exist
    Given a new Package instance
    And a non-existent directory path
    When CreateWithOptions is called with path in that directory
    Then a structured validation error is returned
    And error indicates directory does not exist
    And package configuration fails

  @REQ-API_BASIC-020 @error
  Scenario: CreateWithOptions fails if target directory is not writable
    Given a new Package instance
    And an existing read-only directory
    When CreateWithOptions is called with path in that directory
    Then a structured validation error is returned
    And error indicates directory is not writable
    And package configuration fails

  @REQ-API_BASIC-017 @REQ-API_BASIC-019 @error
  Scenario: CreateWithOptions respects context cancellation
    Given a new Package instance
    And an existing writable directory
    And a cancelled context
    When CreateWithOptions is called
    Then a structured context error is returned
    And package configuration is cancelled

  @REQ-API_BASIC-007 @happy
  Scenario: CreateWithOptions uses Create internally for path validation
    Given a new Package instance
    And an existing writable directory
    When CreateWithOptions is called with valid path and options
    Then Create is called internally to validate path
    And path validation succeeds
    And options are applied to package
    And package remains in memory
