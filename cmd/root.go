// Copyright © 2016 Jagger Wang <jaggerwang@gmail.com>

package cmd

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2/bson"

	"zqcserver/services"
)

var rootFlags struct {
	cfgFile string
}

var rootCmd = &cobra.Command{
	Use:   "zqc",
	Short: "Zai qiu chang app",
	Long:  `Zai qiu chang app.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLog)

	rootCmd.PersistentFlags().StringVarP(&rootFlags.cfgFile, "config", "c", "./config.json", "config file")
	rootCmd.PersistentFlags().StringP("dir.data", "d", "", "runtime data directory")
	rootCmd.PersistentFlags().Bool("server.debug", false, "enable/disable debug mode")
	rootCmd.PersistentFlags().String("log.level", "", "log level")
	rootCmd.PersistentFlags().String("mongodb.zqc.addrs", "", "addrs of zai qiu chang mongodb")

	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.BindPFlags(rootCmd.Flags())

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(dbCmd)

	registerGobTypes()
}

func initConfig() {
	if e := os.Getenv("ZQC_CONFIG_FILE"); e != "" {
		rootFlags.cfgFile = e
	}
	viper.SetConfigFile(rootFlags.cfgFile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println("using config file", viper.ConfigFileUsed())

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config changed")
	})
}

func initLog() {
	log.SetFormatter(&log.JSONFormatter{})

	level, err := log.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		panic(err)
	}
	log.SetLevel(level)

	w, err := os.OpenFile(filepath.Join(viper.GetString("dir.data"), viper.GetString("log.server.file")), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		panic(err)
	}
	log.SetOutput(w)
}

func registerGobTypes() {
	gob.Register(bson.NewObjectId())
	gob.Register(services.VerifyCode{})
	gob.Register(services.User{})
}
