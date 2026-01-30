@domain:core @m2 @REQ-CORE-082 @spec(api_core.md#114-todisplaypath-function) @spec(api_core.md#todisplaypath-behavior) @spec(api_core.md#todisplaypath-usage)
Feature: ToDisplayPath function converts stored paths to display paths

  @REQ-CORE-082 @happy
  Scenario: ToDisplayPath converts stored paths for display
    Given a stored package path with leading slash
    When ToDisplayPath is called with the path
    Then the path is converted to a display-friendly form
    And leading slash is stripped for end-user display
    And the result is suitable for UI or logging
