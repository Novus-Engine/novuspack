@domain:file_mgmt @REQ-FILEMGMT-443 @REQ-FILEMGMT-200 @spec(api_file_mgmt_addition.md#addfilefrommemory-memory-management) @spec(api_file_mgmt_addition.md#22-packageaddfilefrommemory-method)
Feature: AddFileFromMemory Memory Management

  The data parameter is a slice (reference type). Caller must ensure data
  remains valid until Write operation completes or data is processed
  (for encryption cases).

  @REQ-FILEMGMT-443 @happy
  Scenario: Data slice is reference type
    Given an open writable package
    And byte slice created by caller
    When AddFileFromMemory is called with path, data, and nil options
    Then implementation does not copy underlying data
    And slice header is passed through
    And caller retains responsibility for data validity

  @REQ-FILEMGMT-443 @happy
  Scenario: Caller must keep data valid until Write completes
    Given an open writable package
    And byte data added via AddFileFromMemory
    When Write has not yet been called
    Then caller must not modify or free the original data slice
    When Write completes or encryption processing consumes data
    Then data validity requirement is satisfied

  @REQ-FILEMGMT-443 @error
  Scenario: Invalid or freed data before Write may cause undefined behavior or error
    Given an open writable package
    And data added via AddFileFromMemory
    When caller invalidates data before Write or processing completes
    Then implementation may return error or exhibit undefined behavior per spec
    And contract is documented for caller
