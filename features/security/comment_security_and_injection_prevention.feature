@skip @domain:security @m2 @spec(security.md#8-comment-security-and-injection-prevention)
Feature: Comment Security and Injection Prevention

# This feature captures high-level comment security expectations from the security specs.
# Detailed runnable scenarios live in the dedicated security feature files.

  @documentation
  Scenario: Comments are treated as text and never executed
    Given a package comment contains user-provided text
    When the comment is stored and later displayed
    Then the comment is treated as pure text data
    And the system does not execute scripts or commands embedded in the comment

  @documentation
  Scenario: Comment validation enforces UTF-8 and blocks common injection patterns
    Given a comment contains invalid UTF-8, null bytes, or a script tag sequence
    When the comment is validated
    Then validation fails
    And the unsafe content is rejected or sanitized according to the comment security rules
