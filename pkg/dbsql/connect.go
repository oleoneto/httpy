package dbsql

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

func ConnectDatabase(options DBConnectOptions) (SqlBackend, error) {
	var db *sql.DB
	var err error

	if db, err = UseSQLite(options.Filename); err != nil {
		return db, err
	}

	if db == nil {
		return db, fmt.Errorf("database connection failed")
	}

	if options.VerboseLogging {
		loggerAdapter := zerologadapter.New(zerolog.New(os.Stdout))
		db = sqldblogger.OpenDriver(options.Filename, db.Driver(), loggerAdapter)

		if db == nil {
			return db, fmt.Errorf("database logger failed")
		}
	}

	db.Exec("PRAGMA journal_mode=WAL")

	return db, nil
}
