# Development Documentation Directory

This directory contains ALL development documentation for the NovusPack project.

## Technical Specifications Entry

For the full technical specifications, start here:

- `docs/tech_specs/_main.md`

## **CRITICAL REQUIREMENTS**

### **Directory Location**

- **ALL documentation MUST be created ONLY in this `dev_docs/` directory**
- **NO documents should be created outside this directory**
- **NO exceptions are allowed**

### File Naming Convention

ALL documents MUST follow this exact format:

```text
YYYY-MM-DD_document_type_description.md
```

#### Required Elements

- **Date:** ISO format (YYYY-MM-DD)
- **Document Type:** Clear type identifier (plan, analysis, test, review, implementation, etc.)
- **Description:** Specific, descriptive content summary
- **Extension:** Must be `.md` (Markdown format)

#### Valid Examples

- `2024-01-15_development_plan_mlkem_implementation.md`
- `2024-01-15_test_analysis_path_normalization.md`
- `2024-01-15_code_review_quantum_safe_encryption.md`
- `2024-01-15_implementation_notes_file_validation.md`

#### Invalid Examples (DO NOT USE)

- `plan.md` (missing date and description)
- `test_results.md` (missing date and specific description)
- `2024-1-15_notes.md` (incorrect date format)
- `notes.txt` (wrong extension and missing required elements)

### **Document Types**

- **Planning:** Development plans, architecture decisions, design notes
- **Analysis:** Code analysis, gap analysis, dependency mapping
- **Test:** Test plans, test results, test coverage reports
- **Implementation:** Development progress, technical decisions, implementation details
- **Review:** Code reviews, check-in reports, validation results
- **Troubleshooting:** Bug reports, issue tracking, resolution notes

### **Verification Process**

1. **Before creating ANY document:** Verify the filename follows the convention
2. **Before committing:** Verify all filenames are compliant
3. **If non-compliant:** Rename files before committing
4. **No exceptions:** Non-compliant files will be rejected

### **Directory Structure**

```text
dev_docs/
├── README.md (this file)
├── YYYY-MM-DD_document_type_description.md
├── YYYY-MM-DD_document_type_description.md
└── ...
```

## **Compliance is Mandatory**

- **Product Owner Authority:** These requirements are non-negotiable
- **No Deviations:** All documents must follow these rules exactly
- **Quality First:** Proper documentation is required for project acceptance
- **Version Control:** All documentation must be committed with proper naming
