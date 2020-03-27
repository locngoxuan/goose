package goose

import (
	"database/sql"
	"fmt"
)

// SQLDialect abstracts the details of specific SQL dialects
// for goose's few SQL specific statements
type SQLDialect interface {
	preCreateVersionTableSQL() string
	createVersionTableSQL() string                       // sql string to create the db version table
	insertVersionSQL(version int, isApplied bool) string // sql string to insert the initial version table row
	deleteVersionSQL(version int) string                 // sql string to delete version
	migrationSQL(version int) string                     // sql string to retrieve migrations
	dbVersionQuery(inRange []int, db *sql.DB) (*sql.Rows, error)
	booleanValue(value bool) interface{}
}

var dialect SQLDialect = &PostgresDialect{}

var (
	postgresTable = `CREATE TABLE %s (
            	id serial NOT NULL,
                version_id bigint NOT NULL,
                is_applied boolean NOT NULL,
                tstamp timestamp NULL default now(),
                PRIMARY KEY(id)
            );`

	preCreateTableOracle = `CREATE SEQUENCE goose_seq START WITH 1 increment by 1`

	oracleTable = `CREATE TABLE %s (
            	id NUMBER(19) DEFAULT goose_seq.nextval NOT NULL,
                version_id NUMBER(19,0),
                is_applied NUMBER(1) DEFAULT 0 NOT NULL,
                tstamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                PRIMARY KEY(id)
            )`
)

// GetDialect gets the SQLDialect
func GetDialect() SQLDialect {
	return dialect
}

// SetDialect sets the SQLDialect
func SetDialect(d string) error {
	switch d {
	case "postgres":
		dialect = &PostgresDialect{}
	case "oracle":
		dialect = &OracleDialect{}
	default:
		return fmt.Errorf("%q: unknown dialect", d)
	}

	return nil
}

////////////////////////////
// Postgres
////////////////////////////

// PostgresDialect struct.
type PostgresDialect struct{}

func (pg PostgresDialect) preCreateVersionTableSQL() string {
	return ""
}

func (pg PostgresDialect) createVersionTableSQL() string {
	return fmt.Sprintf(postgresTable, TableName())
}

func (pg PostgresDialect) insertVersionSQL(version int, isApplied bool) string {
	return fmt.Sprintf("INSERT INTO %s (version_id, is_applied) VALUES (%d, %v);", TableName(), version, isApplied)
}

func (pg PostgresDialect) dbVersionQuery(inRange []int, db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query(fmt.Sprintf(`SELECT version_id, is_applied from %s 
												where version_id >= %d and version_id < %d ORDER BY id DESC`,
		TableName(), inRange[0], inRange[1]))
	if err != nil {
		return nil, err
	}

	return rows, err
}

func (pg PostgresDialect) booleanValue(value bool) interface{} {
	return value
}

func (m PostgresDialect) migrationSQL(version int) string {
	return fmt.Sprintf("SELECT tstamp, is_applied FROM %s WHERE version_id=%d ORDER BY tstamp DESC LIMIT 1", TableName(), version)
}

func (pg PostgresDialect) deleteVersionSQL(version int) string {
	return fmt.Sprintf("DELETE FROM %s WHERE version_id=%d;", TableName(), version)
}

////////////////////////////
// Postgres
////////////////////////////

// PostgresDialect struct.
type OracleDialect struct{}

func (pg OracleDialect) preCreateVersionTableSQL() string {
	return preCreateTableOracle
}

func (pg OracleDialect) createVersionTableSQL() string {
	return fmt.Sprintf(oracleTable, TableName())
}

func (pg OracleDialect) insertVersionSQL(version int, isApplied bool) string {
	v := pg.booleanValue(isApplied).(int)
	return fmt.Sprintf("INSERT INTO %s (version_id, is_applied) VALUES (%d, %d)", TableName(), version, v)
}

func (pg OracleDialect) booleanValue(value bool) interface{} {
	if value {
		return 1
	}
	return 0
}

func (pg OracleDialect) dbVersionQuery(inRange []int, db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query(fmt.Sprintf(`SELECT version_id, is_applied from %s 
												where version_id >= %d and version_id < %d ORDER BY id DESC`,
		TableName(), inRange[0], inRange[1]))
	if err != nil {
		return nil, err
	}

	return rows, err
}

func (m OracleDialect) migrationSQL(version int) string {
	return fmt.Sprintf("SELECT tstamp, is_applied FROM %s WHERE version_id=%d ORDER BY tstamp DESC LIMIT 1", TableName(), version)
}

func (pg OracleDialect) deleteVersionSQL(version int) string {
	return fmt.Sprintf("DELETE FROM %s WHERE version_id=%d", TableName(), version)
}
