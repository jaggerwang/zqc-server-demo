// Copyright © 2016 Jagger Wang <jaggerwang@gmail.com>

package db

import (
	"fmt"

	"github.com/kr/pretty"
	"github.com/spf13/cobra"

	"jaggerwang.net/zqcserverdemo/models"
)

var listIndexesFlags struct {
	cluster  string
	db       string
	coll     string
	required bool
}

func init() {
	ListIndexesCmd.Flags().StringVar(&listIndexesFlags.cluster, "cluster", "zqc", "which cluster")
	ListIndexesCmd.Flags().StringVar(&listIndexesFlags.db, "db", "zqc", "which db")
	ListIndexesCmd.Flags().StringVar(&listIndexesFlags.coll, "coll", "", "which collection, empty means all in db")
	ListIndexesCmd.Flags().BoolVar(&listIndexesFlags.required, "required", false, "if true list required indexes, else list exist indexes")
}

var ListIndexesCmd = &cobra.Command{
	Use:   "listIndexes",
	Short: "List indexes",
	Long:  `List indexes.`,
	Run: func(cmd *cobra.Command, args []string) {
		if listIndexesFlags.required {
			if listIndexesFlags.db == "zqc" {
				fmt.Printf("%# v", pretty.Formatter(models.ZqcDbIndexes))
			} else {
				fmt.Println("unknown db", listIndexesFlags.db)
			}
		} else {
			collNames, err := models.DbCollNames(listIndexesFlags.cluster, listIndexesFlags.db)
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, collName := range collNames {
				if listIndexesFlags.coll == "" || listIndexesFlags.coll == collName {
					coll, err := models.NewMongoColl(listIndexesFlags.cluster, listIndexesFlags.db, collName)
					if err != nil {
						fmt.Println(err)
						return
					}

					indexes, err := coll.Indexes()
					if err != nil {
						fmt.Println(err)
						return
					}
					fmt.Printf("%# v\n", pretty.Formatter(indexes))
				}
			}
		}
	},
}
