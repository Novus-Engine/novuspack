@domain:signatures @m2 @REQ-SIG-060 @spec(api_signatures.md#11-multiple-signature-management-incremental-signing)
Feature: 1.1 Multiple Signature Management (Incremental Signing) is specified and implemented

  @REQ-SIG-060 @happy
  Scenario: Multiple signature management supports incremental signing
    Given a package that supports multiple signatures
    When incremental signing is performed
    Then multiple signature management is specified and implemented
    And subsequent signatures can be added in order
    And the behavior matches the multiple signature management specification
    And validation supports incremental signature verification
