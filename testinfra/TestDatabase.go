package testinfra

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"hallo/dataSource"
	"os"
	"strings"
)

type TemporaryDatabase struct {
	Ds           *dataSource.DataSource
	Database     *gorm.DB
	DatabaseName string
}

func NewTemporaryDatabase() *TemporaryDatabase {
	mysqlServer := os.Getenv("TEST_MYSQL_SERVER")
	if mysqlServer == "" {
		mysqlServer = "root:root@(127.0.0.1:3306)"
	}

	databaseName := "hallo_test_" + strings.ReplaceAll(uuid.New().String(), "-", "")
	if err := os.Setenv("DATABASE_URL", "mysql://"+mysqlServer+"/"+databaseName+"?charset=utf8mb4&parseTime=True&loc=Local"); err != nil {
		panic(err)
	}

	ds, err := new(dataSource.DataSource).Start()
	if err != nil {
		if ds != nil {
			ds.Stop()
		}
		panic(err)
	}

	return &TemporaryDatabase{Ds: ds, Database: ds.Database, DatabaseName: databaseName}
}

func (ds *TemporaryDatabase) CleanAndDisconnect() {
	DropDatabase(ds.Database, ds.DatabaseName)
	ds.Ds.Stop()
}
