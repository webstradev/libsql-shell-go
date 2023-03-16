package cmd

import (
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"

	"github.com/libsql/libsql-shell-go/pkg/libsql"
	"github.com/libsql/libsql-shell-go/shell"
)

type RootArgs struct {
	statements string
	quiet      bool
}

func NewRootCmd() *cobra.Command {
	var rootArgs RootArgs = RootArgs{}
	var rootCmd = &cobra.Command{
		SilenceUsage: true,
		Use:          "libsql-shell",
		Short:        "A cli for executing SQL statements on a libSQL or SQLite database",
		Args:         cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			dbPath := args[0]
			db, err := libsql.NewDb(dbPath)
			if err != nil {
				return err
			}
			defer db.Close()

			if cmd.Flag("exec").Changed {
				if len(rootArgs.statements) == 0 {
					return fmt.Errorf("no SQL command to execute")
				}

				result, err := db.ExecuteStatements(rootArgs.statements)
				if err != nil {
					return err
				}

				err = libsql.PrintStatementsResult(result, cmd.OutOrStdout(), false)
				if err != nil {
					return err
				}

				return nil
			}

			shellConfig := shell.ShellConfig{
				InF:         cmd.InOrStdin(),
				OutF:        cmd.OutOrStdout(),
				ErrF:        cmd.ErrOrStderr(),
				HistoryMode: shell.PerDatabaseHistory,
				HistoryName: "libsql",
				QuietMode:   rootArgs.quiet,
			}

			return shell.RunShell(db, shellConfig)
		},
	}

	rootCmd.Flags().StringVarP(&rootArgs.statements, "exec", "e", "", "SQL statements separated by ;")
	rootCmd.Flags().BoolVarP(&rootArgs.quiet, "quiet", "q", false, "Don't print welcome message")

	return rootCmd
}

func Init() {
	var rootCmd *cobra.Command = NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}