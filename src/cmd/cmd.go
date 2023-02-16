package cmd

import (
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"

	"github.com/chiselstrike/libsql-shell/src/lib"
)

type RootArgs struct {
	statements string
	dbPath     string
}

func NewRootCmd() *cobra.Command {
	var rootArgs RootArgs = RootArgs{}
	var rootCmd = &cobra.Command{
		Use:          "libsql-shell",
		Short:        "A cli for executing SQL statements on a libSQL or SQLite database",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := lib.OpenOrCreateSQLite3(rootArgs.dbPath)
			if err != nil {
				return err
			}
			defer db.Close()

			result, err := lib.ExecuteStatements(db, rootArgs.statements)
			if err != nil {
				return err
			}
			cmd.Println(result)
			return nil
		},
	}

	rootCmd.Flags().StringVarP(&rootArgs.statements, "exec", "e", "", "SQL statements separated by ;")
	rootCmd.MarkFlagRequired("exec")
	rootCmd.Flags().StringVarP(&rootArgs.dbPath, "db", "d", "", "Path to database file")
	rootCmd.MarkFlagRequired("db")

	return rootCmd
}

func Init() {
	var rootCmd *cobra.Command = NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
