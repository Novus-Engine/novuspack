@domain:core @m2 @REQ-CORE-124 @spec(api_core.md#121-memory-versus-disk-side-effects)
Feature: Memory versus disk side effects define PackageWriter write operations

  @REQ-CORE-124 @happy
  Scenario: PackageWriter write operations separate memory and disk effects
    Given a package with pending in-memory changes
    When PackageWriter write operations are performed
    Then memory effects are applied as specified
    And disk persistence occurs only when write methods are invoked
    And the boundary between memory and disk is clearly defined
