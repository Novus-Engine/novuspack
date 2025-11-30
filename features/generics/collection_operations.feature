@domain:generics @m2 @REQ-GEN-004 @spec(api_generics.md#21-collection-operations)
Feature: Collection Operations

  @REQ-GEN-004 @happy
  Scenario: Filter returns items matching predicate
    Given a collection of items
    And a predicate function
    When Filter is called with items and predicate
    Then filtered items matching predicate are returned
    And type-safe filtering is provided
    And original collection is unchanged

  @REQ-GEN-004 @happy
  Scenario: Map transforms items from one type to another
    Given a collection of items of type T
    And a mapper function from T to U
    When Map is called with items and mapper
    Then transformed items of type U are returned
    And type-safe transformation is provided
    And original collection is unchanged

  @REQ-GEN-004 @happy
  Scenario: Find returns first item matching predicate
    Given a collection of items
    And a predicate function
    When Find is called with items and predicate
    Then first matching item is returned
    And boolean indicates if item was found
    And type-safe finding is provided

  @REQ-GEN-004 @happy
  Scenario: Reduce accumulates result using reducer function
    Given a collection of items
    And an initial accumulator value
    And a reducer function
    When Reduce is called with items, initial, and reducer
    Then accumulated result is returned
    And type-safe reduction is provided
    And reducer is applied to all items

  @REQ-GEN-004 @happy
  Scenario: Collection operations support type-safe generic patterns
    Given generic collection operations
    When operations are used with different types
    Then type safety is enforced at compile time
    And operations work with any comparable type
    And generic patterns are reusable
