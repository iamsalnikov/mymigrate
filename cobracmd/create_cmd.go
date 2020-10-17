package cobracmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/iamsalnikov/mymigrate"

	"github.com/spf13/cobra"
)

// CreateCmd is a cobra command to create new migration file
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new migration file",
}

// CreateRunE is a cobra run function to create new migration file
func CreateRunE(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("please, pass migration name as an argument")
	}

	packageName := "migrations"

	packageFlag := cmd.Flag("package")
	if packageFlag != nil && packageFlag.Value != nil && len(packageFlag.Value.String()) > 0 {
		packageName = packageFlag.Value.String()
	}

	basePath, err := os.Getwd()
	if err != nil {
		return err
	}

	path := ""
	pathFlag := cmd.Flag("path")
	if pathFlag != nil && pathFlag.Value != nil && len(pathFlag.Value.String()) > 0 {
		path = pathFlag.Value.String()
	}

	if !filepath.IsAbs(path) {
		path = filepath.Join(basePath, path)
	}

	dirpath := filepath.Join(path, packageName)
	err = os.MkdirAll(dirpath, 0766)
	if err != nil {
		return err
	}

	template, filename := mymigrate.Template(packageName, args[0])
	migFilePath := filepath.Join(dirpath, filename+".go")
	f, err := os.Create(migFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(template)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "New migration file is here: %s\n", migFilePath)

	return nil
}
