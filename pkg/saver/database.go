package saver

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/defectus/glutton/pkg/iface"
)

const defaultLayout = "INSERT INTO payload(ts, remote, meta, payload) VALUES ($1, $2, $3, $4)"

// DatabaseSaver can save payload to database.
type DatabaseSaver struct {
	db     *sql.DB
	layout string
}

// Configure configures this instance of DatabaseSaver.
func (ds *DatabaseSaver) Configure(settings *iface.Settings) (err error) {
	ds.db, err = sql.Open(settings.SQLDriver, settings.SQLConnectionString)
	ds.layout = settings.SQLayout
	if len(ds.layout) == 0 {
		ds.layout = defaultLayout
	}
	return
}

// Save stored data into the database.
func (ds *DatabaseSaver) Save(payload *iface.PayloadRecord) error {
	_, err := ds.db.Exec(ds.layout, payload.Timestamp, payload.Remote, payload.Meta, payload.Payload)
	return errors.Wrap(err, "error saving payload")
}
