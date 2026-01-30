@domain:basic_ops @m2 @REQ-API_BASIC-142 @spec(api_basic_operations.md#952-package-isreadonly-method)
Feature: Package.IsReadOnly method

  @REQ-API_BASIC-142 @happy
  Scenario: IsReadOnly reports whether the package is in read-only mode
    Given a package opened for reading
    When IsReadOnly is called
    Then it reports true
    Given a package opened for writing
    When IsReadOnly is called
    Then it reports false
    And IsReadOnly reflects the access mode of the current open session

