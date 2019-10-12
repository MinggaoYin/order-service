package startup

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var Db *sql.DB

func Init() {
	log := logrus.WithFields(logrus.Fields{"module": "startup"})

	dbName := Config.Database.DbName
	userName := Config.Database.UserName
	password := Config.Database.Password

	dataSource := fmt.Sprintf("%s:%s@tcp(mysql)/%s", userName, password, dbName)
	//dataSource := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", userName, password, dbName)

	dbConn, err := sql.Open(`mysql`, dataSource)
	if err != nil {
		log.Error("Failed to connect to database")
		os.Exit(1)
	}

	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	Db = dbConn

	log.Info("Successfully connected to database...")

	initTables()
}

var createTableStat = `CREATE TABLE IF NOT EXISTS orders (
    id BIGINT(20) UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
    origin_lat DOUBLE NOT NULL,
    origin_lng DOUBLE NOT NULL,
    destination_lat DOUBLE NOT NULL,
    destination_lng DOUBLE NOT NULL,
    status VARCHAR(20) NOT NULL,
    distance INT UNSIGNED NOT NULL
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8;
`

func initTables() {
	// create order table if not exists
	Db.Exec(createTableStat)
}
