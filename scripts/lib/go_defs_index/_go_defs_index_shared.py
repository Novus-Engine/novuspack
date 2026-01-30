"""
Shared helpers for Go definitions index tooling.
"""

from __future__ import annotations

from typing import Dict


_IMPLEMENTATION_TO_INTERFACE: Dict[str, str] = {
    "filePackage": "Package",
    "readOnlyPackage": "Package",
    "readOnlyPackageImpl": "Package",
    "packageReader": "PackageReader",
    "packageWriter": "PackageWriter",
}


def map_implementation_to_interface(receiver_type: str) -> str:
    """
    Map implementation types to their interface types.

    For example, filePackage -> Package.
    """
    if not receiver_type:
        return receiver_type
    mapped = _IMPLEMENTATION_TO_INTERFACE.get(receiver_type)
    if mapped:
        return mapped
    lower = receiver_type.lower()
    for name, interface in _IMPLEMENTATION_TO_INTERFACE.items():
        if lower == name.lower():
            return interface
    return receiver_type
