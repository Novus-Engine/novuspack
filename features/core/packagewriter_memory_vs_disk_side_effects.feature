@domain:core @m2 @REQ-CORE-054 @spec(api_core.md#121-memory-versus-disk-side-effects)
Feature: PackageWriter memory versus disk side effects are clearly defined

  @REQ-CORE-054 @happy
  Scenario: PackageWriter operations separate in-memory changes from disk persistence
    Given an opened package with pending in-memory changes
    When a PackageWriter operation is performed
    Then changes are reflected in memory immediately
    And changes are not persisted to disk until a write operation is invoked
    And write operations define the persistence boundary
