package database

import (
	"fmt"
	"path"

	"github.com/smtdfc/nagare/core/persistence/database/models"
	"github.com/smtdfc/nagare/shared/paths"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDatabase() (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}

	var err error
	db, err = gorm.Open(sqlite.Open(path.Join(paths.DatabaseDir, "nagare.db")), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Session{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetDatabase() (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}

	fmt.Println("Error")
	return nil, fmt.Errorf("The database connection has not been established. Ensure InitDatabase is called before performing operations")
}
