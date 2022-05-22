package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type pg struct {
	logger   *zap.SugaredLogger
	username string
	password string
	host     string
	port     string
	dbName   string
	sslMode  SslMode
}

type SslMode int64

const (
	UNDEFINED SslMode = iota
	Enable
	Disable
)

// SQL query statements
const sqlCheckTableExists string = "SELECT 1 FROM information_schema.tables WHERE table_schema = 'chronicler' AND" +
	" table_name = '%v';"
const sqlInsertEvent string = `INSERT INTO chronicler.%v (
	id,
	streamname,
	version,
	event)
VALUES (
	'%v',
	'%v',
	'%v',
	'%v')
RETURNING position;`

// SQL exec statements
const sqlCreateTable string = `CREATE TABLE chronicler.%v (
	position SERIAL PRIMARY KEY,
	id TEXT NOT NULL,
	time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	streamname TEXT NOT NULL,
	version TEXT NOT NULL,
	event JSONB NOT NULL
);`

func (mode SslMode) String() string {
	return [...]string{"UNDEFINED", "enable", "disable"}[mode]
}

func Init(logger *zap.SugaredLogger, username string, password string, host string, port string, dbName string, sslMode SslMode) *pg {
	retval := new(pg)
	retval.logger = logger
	retval.username = username
	retval.password = password
	retval.host = host
	retval.port = port
	retval.dbName = dbName
	retval.sslMode = sslMode
	if retval.sslMode == UNDEFINED {
		retval.sslMode = Disable
	}

	return retval
}

// Writes an event to a given stream.  streamName is assumed to be a valid stream name
// (i.e. of the form noun:verb-uuid4)
func (pg pg) WriteEvent(streamName string, event string, uuidString string, version string) (int64, error) {
	connectionString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v",
		pg.username, pg.password, pg.host, pg.port, pg.dbName, pg.sslMode.String())
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return -1, err
	}
	defer db.Close()

	// first, get the table name from the stream name
	tableName, err := pg.getTableNameFromStreamName(streamName)
	if err != nil {
		return -1, err
	}

	// second, make sure the table exists
	result, err := pg.checkTableExists(db, tableName)
	if err != nil {
		return -1, err
	}
	pg.logger.Infof("checkTableExists = %v\n", result)

	// if the table doesn't exist, we need to create the database table
	if !result {
		if err := pg.createTableForStream(db, tableName); err != nil {
			return -1, err
		}
		pg.logger.Infof("Created table [%v]", tableName)
	}

	// finally, write the event to the table
	index, err := pg.writeEventToTable(db, tableName, streamName, event, uuidString, version)
	if err != nil {
		return -1, err
	}

	// return the index (i.e. the primary key) of the event
	return index, nil
}

func (pg pg) writeEventToTable(db *sql.DB, tableName string, streamName string, event string, uuidString, version string) (int64, error) {
	sqlStatement := fmt.Sprintf(sqlInsertEvent, tableName, uuidString, streamName, version, event)
	pg.logger.Debugf("writeEventToTable sql statement = \n%v", sqlStatement)
	result, err := db.Query(sqlStatement)
	if err != nil {
		pg.logger.Error(err)
		return -1, err
	}
	defer result.Close()

	if !result.Next() {
		errString := fmt.Sprintf("No result returned from query: %v", sqlStatement)
		pg.logger.Error(errString)
		return -1, errors.New(errString)
	}

	var position int64
	result.Scan(&position)
	if err != nil {
		pg.logger.Error(err)
		return -1, err
	}

	if err := result.Err(); err != nil {
		pg.logger.Error(err)
		return -1, err
	}

	pg.logger.Infof("Wrote event with uuid [%v] to table [%v], position at [%v]", uuidString, tableName, position)
	return position, nil
}

func (pg pg) createTableForStream(db *sql.DB, tableName string) error {
	sqlStatement := fmt.Sprintf(sqlCreateTable, tableName)
	pg.logger.Debugf("createTableForStream sql statement = \n%v", sqlStatement)
	_, err := db.Exec(sqlStatement)
	if err != nil {
		pg.logger.Error(err)
		return err
	}

	return nil
}

// This function returns the table name from a stream name.  The stream name is assumed to be of the form
// noun:verb-uuid4.  If the stream name is not of this form, an error is returned.
func (pg pg) getTableNameFromStreamName(streamName string) (string, error) {
	categoryIndex := strings.IndexAny(streamName, "-")
	if categoryIndex <= 0 {
		errorString := fmt.Sprintf("Invalid streamName: %v", streamName)
		pg.logger.Error(errorString)
		return "", errors.New(errorString)
	}
	tableName := strings.Replace(streamName[:categoryIndex], ":", "_", -1)

	return tableName, nil
}

// This function identifies whether or not a database table exists in the database.
// Returns true if a table exists, false otherwise
func (pg pg) checkTableExists(db *sql.DB, tableName string) (bool, error) {
	retval := false

	sqlStatement := fmt.Sprintf(sqlCheckTableExists, tableName)
	pg.logger.Debugf("checkTableExists sql statement = \n%v", sqlStatement)
	result, err := db.Query(sqlStatement)
	if err != nil {
		pg.logger.Error(err)
		return false, err
	}
	defer result.Close()

	for result.Next() {
		var exists int
		if err := result.Scan(&exists); err != nil {
			pg.logger.Error(err)
			return false, err
		}
		if exists == 1 {
			retval = true
			break
		}
	}
	if err := result.Err(); err != nil {
		pg.logger.Error(err)
		return false, err
	}

	return retval, nil
}

func (pg pg) ConditionallyWriteEvent(streamName string, event string, version string, previousPosition uint64) error {
	return fmt.Errorf("Unimplemented")
}

func (pg pg) ReadEvent(streamName string, position uint64) (string, error) {
	return "", fmt.Errorf("Unimplemented")
}

func (pg pg) ReadAllEvents(streamName string) (string, error) {
	return "", fmt.Errorf("Unimplemented")
}

func (pg pg) ReadEventsFromPosition(streamName string, position uint64) (string, error) {
	return "", fmt.Errorf("Unimplemented")
}

func (pg pg) ReadEventsFromTime(streamName string, datetime time.Time) (string, error) {
	return "", fmt.Errorf("Unimplemented")
}
