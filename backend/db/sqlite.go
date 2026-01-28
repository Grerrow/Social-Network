package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

var Database *sql.DB

func InitDB() error {
	// ------------------------------------------------------------------------------------------------------
	// CHECKING if the "./db" directory exists, if not we create it

	_, statErr := os.Stat("./db")
	// os.Stat returns info about a file/folder or an error if it doesnt exist
	// we only need the error here, so we can use it in next check

	if os.IsNotExist(statErr) {
		// os.IsNotExist only takes as arguments error values (errors from os.Stat, os.Open etc)
		// it returns True if the error from os.Stat indicates that the file/directory doesn't exist or
		// False for any other error (the file/folder exists or permission issues etc)

		makeDirErr := os.MkdirAll("./db", 0o755)
		if makeDirErr != nil {
			return fmt.Errorf("failed to create database directory: %v", makeDirErr)
		}
	}

	// ------------------------------------------------------------------------------------------------------
	// OPENING DATABASE CONNECTION and storing it in the variable "database" we created in start

	var openErr error

	Database, openErr = sql.Open("sqlite3", "./db/social-network.db")
	// sql.Open returns: *sql.DB, error
	if openErr != nil {
		return fmt.Errorf("failed to open database: %v", openErr)
	}

	pingErr := Database.Ping()
	// Ping() is a method of the sql.DB struct (func (db *sql.DB) Ping() error)
	if pingErr != nil {
		return fmt.Errorf("failed to connect to database: %v", pingErr)
	}

	_, foreignKeyErr := Database.Exec("PRAGMA foreign_keys = ON;")
	if foreignKeyErr != nil {
		return fmt.Errorf("failed to enable foreign keys: %v", foreignKeyErr)
	}

	// ------------------------------------------------------------------------------------------------------
	// APPLYING MIGRATIONS
	migrationErr := applyMigrations()
	if migrationErr != nil {
		return fmt.Errorf("failed to apply migrations: %v", migrationErr)
	}

	return nil
}

func CloseDB() error {
	if Database != nil {
		return Database.Close()
	}
	return nil
}

func applyMigrations() error {
	// 1. DRIVER: create driver instance (of database.Driver type), driver: its how the DB talks to the migrator) -----
	driver, err := sqlite3.WithInstance(Database, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %v", err)
	}

	// 2. MIGRATE OBJECT: create a migrate instance (of *migrate.Migrate type) ----------------------------------------
	migration, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations/sqlite",
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %v", err)
	}

	// 3. APPLY MIGRATIONS: apply all pending migrations --------------------------------------------------------------
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		// If there are no new migrations to apply, ErrNoChange is returned
		return fmt.Errorf("could not apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}

/* ==================================================================================================================
notes:
	with the _ in start, we only run the package's init function
	we do NOT make its exported functions, types, and variables available to use directly
	we do not have to use the package elsewhere in the project


	Exec() - method of sql.DB struct
	executes any SQL statement that does NOT return rows, so it works for
	enabling foreign keys and creating tables
	(by default, foreign key constraints are disabled in SQLite)


	Down migrations are executed manually (not auto-applied) - use CLI tool or separate script for rollback if needed
*/
