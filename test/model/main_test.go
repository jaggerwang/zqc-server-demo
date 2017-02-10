package model

import (
	"os"
	"testing"

	"zqc/test"
)

func TestMain(m *testing.M) {
	test.EmptyDb()

	result := m.Run()

	os.Exit(result)
}
