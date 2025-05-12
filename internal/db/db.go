package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Log struct {
	ID        uint `gorm:"primaryKey"`
	Project   string
	Status    int
	Message   string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Project struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique"`
	Url       string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

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

	err = db.AutoMigrate(&Project{})
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
	oneHourAgo := time.Now().Add(-1 * time.Hour)

	if err := DB.Find(&logs).Error; err != nil {
		fmt.Println(err)
	}

	for _, log := range logs {
		if log.CreatedAt.Before(oneHourAgo) || len(groupedLogs[log.Project]) == 60 {
			continue
		}

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

func CleanupLogs() {
	result := DB.Where("created_at < datetime('now', '-1 hour')").Delete(&Log{})

	if result.Error != nil {
		log.Printf("[CRON] Failed to delete old logs: %v", result.Error)
	} else {
		log.Printf("[CRON] Deleted %d old logs", result.RowsAffected)
	}
}

func CreateProject(project Project) {
	result := DB.Create(&project)

	if result.Error != nil {
		fmt.Println(result.Error)
	}
}

func ReadProjects() []Project {
	var projects []Project

	if err := DB.Find(&projects).Error; err != nil {
		fmt.Println(err)
	}

	return projects
}
