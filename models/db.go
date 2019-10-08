package models

import (
	"database/sql"
	"fmt"
	"os"

	"order-service/startup"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var Db *sql.DB

func init() {
	log := logrus.WithFields(logrus.Fields{"module": "model"})

	dbName := startup.Config.Database.DbName
	userName := startup.Config.Database.UserName
	password := startup.Config.Database.Password

	dataSource := fmt.Sprintf("%s:%s@tcp(mysql)/%s", userName, password, dbName)
	//dataSource := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", userName, password, dbName)

	var err error
	Db, err = sql.Open("mysql", dataSource)
	if err != nil {
		log.Error("Failed to connect to database")
		os.Exit(1)
	}

	err = Db.Ping()
	if err != nil {
		log.Error("Failed to ping database")
		os.Exit(1)
	}

	log.Info("Successfully connected to database...")
}
