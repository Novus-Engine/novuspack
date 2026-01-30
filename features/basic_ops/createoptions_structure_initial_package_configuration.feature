@domain:basic_ops @m2 @REQ-API_BASIC-132 @spec(api_basic_operations.md#631-createoptions-structure)
Feature: CreateOptions structure defines initial configuration

  @REQ-API_BASIC-132 @happy
  Scenario: CreateOptions exposes fields for initial package configuration
    Given package creation with options
    When CreateOptions is used to configure creation
    Then CreateOptions provides initial package configuration fields
    And fields represent supported configuration knobs for creation behavior
    And default values are well-defined when fields are not set
    And CreateWithOptions and related constructors accept CreateOptions-derived configuration
    And options affect package state prepared for later writing

