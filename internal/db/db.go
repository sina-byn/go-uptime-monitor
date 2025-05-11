package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Log struct {
	ID uint `gorm:"primaryKey"`
	Project string
	Status uint
	Message string
	CreatedAt time.Time `gorm:"autoCreateTime"`
};

var DB *gorm.DB

func InitDB() {
	db, err := gorm.Open(sqlite.Open("uptime.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database")
	}

	DB = db

	err = db.AutoMigrate(&Log{})
	if err != nil {
		log.Fatalf("failed to migrate to schema")
	}
}

func CreateLog(log Log) {
	result := DB.Create(&log)
	
	if result.Error != nil {
		fmt.Println(result.Error)
	}
}

func ReadLog(id uint) *Log {
	var log Log

	if err := DB.First(&log, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("Log not found")
		} else {
			fmt.Println("Error retrieving log:", err)
		}

		return nil
	}

	return &log
}

func ReadLogs() map[string][]Log {
	var logs []Log
	groupedLogs := make(map[string][]Log)

	if err := DB.Group("project").Find(&logs).Error; err != nil {
		fmt.Println(err)
	}

	for _, log := range(logs) {
		groupedLogs[log.Project] = append(groupedLogs[log.Project], log)
	}
	
	return groupedLogs
}

func ReadProjectLogs(project string) []Log {
	var logs []Log

	if err := DB.Where("project = ?", project).Find(&logs).Error; err != nil {
		fmt.Println(err)
	}

	return logs
}