package cmd

import (
	"context"
	"fmt"
	"os"

	novuspack "github.com/novus-engine/novuspack/api/go"
	"github.com/spf13/cobra"
)

var commentCmd = &cobra.Command{
	Use:   "comment <package path> [--set \"...\" | --clear]",
	Short: "Get or set package comment",
	Long:  "With no flags, prints the package comment. Use --set to set, --clear to remove. Changes are written to disk.",
	Args:  cobra.ExactArgs(1),
	RunE:  runComment,
}

var (
	commentSet   string
	commentClear bool
)

func init() {
	commentCmd.Flags().StringVar(&commentSet, "set", "", "Set comment to this string")
	commentCmd.Flags().BoolVar(&commentClear, "clear", false, "Clear the comment")
}

func runComment(_ *cobra.Command, args []string) error {
	path := args[0]
	ctx := context.Background()

	if commentSet != "" && commentClear {
		return fmt.Errorf("cannot use both --set and --clear")
	}

	pkg, err := novuspack.OpenPackage(ctx, path)
	if err != nil {
		return fmt.Errorf("open package: %w", err)
	}
	defer func() { _ = pkg.Close() }()

	if commentClear {
		if err := pkg.ClearComment(); err != nil {
			return fmt.Errorf("clear comment: %w", err)
		}
		if err := pkg.Write(ctx); err != nil {
			return fmt.Errorf("write: %w", err)
		}
		_, _ = fmt.Fprintf(os.Stdout, "Comment cleared\n")
		return nil
	}
	if commentSet != "" {
		if err := pkg.SetComment(commentSet); err != nil {
			return fmt.Errorf("set comment: %w", err)
		}
		if err := pkg.Write(ctx); err != nil {
			return fmt.Errorf("write: %w", err)
		}
		_, _ = fmt.Fprintf(os.Stdout, "Comment set\n")
		return nil
	}
	comment := pkg.GetComment()
	if comment == "" {
		_, _ = fmt.Fprintln(os.Stdout, "(no comment)")
	} else {
		_, _ = fmt.Fprintf(os.Stdout, "%s\n", comment)
	}
	return nil
}
