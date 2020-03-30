package goose

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"
)

const VERSION = "v2.7.0-rc3"

var (
	duplicateCheckOnce sync.Once
	minVersion         = int64(0)
	maxVersion         = int64((1 << 63) - 1)
	timestampFormat    = "20060102150405"
	verbose            = false
)

// SetVerbose set the goose verbosity mode
func SetVerbose(v bool) {
	verbose = v
}

var productRange = []int{0, 200000}
var projectRange = []int{200000, 300000}
var unknownRange = []int{0, int(maxVersion)}

// Run runs a goose command.
func Run(command, dbLevel string, db *sql.DB, dir string, args ...string) error {
	log.Printf("goose command=%s dbLevel=%s dir=%s args=%v ", command, dbLevel, dir, args)
	r := unknownRange
	if dbLevel == "product" {
		r = productRange
	} else if dbLevel == "project" {
		r = projectRange
	}
	switch command {
	case "up":
		if err := Up(db, dir, r); err != nil {
			return err
		}
	case "up-by-one":
		if err := UpByOne(db, dir, r); err != nil {
			return err
		}
	case "up-to":
		if len(args) == 0 {
			return fmt.Errorf("up-to must be of form: goose [OPTIONS] DRIVER DBSTRING up-to VERSION")
		}
		version, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("version must be a number (got '%s')", args[0])
		}
		if err := UpTo(db, dir, r, version); err != nil {
			return err
		}
	case "create":
		if len(args) == 0 {
			return fmt.Errorf("create must be of form: goose [OPTIONS] DRIVER DBSTRING create NAME [go|sql]")
		}

		if err := Create(db, dir, args[0], "sql"); err != nil {
			return err
		}
	case "down":
		if err := Down(db, dir, r); err != nil {
			return err
		}
	case "down-to":
		if len(args) == 0 {
			return fmt.Errorf("down-to must be of form: goose [OPTIONS] DRIVER DBSTRING down-to VERSION")
		}

		version, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("version must be a number (got '%s')", args[0])
		}
		if err := DownTo(db, dir, r, version); err != nil {
			return err
		}

		if err := Fix(dir); err != nil {
			return err
		}
	case "reset":
		if err := Reset(db, dir, r); err != nil {
			return err
		}
	case "status":
		if err := Status(db, dir); err != nil {
			return err
		}
	case "version":
		if err := Version(db, dir); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%q: no such command", command)
	}
	return nil
}
