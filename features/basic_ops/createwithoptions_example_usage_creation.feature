@domain:basic_ops @m2 @REQ-API_BASIC-012 @REQ-API_BASIC-014 @spec(api_basic_operations.md#434-createwithoptions-example-usage)
Feature: CreateWithOptions example usage

  @REQ-API_BASIC-012 @REQ-API_BASIC-014 @happy
  Scenario: CreateWithOptions example demonstrates package creation with options
    Given a package needs to be created with specific options
    When CreateWithOptions is called with Comment, VendorID, and AppID
    Then package is created in memory
    And options are applied to package
    And package comment is set
    And VendorID is set
    And AppID is set

  @REQ-API_BASIC-012 @REQ-API_BASIC-014 @happy
  Scenario: CreateWithOptions example shows error checking pattern
    Given CreateWithOptions operation
    When CreateWithOptions is called
    And error return value is checked
    Then error is handled appropriately
    And error checking pattern is demonstrated

  @REQ-API_BASIC-012 @REQ-API_BASIC-014 @happy
  Scenario: CreateWithOptions example demonstrates configuration before writing
    Given CreateWithOptions is used
    When package is configured
    Then package remains in memory
    And package is not written to disk
    And Write method must be called to save
    And example demonstrates proper workflow
