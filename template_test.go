package mymigrate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var template = `package %s

import (
	"database/sql"

	"github.com/iamsalnikov/mymigrate"
)

func init() {
	mymigrate.Add(
		"%s",
		func(db *sql.DB) error {
			// TODO: write UP logic
			return nil
		},
		func(db *sql.DB) error {
			// TODO: write down logic

			return nil
		},
	)
}
`

func TestTemplate(t *testing.T) {
	type testCase struct {
		pkg    string
		name   string
		expPkg string
	}

	testCases := map[string]testCase{
		"empty package name": {
			pkg:    "",
			name:   "mig-001",
			expPkg: "migrations",
		},
		"specified package name": {
			pkg:    "hello",
			name:   "mig-001",
			expPkg: "hello",
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			tpl, name := Template(tc.pkg, tc.name)
			exp := fmt.Sprintf(template, tc.expPkg, name)

			assert.EqualValues(t, exp, tpl)
			assert.Contains(t, tpl, name)
		})
	}

}
