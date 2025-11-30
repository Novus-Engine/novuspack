@domain:generics @m2 @REQ-GEN-012 @spec(api_generics.md#13-collection-interface)
Feature: Generics Collection Interface

  @REQ-GEN-012 @happy
  Scenario: Collection interface provides type-safe collection operations
    Given a Collection interface
    When Collection is implemented for a type
    Then type-safe collection operations are provided
    And Add, Remove, Contains operations are available
    And Size, Clear, ToSlice operations are available

  @REQ-GEN-012 @happy
  Scenario: Collection Add adds items to collection
    Given a Collection instance
    And an item to add
    When Add is called with item
    Then item is added to collection
    And error is returned if addition fails
    And type-safe addition is provided

  @REQ-GEN-012 @happy
  Scenario: Collection Remove removes items from collection
    Given a Collection instance
    And an item to remove
    When Remove is called with item
    Then item is removed from collection
    And error is returned if removal fails
    And type-safe removal is provided

  @REQ-GEN-012 @happy
  Scenario: Collection Contains checks item presence
    Given a Collection instance
    And an item to check
    When Contains is called with item
    Then boolean indicates if item is in collection
    And type-safe checking is provided

  @REQ-GEN-012 @happy
  Scenario: Collection Size returns collection size
    Given a Collection instance
    When Size is called
    Then number of items in collection is returned
    And type-safe size querying is provided

  @REQ-GEN-012 @happy
  Scenario: Collection Clear removes all items
    Given a Collection instance with items
    When Clear is called
    Then all items are removed
    And collection is empty
    And type-safe clearing is provided

  @REQ-GEN-012 @happy
  Scenario: Collection ToSlice converts to slice
    Given a Collection instance with items
    When ToSlice is called
    Then slice of items is returned
    And type-safe conversion is provided
    And all collection items are included
