@skip @domain:metadata @m2 @spec(metadata.md#134-inheritance-examples)
Feature: Metadata Usage Examples

# This feature captures usage-oriented metadata scenarios derived from the metadata spec.
# Detailed runnable scenarios live in the dedicated metadata feature files.

  @documentation
  Scenario: Path metadata supports inheritance of properties from parent paths
    Given path metadata entries exist for a parent directory and a child directory
    When effective tags are resolved for a file path under the child directory
    Then properties from the parent are inherited
    And more specific child properties override inherited values as specified

  @documentation
  Scenario: Inheritance priority influences which metadata entries win
    Given multiple path metadata entries could apply to the same file path
    When inheritance resolution is performed
    Then exact path matches take priority
    And higher priority entries override lower priority entries
