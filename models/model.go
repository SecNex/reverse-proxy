package models

import (
	"time"

	"gorm.io/gorm"
)

type Website struct {
	gorm.Model
	Domain   string `gorm:"uniqueIndex;not null"`
	Protocol string `gorm:"not null;default:'http'"`
	Host     string `gorm:"not null"`
	Port     int    `gorm:"not null;default:80"`
	SSL      bool   `gorm:"default:false"`
	Active   bool   `gorm:"default:true"`
	Email    string `gorm:"not null"`
	LastSeen time.Time
}

type WebsiteConfig struct {
	Domain   string
	Protocol string
	Host     string
	Port     int
	SSL      bool
	Active   bool
	Email    string
}
