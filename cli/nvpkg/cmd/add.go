package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"

	novuspack "github.com/novus-engine/novuspack/api/go"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <package path> <file or dir> [file or dir ...]",
	Short: "Add files or directories to a NovusPack package",
	Args:  cobra.MinimumNArgs(2),
	RunE:  runAdd,
}

var (
	addStoredPath          string
	addBasePath            string
	addPreserveDepth       int
	addFlatten             bool
	addNoFollowSymlinks    bool
	addPreservePermissions bool
	addPreserveOwnership   bool
)

func init() {
	addCmd.Flags().StringVar(&addStoredPath, "as", "", "Store under this path (single file only)")
	addCmd.Flags().StringVar(&addBasePath, "base-path", "", "Strip this prefix from source paths (at most one of --as, --base-path, --preserve-depth, --flatten)")
	addCmd.Flags().IntVar(&addPreserveDepth, "preserve-depth", 0, "Keep N directory levels from source (0=off)")
	addCmd.Flags().BoolVar(&addFlatten, "flatten", false, "Store all files at package root")
	addCmd.Flags().BoolVar(&addNoFollowSymlinks, "no-follow-symlinks", false, "Do not follow symlinks; reject them")
	addCmd.Flags().BoolVar(&addPreservePermissions, "preserve-permissions", false, "Store Unix permission bits")
	addCmd.Flags().BoolVar(&addPreserveOwnership, "preserve-ownership", false, "Store UID/GID (implies --preserve-permissions)")
}

func runAdd(_ *cobra.Command, args []string) error {
	ctx := context.Background()
	pkgPath := args[0]
	sources := args[1:]

	pkg, err := openOrCreatePackage(ctx, pkgPath)
	if err != nil {
		return err
	}
	defer func() { _ = pkg.Close() }()

	opts, err := buildAddFileOptions(nil)
	if err != nil {
		return err
	}
	if err := addSources(ctx, pkg, sources, opts); err != nil {
		return err
	}
	if err := pkg.Write(ctx); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Added %d item(s) to %s\n", len(sources), pkgPath)
	return nil
}

func openOrCreatePackage(ctx context.Context, pkgPath string) (novuspack.Package, error) {
	if _, statErr := os.Stat(pkgPath); statErr != nil && os.IsNotExist(statErr) {
		pkg, err := novuspack.NewPackage()
		if err != nil {
			return nil, fmt.Errorf("new package: %w", err)
		}
		if err := pkg.Create(ctx, pkgPath); err != nil {
			_ = pkg.Close()
			return nil, fmt.Errorf("create: %w", err)
		}
		return pkg, nil
	}
	return openPackage(ctx, pkgPath, false)
}

// openPackage opens an existing package, optionally read-only. Used by list, read, info, extract, validate.
func openPackage(ctx context.Context, path string, readOnly bool) (novuspack.Package, error) {
	if readOnly {
		return novuspack.OpenPackageReadOnly(ctx, path)
	}
	return novuspack.OpenPackage(ctx, path)
}

const (
	flagTrue = "true"
	flagYes  = "yes"
)

// flagBool returns true if v is a truthy flag value (e.g. "1", "true", "yes").
func flagBool(v string) bool {
	return v == "1" || v == flagTrue || v == flagYes
}

// buildAddFileOptions builds AddFileOptions from package vars (when flags is nil) or from interactive flags map.
// At most one of StoredPath, BasePath, PreserveDepth, FlattenPaths may be set.
func buildAddFileOptions(flags map[string]string) (*novuspack.AddFileOptions, error) {
	opts := &novuspack.AddFileOptions{}
	var pathOptCount int
	var err error
	if flags != nil {
		pathOptCount, err = applyAddFileOptionsFromFlags(flags, opts)
		if err != nil {
			return nil, err
		}
	} else {
		pathOptCount = applyAddFileOptionsFromVars(opts)
	}
	if pathOptCount > 1 {
		return nil, fmt.Errorf("at most one of --as, --base-path, --preserve-depth, --flatten may be set")
	}
	return opts, nil
}

func applyAddFileOptionsFromFlags(flags map[string]string, opts *novuspack.AddFileOptions) (int, error) {
	var n int
	if v := flags["as"]; v != "" {
		opts.StoredPath.Set(v)
		n++
	}
	if v := flags["base-path"]; v != "" {
		opts.BasePath.Set(v)
		n++
	}
	if v := flags["preserve-depth"]; v != "" {
		d, err := strconv.Atoi(v)
		if err != nil || d < 0 {
			return 0, fmt.Errorf("invalid --preserve-depth: %s", v)
		}
		opts.PreserveDepth.Set(d)
		n++
	}
	if flagBool(flags["flatten"]) {
		opts.FlattenPaths.Set(true)
		n++
	}
	if flagBool(flags["no-follow-symlinks"]) {
		opts.FollowSymlinks.Set(false)
	}
	if flagBool(flags["preserve-permissions"]) {
		opts.PreservePermissions.Set(true)
	}
	if flagBool(flags["preserve-ownership"]) {
		opts.PreserveOwnership.Set(true)
	}
	return n, nil
}

func applyAddFileOptionsFromVars(opts *novuspack.AddFileOptions) int {
	var n int
	if addStoredPath != "" {
		opts.StoredPath.Set(addStoredPath)
		n++
	}
	if addBasePath != "" {
		opts.BasePath.Set(addBasePath)
		n++
	}
	if addPreserveDepth > 0 {
		opts.PreserveDepth.Set(addPreserveDepth)
		n++
	}
	if addFlatten {
		opts.FlattenPaths.Set(true)
		n++
	}
	if addNoFollowSymlinks {
		opts.FollowSymlinks.Set(false)
	}
	if addPreservePermissions {
		opts.PreservePermissions.Set(true)
	}
	if addPreserveOwnership {
		opts.PreserveOwnership.Set(true)
	}
	return n
}

func addSources(ctx context.Context, pkg novuspack.Package, sources []string, opts *novuspack.AddFileOptions) error {
	for _, src := range sources {
		info, err := os.Stat(src)
		if err != nil {
			return fmt.Errorf("stat %s: %w", src, err)
		}
		if info.IsDir() {
			_, err = pkg.AddDirectory(ctx, src, opts)
			if err != nil {
				return fmt.Errorf("add directory %s: %w", src, err)
			}
		} else {
			_, err = pkg.AddFile(ctx, src, opts)
			if err != nil {
				return fmt.Errorf("add file %s: %w", src, err)
			}
		}
	}
	return nil
}
