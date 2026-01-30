@domain:generics @m2 @REQ-GEN-034 @spec(api_generics.md#1336-unicode-normalization)
Feature: Path NFC Normalization

  @REQ-GEN-034 @happy
  Scenario: Paths normalized to NFC before storage
    Given a file path with Unicode characters
    When the path is stored in the package
    Then the path is normalized to NFC form
    And NFC normalization occurs before validation
    And stored path uses composed Unicode form

  @REQ-GEN-034 @happy
  Scenario: NFD input converted to NFC
    Given a file path in NFD (decomposed) form
    When AddFile is called with NFD path
    Then path is converted to NFC before storage
    And stored path uses composed form
    And path is consistently stored in NFC

  @REQ-GEN-035 @happy
  Scenario: NFC normalization ensures consistent lookups
    Given a file added with Unicode path
    When the same path is queried in different Unicode forms
    Then file is found consistently
    And NFD queries match NFC stored paths
    And lookups work across different input forms

  @REQ-GEN-036 @happy
  Scenario: NFC normalization resolves macOS vs Windows differences
    Given a file path with accented characters
    When path is added from macOS system (NFD default)
    Then path is stored in NFC form
    And Windows systems (NFC default) can locate file
    And cross-platform compatibility is ensured

  @REQ-GEN-035 @happy
  Scenario: Consistent lookups across platforms
    Given a package created on macOS with Unicode paths
    When package is opened on Windows or Linux
    Then all paths are found correctly
    And no duplicate entries for visually identical paths
    And platform-specific normalization is handled

  @REQ-GEN-034 @happy
  Scenario: Example path normalization - cafe vs café
    Given a path containing "café" in NFD form
    When path is stored in package
    Then path is normalized to NFC "café"
    And both NFD and NFC inputs find same file
    And deduplication works correctly

  @REQ-GEN-034 @happy
  Scenario: NFC normalization implemented using Unicode norm package
    Given Go implementation of path normalization
    When NFC normalization is performed
    Then golang.org/x/text/unicode/norm package is used
    And normalization follows Unicode NFC standard
    And normalization is consistent and correct
