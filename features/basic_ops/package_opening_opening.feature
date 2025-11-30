@domain:basic_ops @m2 @REQ-API_BASIC-039 @spec(api_basic_operations.md#5-package-opening)
Feature: Package Opening

  @REQ-API_BASIC-039 @happy
  Scenario: Package opening uses Open method
    Given a NovusPack package instance
    And a valid context
    And an existing package file
    When Open is called with file path
    Then package file is opened for reading and writing
    And package is ready for operations

  @REQ-API_BASIC-039 @happy
  Scenario: Package opening uses OpenWithValidation method
    Given a NovusPack package instance
    And a valid context
    And an existing package file
    When OpenWithValidation is called
    Then package file is opened
    And full package validation is performed
    And package integrity is verified before operations

  @REQ-API_BASIC-039 @happy
  Scenario: Package opening validates header and structure
    Given a NovusPack package instance
    And a valid context
    And an existing valid package file
    When package is opened
    Then package header is validated
    And basic package structure is validated
    And package format is confirmed

  @REQ-API_BASIC-039 @happy
  Scenario: Package opening loads package metadata
    Given a NovusPack package instance
    And a valid context
    And an existing package file
    When package is opened
    Then package metadata is loaded
    And comment, VendorID, AppID are available
    And package information is accessible

  @REQ-API_BASIC-039 @happy
  Scenario: Package opening prepares file operations
    Given a NovusPack package instance
    And a valid context
    And an existing package file
    When package is opened
    Then file entries are read
    And file operations are prepared
    And package state is set up for subsequent operations

  @REQ-API_BASIC-039 @error
  Scenario: Package opening handles file not found errors
    Given a NovusPack package instance
    And a valid context
    And a non-existent package file path
    When Open is called
    Then validation error is returned
    And error indicates package file not found

  @REQ-API_BASIC-039 @error
  Scenario: Package opening handles invalid format errors
    Given a NovusPack package instance
    And a valid context
    And a file with corrupted or invalid format
    When Open is called
    Then validation error is returned
    And error indicates corrupted or invalid format

  @REQ-API_BASIC-039 @error
  Scenario: Package opening handles unsupported version errors
    Given a NovusPack package instance
    And a valid context
    And a package file with unsupported version
    When Open is called
    Then unsupported error is returned
    And error indicates package version not supported
