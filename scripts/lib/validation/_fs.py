"""File and path discovery for validation scripts."""

import sys
from pathlib import Path
from typing import Dict, List, Optional, Set


_DEFAULT_EXCLUDE_DIRS: Set[str] = {
    'node_modules', 'vendor', 'tmp', '.git', '.venv', 'venv',
    '__pycache__', '.pytest_cache', 'dist', 'build',
    '.idea', '.vscode', '.cache'
}


def is_in_dot_directory(path: Path) -> bool:
    """
    Check if a path contains any directory starting with '.'.

    Args:
        path: Path object to check

    Returns:
        True if path contains any directory starting with '.' (except '.' itself), False otherwise
    """
    for part in path.parts:
        if part.startswith('.') and part != '.':
            return True
    return False


def _resolve_search_dir(
    root_dir: Optional[Path],
    default_dir: Optional[Path],
) -> Path:
    """Resolve search directory from root_dir and default_dir."""
    if root_dir is not None:
        return root_dir
    if default_dir is not None:
        return default_dir
    return Path('.')


def _collect_md_from_target_paths(
    target_paths: List[str],
    md_files: List[Path],
    verbose: bool,
) -> None:
    """Append markdown files from target paths to md_files."""
    for target_path in target_paths:
        target = Path(target_path)
        if not target.exists():
            if verbose:
                print(f"Warning: Target path does not exist: {target_path}", file=sys.stderr)
            continue
        if target.is_file():
            if target.suffix == '.md' and not is_in_dot_directory(target):
                md_files.append(target)
            elif verbose:
                print(
                    f"Warning: Target file is not a markdown file: {target_path}",
                    file=sys.stderr
                )
        else:
            for md_file in target.rglob('*.md'):
                if not is_in_dot_directory(md_file):
                    md_files.append(md_file)


def _collect_md_from_search_dir(
    search_dir: Path,
    exclude_dirs: Set[str],
    default_dir: Optional[Path],
    *,
    root_dir: Optional[Path],
    md_files: List[Path],
    verbose: bool,
) -> None:
    """Append markdown files from search_dir to md_files."""
    if not search_dir.exists():
        if verbose:
            print(f"Error: Search directory does not exist: {search_dir}", file=sys.stderr)
        return
    if default_dir is not None and root_dir is None:
        md_files.extend(
            f for f in sorted(search_dir.glob('*.md'))
            if not is_in_dot_directory(f)
        )
    else:
        for md_file in search_dir.rglob('*.md'):
            if any(excluded in md_file.parts for excluded in exclude_dirs):
                continue
            if is_in_dot_directory(md_file):
                continue
            md_files.append(md_file)


def find_markdown_files(
    target_paths: Optional[List[str]] = None,
    root_dir: Optional[Path] = None,
    default_dir: Optional[Path] = None,
    *,
    exclude_dirs: Optional[Set[str]] = None,
    verbose: bool = False,
    return_strings: bool = False,
) -> List[Path]:
    """
    Find markdown files in the repository or target paths.

    Args:
        target_paths: Optional list of specific files or directories to check
        root_dir: Root directory to search from (when target_paths is None)
        default_dir: Default directory to search if target_paths is None and root_dir is None
        exclude_dirs: Set of directory names to exclude when scanning root_dir
        verbose: Whether to show detailed progress
        return_strings: If True, return list of strings instead of Path objects

    Returns:
        List of Path objects (or strings if return_strings=True) for markdown files found
    """
    md_files: List[Path] = []
    exclude = exclude_dirs if exclude_dirs is not None else _DEFAULT_EXCLUDE_DIRS
    if target_paths:
        _collect_md_from_target_paths(target_paths, md_files, verbose)
    else:
        search_dir = _resolve_search_dir(root_dir, default_dir)
        _collect_md_from_search_dir(
            search_dir, exclude, default_dir,
            root_dir=root_dir, md_files=md_files, verbose=verbose,
        )
    if return_strings:
        return sorted([str(f) for f in md_files])
    return sorted(md_files)


def _collect_feature_from_target_paths(
    target_paths: List[str],
    feature_files: List[Path],
    verbose: bool,
) -> None:
    """Append feature files from target paths to feature_files."""
    for target_path in target_paths:
        target = Path(target_path)
        if not target.exists():
            if verbose:
                print(f"Warning: Target path does not exist: {target_path}", file=sys.stderr)
            continue
        if target.is_file():
            if target.suffix == '.feature' and not is_in_dot_directory(target):
                feature_files.append(target)
            elif verbose:
                print(
                    f"Warning: Target file is not a .feature file: {target_path}",
                    file=sys.stderr
                )
        else:
            for feature_file in target.rglob('*.feature'):
                if not is_in_dot_directory(feature_file):
                    feature_files.append(feature_file)


def _collect_feature_from_search_dir(
    search_dir: Path,
    exclude_dirs: Set[str],
    default_dir: Optional[Path],
    *,
    root_dir: Optional[Path],
    feature_files: List[Path],
    verbose: bool,
) -> None:
    """Append feature files from search_dir to feature_files."""
    if not search_dir.exists():
        if verbose:
            print(f"Error: Search directory does not exist: {search_dir}", file=sys.stderr)
        return
    if default_dir is not None and root_dir is None:
        feature_files.extend(
            f for f in sorted(search_dir.rglob('*.feature'))
            if not is_in_dot_directory(f)
            and not any(excluded in f.parts for excluded in exclude_dirs)
        )
    else:
        for feature_file in search_dir.rglob('*.feature'):
            if any(excluded in feature_file.parts for excluded in exclude_dirs):
                continue
            if is_in_dot_directory(feature_file):
                continue
            feature_files.append(feature_file)


def find_feature_files(
    target_paths: Optional[List[str]] = None,
    root_dir: Optional[Path] = None,
    default_dir: Optional[Path] = None,
    *,
    exclude_dirs: Optional[Set[str]] = None,
    verbose: bool = False,
    return_strings: bool = False,
) -> List[Path]:
    """
    Find feature files (.feature) in the repository or target paths.

    Args:
        target_paths: Optional list of specific files or directories to check
        root_dir: Root directory to search from (when target_paths is None)
        default_dir: Default directory to search if target_paths is None and root_dir is None
        exclude_dirs: Set of directory names to exclude when scanning root_dir
        verbose: Whether to show detailed progress
        return_strings: If True, return list of strings instead of Path objects

    Returns:
        List of Path objects (or strings if return_strings=True) for feature files found
    """
    feature_files: List[Path] = []
    exclude = exclude_dirs if exclude_dirs is not None else _DEFAULT_EXCLUDE_DIRS
    if target_paths:
        _collect_feature_from_target_paths(target_paths, feature_files, verbose)
    else:
        search_dir = _resolve_search_dir(root_dir, default_dir)
        _collect_feature_from_search_dir(
            search_dir, exclude, default_dir,
            root_dir=root_dir, feature_files=feature_files, verbose=verbose,
        )
    if return_strings:
        return sorted([str(f) for f in feature_files])
    return sorted(feature_files)


class FileContentCache:
    """
    Cache for file contents to avoid repeated reads.

    This class provides efficient caching of file contents to reduce I/O overhead
    when the same files are read multiple times during validation.
    """

    def __init__(self):
        """Initialize an empty cache."""
        self._cache: Dict[Path, str] = {}
        self._lines_cache: Dict[Path, List[str]] = {}

    def get_content(self, file_path: Path) -> str:
        """
        Get file content, using cache if available.

        Args:
            file_path: Path to the file to read

        Returns:
            File content as string

        Raises:
            IOError: If file cannot be read
        """
        if file_path not in self._cache:
            self._cache[file_path] = file_path.read_text(encoding='utf-8')
        return self._cache[file_path]

    def get_lines(self, file_path: Path) -> List[str]:
        """
        Get file content as list of lines, using cache if available.

        Args:
            file_path: Path to the file to read

        Returns:
            File content as list of lines (without newline characters)

        Raises:
            IOError: If file cannot be read
        """
        if file_path not in self._lines_cache:
            content = self.get_content(file_path)
            self._lines_cache[file_path] = content.split('\n')
        return self._lines_cache[file_path]

    def clear(self):
        """Clear all cached content."""
        self._cache.clear()
        self._lines_cache.clear()

    def has(self, file_path: Path) -> bool:
        """
        Check if file content is cached.

        Args:
            file_path: Path to check

        Returns:
            True if file is cached, False otherwise
        """
        return file_path in self._cache
