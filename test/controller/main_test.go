package controller

import (
	"os"
	"testing"

	"zqc/test"
)

func TestMain(m *testing.M) {
	test.emptyDb()

	result := m.Run()

	os.Exit(result)
}
