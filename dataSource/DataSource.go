package dataSource

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"hallo/infra"
	"log"
	"os"
	"strings"
)

type DataSource struct {
	Database *gorm.DB
}

func (ds *DataSource) Start() (*DataSource, error) {
	// mysql://user:secret@(mysql-svc:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	databaseUrl := os.ExpandEnv(os.Getenv("DATABASE_URL"))

	if databaseUrl == "" {
		// databaseUrl = "sqlite3://file::memory:?cache=shared" // sqlite3:///tmp/gorm.db
		return nil, errors.New("database url is empty")
	}

	slice := strings.Split(databaseUrl, "://")
	driver := slice[0]
	driverArgs := slice[1]

	if driver == "mysql" {
		prepareMysqlDatabase(driverArgs)
	}

	ds.Database = connect(driver, driverArgs)
	if strings.ToUpper(os.Getenv("ENABLE_DEBUG")) == "TRUE" {
		ds.Database.LogMode(true)
	}

	infra.Migrate(ds.Database)

	return ds, nil
}

func (ds *DataSource) Stop() {
	if ds.Database != nil {
		err := ds.Database.Close()
		if err != nil {
			log.Fatalf("failed to close DB: %v", err)
		}
		ds.Database = nil
	}
}

func connect(driver, driverArgs string) *gorm.DB {
	db, err := gorm.Open(driver, driverArgs)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	err = db.DB().Ping()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	return db
}

func prepareMysqlDatabase(mysqlDriverArgs string) {
	// root:xxx@(test.xxx.com:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	databaseName, rootDriverArgs := extractDatabaseName(mysqlDriverArgs)

	db, err := gorm.Open("mysql", rootDriverArgs)
	if err != nil {
		log.Fatalf("[prepare] failed to open database: %v", err)
	}
	err = db.DB().Ping()
	if err != nil {
		log.Fatalf("[prepare] failed to connect to database: %v", err)
	}
	initSql := "CREATE DATABASE IF NOT EXISTS `" + databaseName + "` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_unicode_ci;"
	err = db.Exec(initSql).Error
	if err != nil {
		log.Fatalf("[prepare] failed to prepare database with sql '%s': %v", initSql, err)
	}
}

func extractDatabaseName(mysqlDriverArgs string) (string, string) {
	nameIndex := strings.IndexRune(mysqlDriverArgs, '/')
	paramsIndex := strings.IndexRune(mysqlDriverArgs, '?')

	if nameIndex > 0 && paramsIndex > nameIndex {
		return mysqlDriverArgs[nameIndex+1 : paramsIndex], mysqlDriverArgs[0:nameIndex+1] + mysqlDriverArgs[paramsIndex:]
	}
	if nameIndex < 0 {
		return "", mysqlDriverArgs
	}
	if nameIndex > 0 && paramsIndex < 0 {
		return mysqlDriverArgs[nameIndex+1:], mysqlDriverArgs[0 : nameIndex+1]
	}
	log.Fatalf("bad mysql driver args %s", mysqlDriverArgs)
	return "", ""
}
