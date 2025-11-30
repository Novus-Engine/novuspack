@domain:metadata @m2 @REQ-META-056 @REQ-META-057 @spec(api_metadata.md#111-writeto-parameters)
Feature: Metadata Parameter Specification

  @REQ-META-056 @happy
  Scenario: WriteTo accepts writer parameter
    Given a PackageComment with content
    And a writer
    When WriteTo is called with writer
    Then comment is written to writer
    And written data matches comment content
    And number of bytes written is returned

  @REQ-META-056 @happy
  Scenario: WriteTo writes complete comment content
    Given a PackageComment with content
    And a writer
    When WriteTo is called with writer
    Then all comment bytes are written
    And write operation completes successfully
    And writer receives UTF-8 encoded content

  @REQ-META-056 @happy
  Scenario: WriteTo handles empty comments
    Given an empty PackageComment
    And a writer
    When WriteTo is called with writer
    Then zero bytes are written
    And write operation completes successfully
    And no error is returned

  @REQ-META-057 @happy
  Scenario: ReadFrom accepts reader parameter
    Given a reader with comment data
    And a PackageComment
    When ReadFrom is called with reader
    Then comment is read from reader
    And comment content matches reader data
    And number of bytes read is returned

  @REQ-META-057 @happy
  Scenario: ReadFrom reads complete comment content
    Given a reader with comment data
    And a PackageComment
    When ReadFrom is called with reader
    Then all available data is read
    And read operation completes successfully
    And comment content is UTF-8 encoded

  @REQ-META-057 @happy
  Scenario: ReadFrom handles empty readers
    Given an empty reader
    And a PackageComment
    When ReadFrom is called with reader
    Then zero bytes are read
    And comment remains empty
    And no error is returned

  @REQ-META-056 @REQ-META-057 @error
  Scenario: WriteTo fails with nil writer
    Given a PackageComment with content
    And a nil writer
    When WriteTo is called with nil writer
    Then structured validation error is returned
    And error indicates invalid writer

  @REQ-META-056 @REQ-META-057 @error
  Scenario: ReadFrom fails with nil reader
    Given a nil reader
    And a PackageComment
    When ReadFrom is called with nil reader
    Then structured validation error is returned
    And error indicates invalid reader

  @REQ-META-057 @error
  Scenario: ReadFrom fails with invalid encoding
    Given a reader with invalid UTF-8 data
    And a PackageComment
    When ReadFrom is called with reader
    Then structured validation error is returned
    And error indicates encoding issue
