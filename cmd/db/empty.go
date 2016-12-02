// Copyright © 2016 Jagger Wang <jaggerwang@gmail.com>

package db

import (
	"fmt"

	"github.com/spf13/cobra"

	"zqcserver/models"
)

var emptyFlags struct {
	cluster string
	db      string
	coll    string
}

func init() {
	EmptyCmd.Flags().StringVar(&emptyFlags.cluster, "cluster", "zqc", "which cluster")
	EmptyCmd.Flags().StringVar(&emptyFlags.db, "db", "zqc", "which db")
	EmptyCmd.Flags().StringVar(&emptyFlags.coll, "coll", "", "which collection, empty means all in db")
}

var EmptyCmd = &cobra.Command{
	Use:   "empty",
	Short: "Empty all collections in db",
	Long:  `Empty all collections in db.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := models.EmptyDb(emptyFlags.cluster, emptyFlags.db, emptyFlags.coll)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("empty db ok")
		}
	},
}
