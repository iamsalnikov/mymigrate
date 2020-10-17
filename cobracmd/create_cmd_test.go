package cobracmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

type StringValue struct {
	Value string
}

func (s StringValue) String() string {
	return s.Value
}

func (s StringValue) Set(s2 string) error {
	return nil
}

func (s StringValue) Type() string {
	return "string"
}

func TestCreateRunE_CreatesPackageDirectory(t *testing.T) {
	type testCase struct {
		expPath       string
		packagePassed bool
		packageName   string
		pathPassed    bool
		pathName      string
	}

	wd, _ := os.Getwd()

	testCases := map[string]testCase{
		"path and package aren't passed": {
			packagePassed: false,
			packageName:   "",
			pathPassed:    false,
			pathName:      "",
			expPath:       filepath.Join(wd, "migrations"),
		},
		"empty path passed and package isn't passed": {
			packagePassed: false,
			packageName:   "",
			pathPassed:    true,
			pathName:      "",
			expPath:       filepath.Join(wd, "migrations"),
		},
		"filled path passed and package isn't passed": {
			packagePassed: false,
			packageName:   "",
			pathPassed:    true,
			pathName:      "hello",
			expPath:       filepath.Join(wd, "hello", "migrations"),
		},
		"path isn't passed and empty package passed": {
			packagePassed: true,
			packageName:   "",
			pathPassed:    false,
			pathName:      "",
			expPath:       filepath.Join(wd, "migrations"),
		},
		"path isn't passed and filled package passed": {
			packagePassed: true,
			packageName:   "wat",
			pathPassed:    false,
			pathName:      "",
			expPath:       filepath.Join(wd, "wat"),
		},
		"path is passed and filled package is passed": {
			packagePassed: true,
			packageName:   "wat",
			pathPassed:    true,
			pathName:      "hello",
			expPath:       filepath.Join(wd, "hello", "wat"),
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			_ = os.Remove(tc.expPath)

			cmd := &cobra.Command{}
			if tc.packagePassed {
				packageFlag := &pflag.Flag{
					Name:  "package",
					Value: StringValue{Value: tc.packageName},
				}

				cmd.Flags().AddFlag(packageFlag)
			}

			if tc.pathPassed {
				packageFlag := &pflag.Flag{
					Name:  "path",
					Value: StringValue{Value: tc.pathName},
				}

				cmd.Flags().AddFlag(packageFlag)
			}

			err := CreateRunE(cmd, []string{"mig"})
			assert.Nil(t, err)
			assert.DirExists(t, tc.expPath)

			relpath, _ := filepath.Rel(wd, tc.expPath)
			fmt.Println("splitted", relpath)
			if len(relpath) > 0 {
				splittedRelpath := strings.Split(relpath, string(os.PathSeparator))
				if splittedRelpath != nil && len(splittedRelpath) > 0 {
					fmt.Println(splittedRelpath[0])
					_ = os.RemoveAll(splittedRelpath[0])
				}
			}

			_ = os.RemoveAll(tc.expPath)
		})
	}
}

func TestCreateRunE_CreatesMigrationFile(t *testing.T) {
	wd, _ := os.Getwd()

	type testCase struct {
		isPackagePassed  bool
		packageName      string
		migName          string
		dirWithMigration string
	}

	testCases := map[string]testCase{
		"package isn't passed": {
			isPackagePassed:  false,
			packageName:      "",
			migName:          "hello",
			dirWithMigration: filepath.Join(wd, "migrations"),
		},
		"empty package is passed": {
			isPackagePassed:  true,
			packageName:      "",
			migName:          "hello",
			dirWithMigration: filepath.Join(wd, "migrations"),
		},
		"package is passed": {
			isPackagePassed:  true,
			packageName:      "wat",
			migName:          "hello",
			dirWithMigration: filepath.Join(wd, "wat"),
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			cmd := &cobra.Command{}

			if tc.isPackagePassed {
				packageFlag := &pflag.Flag{
					Name:  "package",
					Value: StringValue{Value: tc.packageName},
				}

				cmd.Flags().AddFlag(packageFlag)
			}

			err := CreateRunE(cmd, []string{tc.migName})
			assert.Nil(t, err)

			// Try to find migration file
			fileFound := false
			//migFile := ""
			_ = filepath.Walk(tc.dirWithMigration, func(path string, info os.FileInfo, err error) error {
				fmt.Println(path)
				if info.IsDir() {
					return nil
				}

				if strings.HasSuffix(path, tc.migName+".go") {
					fileFound = true
					//migFile = path
				}

				return nil
			})

			assert.True(t, fileFound, "файл с миграцией найден")
			_ = os.RemoveAll(tc.dirWithMigration)
		})
	}
}

func TestCreateRunE_Out(t *testing.T) {
	out := bytes.NewBufferString("")

	cmd := &cobra.Command{}
	cmd.SetOut(out)
	err := CreateRunE(cmd, []string{"hello"})

	assert.Nil(t, err)

	outStr := out.String()
	assert.Contains(t, outStr, "New migration file is here:")
	assert.Contains(t, outStr, "hello.go")
}
