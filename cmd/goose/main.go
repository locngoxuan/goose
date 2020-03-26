package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"scm.wcs.fortna.com/lngo/goose"
)

var (
	flags   = flag.NewFlagSet("goose", flag.ExitOnError)
	dir     = flags.String("dir", ".", "directory with migration files")
	verbose = flags.Bool("v", false, "enable verbose mode")
	help    = flags.Bool("h", false, "print help")
	version = flags.Bool("version", false, "print version")
	dbLevel = flags.String("db-level", "", "level of database, i.e. product, project")
)

func main() {
	flags.Usage = usage
	err := flags.Parse(os.Args[1:])
	if err != nil {
		flags.Usage()
		return
	}

	if *version {
		fmt.Println(goose.VERSION)
		return
	}
	if *verbose {
		goose.SetVerbose(true)
	}

	args := flags.Args()
	if len(args) == 0 || *help {
		flags.Usage()
		return
	}

	switch args[0] {
	case "create":
		if err := goose.Run("create", "", nil, *dir, args[1:]...); err != nil {
			log.Fatalf("goose run: %v", err)
		}
		return
	case "fix":
		if err := goose.Run("fix", "", nil, *dir); err != nil {
			log.Fatalf("goose run: %v", err)
		}
		return
	}

	if len(args) < 3 {
		flags.Usage()
		return
	}

	driver, dbstring, command := args[0], args[1], args[2]

	db, err := goose.OpenDBWithDriver(driver, dbstring)
	if err != nil {
		log.Fatalf("-dbstring=%q: %v\n", dbstring, err)
	}

	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}
	if err := goose.Run(command, *dbLevel, db, *dir, arguments...); err != nil {
		log.Fatalf("goose run: %v", err)
	}
}

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `Usage: goose [OPTIONS] DRIVER DBSTRING COMMAND

Drivers:
    postgres
    oracle

Examples:
    goose create add_new_table sql

    goose postgres "postgres://username:password@host[:port]/table?sslmode=disable" status
    goose oracle "username/password@[//]host[:port][/service_name][:server][/instance_name]" status

Options:
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations
`
)
