@domain:security @m2 @REQ-SEC-061 @spec(security.md#823-sanitization-methods)
Feature: Input Sanitization Methods

  @REQ-SEC-061 @happy
  Scenario: Sanitization provides HTML escaping
    Given an open NovusPack package
    And a valid context
    And comment data with HTML characters
    When sanitization is performed
    Then HTML characters (<, >, &, ", ') are converted to HTML entities
    And HTML escaping prevents script injection
    And sanitized content is safe for display

  @REQ-SEC-061 @happy
  Scenario: Sanitization provides URL encoding
    Given an open NovusPack package
    And a valid context
    And comment data in URL contexts
    When sanitization is performed
    Then special characters are URL encoded
    And URL encoding prevents URL-based attacks
    And sanitized content is safe for URL contexts

  @REQ-SEC-061 @happy
  Scenario: Sanitization provides character removal
    Given an open NovusPack package
    And a valid context
    And comment data with dangerous characters
    When sanitization is performed
    Then dangerous characters are removed entirely
    And removed characters cannot cause security issues
    And sanitized content contains only safe characters

  @REQ-SEC-061 @happy
  Scenario: Sanitization provides character replacement
    Given an open NovusPack package
    And a valid context
    And comment data with dangerous characters
    When sanitization is performed
    Then dangerous characters are replaced with safe alternatives
    And character replacement maintains content readability
    And replaced content is safe for processing

  @REQ-SEC-061 @happy
  Scenario: Sanitization provides length truncation
    Given an open NovusPack package
    And a valid context
    And comment data exceeding length limits
    When sanitization is performed
    Then overly long content is truncated to safe limits
    And truncation preserves content up to safe limit
    And truncated content is within maximum length

  @REQ-SEC-061 @happy
  Scenario: Sanitization applies appropriate method based on content type
    Given an open NovusPack package
    And a valid context
    And comment data with various content types
    When sanitization is performed
    Then HTML content uses HTML escaping
    And URL content uses URL encoding
    And dangerous content uses character removal or replacement
    And oversized content uses length truncation
