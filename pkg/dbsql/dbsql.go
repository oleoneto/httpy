package dbsql

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mattn/go-sqlite3"
)

func UseSQLite(dbname string) (*sql.DB, error) {
	if dbname == "" {
		return nil, fmt.Errorf("no database name provided")
	}

	var regex = func(re, s string) (bool, error) { return regexp.MatchString(re, s) }

	sql.Register("sqlite3_ext",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				return conn.RegisterFunc("regexp", regex, true)
			},
		},
	)

	options := []string{"fts5=on", "_cslike=off", "_fk=on", "_ignore_check_constraints=off", "_journal=WAL"}
	d, err := sql.Open("sqlite3_ext", fmt.Sprintf("%s?%s", dbname, strings.Join(options, "&")))
	if err != nil {
		return nil, err
	}

	d.Exec(`PRAGMA busy_timeout = 5000`)

	return d, nil
}
