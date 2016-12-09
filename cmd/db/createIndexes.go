// Copyright © 2016 Jagger Wang <jaggerwang@gmail.com>

package db

import (
	"fmt"

	"github.com/spf13/cobra"

	"zqc/models"
)

var createIndexesFlags struct {
	cluster string
	db      string
	coll    string
	pos     int
}

func init() {
	CreateIndexesCmd.Flags().StringVar(&createIndexesFlags.cluster, "cluster", "zqc", "which cluster")
	CreateIndexesCmd.Flags().StringVar(&createIndexesFlags.db, "db", "zqc", "which db")
	CreateIndexesCmd.Flags().StringVar(&createIndexesFlags.coll, "coll", "", "which collection, empty means all collections in db")
	CreateIndexesCmd.Flags().IntVar(&createIndexesFlags.pos, "pos", -1, "which index, postion in index array, -1 means all")
}

var CreateIndexesCmd = &cobra.Command{
	Use:   "createIndexes",
	Short: "Create indexes",
	Long:  `Create indexes.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := models.CreateDbIndexes(createIndexesFlags.cluster, createIndexesFlags.db, createIndexesFlags.coll, createIndexesFlags.pos)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("create indexes ok")
		}
	},
}
