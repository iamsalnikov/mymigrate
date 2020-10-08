package mymigrate

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
)

func Test_MigratorNewNames(t *testing.T) {
	testError := errors.New("test error")

	type testCase struct {
		newNames      map[string]UpFunc
		appliedNames  []string
		applyErr      error
		expectedErr   error
		expectedNames []string
	}

	testCases := []testCase{
		{
			newNames:      map[string]UpFunc{},
			appliedNames:  []string{},
			applyErr:      nil,
			expectedErr:   nil,
			expectedNames: []string{},
		},
		{
			newNames: map[string]UpFunc{
				"0": UpFunc(func(db *sql.DB) error { return nil }),
			},
			appliedNames:  []string{},
			applyErr:      nil,
			expectedErr:   nil,
			expectedNames: []string{"0"},
		},
		{
			newNames: map[string]UpFunc{
				"0": UpFunc(func(db *sql.DB) error { return nil }),
			},
			appliedNames:  []string{"0"},
			applyErr:      nil,
			expectedErr:   nil,
			expectedNames: []string{},
		},
		{
			newNames: map[string]UpFunc{
				"0": UpFunc(func(db *sql.DB) error { return nil }),
				"1": UpFunc(func(db *sql.DB) error { return nil }),
				"2": UpFunc(func(db *sql.DB) error { return nil }),
			},
			appliedNames:  []string{"0", "2"},
			applyErr:      nil,
			expectedErr:   nil,
			expectedNames: []string{"1"},
		},
		{
			newNames: map[string]UpFunc{
				"2": UpFunc(func(db *sql.DB) error { return nil }),
				"1": UpFunc(func(db *sql.DB) error { return nil }),
				"0": UpFunc(func(db *sql.DB) error { return nil }),
			},
			appliedNames:  []string{"0"},
			applyErr:      nil,
			expectedErr:   nil,
			expectedNames: []string{"1", "2"},
		},
		{
			newNames: map[string]UpFunc{
				"0": UpFunc(func(db *sql.DB) error { return nil }),
			},
			appliedNames:  []string{"2"},
			applyErr:      testError,
			expectedErr:   testError,
			expectedNames: []string{},
		},
	}

	for i, c := range testCases {
		resetMigrations()
		resetAppliedFunc()
		resetMarkAppliedFunc()

		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			getApplied = func(db *sql.DB) ([]string, error) {
				return c.appliedNames, c.applyErr
			}

			for name, mig := range c.newNames {
				Add(name, mig)
			}

			result, err := NewNames()
			if err != c.expectedErr {
				t.Errorf("I expected to get error \"%v\" but got \"%v\"", c.expectedErr, err)
				return
			}

			if !isEqualSlices(result, c.expectedNames) {
				t.Errorf("I expected to get names %v but got %v", c.expectedNames, result)
			}
		})
	}

}

func Test_MigratorApply(t *testing.T) {
	applyErr := errors.New("apply err")
	markAppliedErr := errors.New("mark applied err")
	err1 := errors.New("err1")
	err2 := errors.New("err2")

	type testCase struct {
		migrations       map[string]UpFunc
		applyErr         error
		markAppliedErr   error
		expectedErr      error
		expectMarkedCall map[string]bool
	}

	testCases := []testCase{
		{
			migrations:       map[string]UpFunc{},
			applyErr:         nil,
			markAppliedErr:   nil,
			expectedErr:      nil,
			expectMarkedCall: map[string]bool{},
		},
		{
			migrations:       map[string]UpFunc{},
			applyErr:         applyErr,
			markAppliedErr:   nil,
			expectedErr:      applyErr,
			expectMarkedCall: map[string]bool{},
		},
		{
			migrations:       map[string]UpFunc{},
			applyErr:         nil,
			markAppliedErr:   markAppliedErr,
			expectedErr:      nil,
			expectMarkedCall: map[string]bool{},
		},
		{
			migrations: map[string]UpFunc{
				"0": UpFunc(func(db *sql.DB) error { return nil }),
			},
			applyErr:         applyErr,
			markAppliedErr:   markAppliedErr,
			expectedErr:      applyErr,
			expectMarkedCall: map[string]bool{},
		},
		{
			migrations: map[string]UpFunc{
				"0": UpFunc(func(db *sql.DB) error { return nil }),
			},
			applyErr:         nil,
			markAppliedErr:   markAppliedErr,
			expectedErr:      markAppliedErr,
			expectMarkedCall: map[string]bool{"0": true},
		},
		{
			migrations: map[string]UpFunc{
				"0": UpFunc(func(db *sql.DB) error { return err1 }),
				"1": UpFunc(func(db *sql.DB) error { return err2 }),
			},
			applyErr:         nil,
			markAppliedErr:   nil,
			expectedErr:      err1,
			expectMarkedCall: map[string]bool{},
		},
		{
			migrations: map[string]UpFunc{
				"0": UpFunc(func(db *sql.DB) error { return nil }),
				"1": UpFunc(func(db *sql.DB) error { return err2 }),
			},
			applyErr:         nil,
			markAppliedErr:   nil,
			expectedErr:      err2,
			expectMarkedCall: map[string]bool{"0": true},
		},
		{
			migrations: map[string]UpFunc{
				"0": UpFunc(func(db *sql.DB) error { return err1 }),
				"1": UpFunc(func(db *sql.DB) error { return nil }),
			},
			applyErr:         nil,
			markAppliedErr:   nil,
			expectedErr:      err1,
			expectMarkedCall: map[string]bool{},
		},
		{
			migrations: map[string]UpFunc{
				"0": UpFunc(func(db *sql.DB) error { return nil }),
				"1": UpFunc(func(db *sql.DB) error { return nil }),
			},
			applyErr:         nil,
			markAppliedErr:   nil,
			expectedErr:      nil,
			expectMarkedCall: map[string]bool{"0": true, "1": true},
		},
	}

	for i, c := range testCases {
		resetMigrations()
		resetAppliedFunc()
		resetMarkAppliedFunc()

		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			// we've already tested NewNames() function
			// so, here we will return always empty slice
			getApplied = func(db *sql.DB) ([]string, error) {
				return []string{}, c.applyErr
			}

			markedCall := make(map[string]bool)
			markApplied = func(db *sql.DB, name string) error {
				if !c.expectMarkedCall[name] {
					t.Errorf("I didn't excpect that mig '%s' will be marked as aplied", name)
				}

				markedCall[name] = true

				return c.markAppliedErr
			}

			for name, mig := range c.migrations {
				Add(name, mig)
			}

			err := Apply()
			if err != c.expectedErr {
				t.Errorf("I expected to get error \"%v\" but got \"%v\"", c.expectedErr, err)
			}

			if len(markedCall) != len(c.expectMarkedCall) {
				t.Errorf("I expected that these migrations (%+v) we will try to mark as applied. "+
					"But we tried to mark as applied %+v", c.expectMarkedCall, markedCall)
			}
		})
	}
}

func isEqualSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for j, n := range a {
		if n != b[j] {
			return false
		}
	}

	return true
}
