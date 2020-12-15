package testinfra

import (
	"github.com/jinzhu/gorm"
	"log"
	"strings"
)

func DropDatabase(db *gorm.DB, databaseToDrop string) {
	if strings.Contains(strings.ToLower(databaseToDrop), "test") {
		if err := db.Exec("DROP DATABASE " + databaseToDrop + "").Error; err != nil {
			log.Fatalln("failed to drop database: " + databaseToDrop)
		} else {
			log.Println("databased dropped: " + databaseToDrop)
		}
	}
}
