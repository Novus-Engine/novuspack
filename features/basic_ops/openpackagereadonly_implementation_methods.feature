@domain:basic_ops @m2 @REQ-API_BASIC-221 @spec(api_basic_operations.md#7215-implementation-methods) @spec(api_basic_operations.md#116-readonlypackage-implementation-methods)
Feature: OpenPackageReadOnly implementation methods

  @REQ-API_BASIC-221 @happy
  Scenario: OpenPackageReadOnly reuses OpenPackage and wraps the returned Package
    Given a package open operation in read-only mode
    When OpenPackageReadOnly is implemented
    Then it reuses the OpenPackage logic to parse and load the package
    And it wraps the returned Package in a read-only enforcement wrapper
    And it avoids duplicating parsing and validation logic
    And it preserves correct in-memory state initialization
    And it returns a Package interface that enforces read-only behavior

