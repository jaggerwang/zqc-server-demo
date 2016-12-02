package model

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"

	"zqcserver/models"
)

func TestMain(m *testing.M) {
	addrs, ctn, err := startMongo(10, 1*time.Second)
	if err != nil {
		panic(err)
	}
	viper.Set("mongodb.zqc", map[string]interface{}{
		"addrs": addrs,
	})

	err = models.CreateDbIndexes("zqc", "zqc", "", -1)

	result := m.Run()

	removeMongo(ctn)

	os.Exit(result)
}
