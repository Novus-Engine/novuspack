@domain:compression @m2 @REQ-COMPR-083 @spec(api_package_compression.md#1341-compression)
Feature: Compression memory usage characteristics

  @REQ-COMPR-083 @happy
  Scenario: Compression memory usage characteristics are defined
    Given compression workloads for packages of varying size
    When memory usage is evaluated
    Then compression memory consumption characteristics are defined
    And memory usage depends on strategy and buffer sizing
    And streaming workflows bound peak memory usage by buffers
    And memory usage expectations inform configuration choices
    And memory usage characteristics are considered a non-functional concern

