@domain:file_mgmt @extraction @security @performance @REQ-FILEMGMT-396 @spec(api_file_mgmt_extraction.md#241-extraction-ordering-for-key-lifetime-minimization)
Feature: Extraction ordering minimizes decryption key lifetime

  As a package user
  I want multi-file extraction to prioritize encrypted files
  So that decryption keys can be cleared from memory as soon as possible

  @REQ-FILEMGMT-396 @happy
  Scenario: Multi-file extraction prioritizes encrypted files to minimize key lifetime
    Given a package with encrypted and unencrypted files
    And decryption keys are available
    When I extract the package
    Then encrypted files should be scheduled ahead of unencrypted files
    And per-file decryption key material should be cleared after each encrypted file is extracted
