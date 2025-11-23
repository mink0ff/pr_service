package integration

import (
	"os"
	"testing"

	"github.com/mink0ff/pr_service/tests/utils"
)

var ts *utils.TestServices

func TestMain(m *testing.M) {
	ts = utils.InitTestServices()
	defer ts.Teardown()

	utils.TruncateTables(ts.DB)

	code := m.Run()
	os.Exit(code)
}
