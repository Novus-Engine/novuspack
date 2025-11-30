@domain:file_format @m2 @REQ-FILEFMT-050 @spec(package_file_format.md#4123-multiple-path-support-with-per-path-metadata)
Feature: Multiple Path Support with Per-Path Metadata

  @REQ-FILEFMT-050 @happy
  Scenario: Multiple path support enables multiple paths pointing to same content
    Given a NovusPack package
    And file entry has multiple paths
    When multiple path support is used
    Then multiple paths point to same content
    And path aliasing is enabled
    And storage efficiency is improved

  @REQ-FILEFMT-050 @happy
  Scenario: Per-path metadata supports individual path attributes
    Given a NovusPack package
    And file entry has multiple paths
    When per-path metadata is used
    Then each path has its own metadata
    And path permissions can be different
    And path timestamps can be different
    And individual path attributes are maintained

  @REQ-FILEFMT-050 @happy
  Scenario: Multiple path support enables hard links and symbolic links
    Given a NovusPack package
    And file entry represents linked files
    When multiple path support is used
    Then hard links are supported
    And symbolic links are supported
    And link structure is preserved efficiently
