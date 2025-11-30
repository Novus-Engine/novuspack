@domain:compression @m2 @REQ-COMPR-062 @spec(api_package_compression.md#1023-compressing-signed-packages-workflow)
Feature: Compressing Signed Packages Workflow

  @REQ-COMPR-062 @happy
  Scenario: Workflow removes signatures first if package is signed
    Given a signed NovusPack package that needs compression
    When proper workflow is followed
    Then signatures are removed first
    And package is prepared for compression
    And signatures are cleared before compression

  @REQ-COMPR-062 @happy
  Scenario: Workflow allows changes to package after signature removal
    Given a signed package with signatures removed
    When workflow continues
    Then changes can be made to package
    And package modifications are allowed
    And package content can be updated

  @REQ-COMPR-062 @happy
  Scenario: Workflow allows recompression after changes
    Given a modified package that was previously signed
    When compression is desired
    Then package can be recompressed
    And compression operation succeeds
    And package is compressed

  @REQ-COMPR-062 @happy
  Scenario: Workflow requires re-signing after compression
    Given a compressed package that was previously signed
    When workflow completes
    Then package must be re-signed
    And re-signing provides new signatures
    And new signatures validate compressed content
