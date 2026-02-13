package migration

import (
	"fmt"
	"io/fs"
)

// TestMigrationFiles prints the name of every file in MigrationFS for debugging.
// It does not execute any migration.
func TestMigrationFiles() error {
	return fs.WalkDir(MigrationFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fmt.Println(path)
		}
		return nil
	})
}