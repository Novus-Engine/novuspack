"""
Core constants and path/workspace helpers for validation scripts.
"""

import importlib
import types
from pathlib import Path
from typing import Optional, List

# Standard directory names used across validation scripts
DOCS_DIR = 'docs'
TECH_SPECS_DIR = 'tech_specs'
REQUIREMENTS_DIR = 'requirements'
FEATURES_DIR = 'features'


def get_validation_exit_code(has_errors, no_fail=False):
    """
    Get the appropriate exit code for validation scripts.

    Args:
        has_errors: True if validation errors were found, False otherwise
        no_fail: If True, always return 0 (even if errors were found)

    Returns:
        0 if no errors found or no_fail is True, 1 if errors were found
    """
    if no_fail:
        return 0
    return 0 if not has_errors else 1


def get_workspace_root() -> Path:
    """
    Get the workspace root directory (parent of scripts directory).

    Returns:
        Path to workspace root
    """
    # From scripts/lib/validation/_core.py: validation -> lib -> scripts -> repo
    script_dir = Path(__file__).parent
    return script_dir.parent.parent.parent


def import_module_with_fallback(module_name: str, _script_dir: Path) -> types.ModuleType:
    """
    Import a module by name.

    Args:
        module_name: Name of module to import (e.g., '_validation_utils')
        _script_dir: Directory containing the module file (unused)

    Returns:
        Imported module
    """
    return importlib.import_module(module_name)


def parse_paths(path_str: Optional[str]) -> Optional[List[str]]:
    """
    Parse comma-separated path string into list of paths.

    Args:
        path_str: Comma-separated string of paths, or None

    Returns:
        List of trimmed path strings, or None if path_str is None/empty
    """
    if not path_str:
        return None
    return [p.strip() for p in path_str.split(',') if p.strip()]
