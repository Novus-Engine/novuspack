@skip @domain:writing @m2 @spec(api_writing.md#11-safewrite-method-signature)
Feature: Writing Definitions

# This feature captures key writing API constraints and expectations from the writing spec.
# Detailed runnable scenarios live in the dedicated writing feature files.

  @REQ-WRITE-012 @interface
  Scenario: SafeWrite method signature supports context and overwrite control
    Given a Package instance configured for writing to a target path
    When SafeWrite is invoked
    Then SafeWrite accepts a context for cancellation and deadlines
    And SafeWrite accepts an overwrite flag for controlling replacement behavior

  @REQ-WRITE-013 @constraint
  Scenario: SafeWrite uses same-directory temporary file for atomic replace
    Given a target package path on a filesystem that supports atomic rename
    When SafeWrite writes the package to disk
    Then the temporary file is created in the same directory as the target file
    And the final replace is atomic when overwriting an existing file
