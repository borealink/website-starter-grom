package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Database wraps the GORM database instance.
type Database struct {
	db *gorm.DB
}

// NewDatabase creates and initializes a new SQLite database connection.
func NewDatabase(
	path string,
) (*Database, error) {

	// Open the SQLite database using GORM.
	db, err := gorm.Open(
		sqlite.Open(path),
		&gorm.Config{},
	)

	// Return error if connection fails.
	if err != nil {
		return nil, err
	}

	// Return the database wrapper instance.
	return &Database{
		db: db,
	}, nil
}

// Get returns the underlying GORM database instance.
func (database *Database) Get() *gorm.DB {
	return database.db
}
