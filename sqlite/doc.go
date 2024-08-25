// Package sqlite implements a connect hook around the sqlite3 driver so that the underlying connection can be fetched from the driver for more advanced operations such as backups. See: https://github.com/mattn/go-sqlite3/blob/master/_example/hook/hook.go.
// To use make sure you import this package so that the init code registers the driver: import _ github.com/theopenlane/utils/sqlite
// Then you can use sql.Open in the same way you would with sqlite3: sql.Open("_sqlite3", "path/to/database.db")
package sqlite
