// Copyright © 2016 Jagger Wang <jaggerwang@gmail.com>

package cmd

import (
	"github.com/spf13/cobra"

	"jaggerwang.net/zqcserverdemo/cmd/db"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database admin",
	Long:  `Database admin.`,
}

func init() {
	dbCmd.AddCommand(db.CreateIndexesCmd)
	dbCmd.AddCommand(db.ListIndexesCmd)
	dbCmd.AddCommand(db.EmptyCmd)
}
