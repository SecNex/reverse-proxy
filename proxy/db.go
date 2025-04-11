package proxy

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/secnex/reverse-proxy/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBManager struct {
	db *gorm.DB
}

func NewDBManager() (*DBManager, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "secnex"
	}
	dbsslmode := os.Getenv("DB_SSLMODE")
	if dbsslmode == "" {
		dbsslmode = "disable"
	} else if dbsslmode == "true" {
		dbsslmode = "enable"
	} else if dbsslmode == "false" {
		dbsslmode = "disable"
	}

	if os.Getenv("DB_RESET") == "true" {
		log.Printf("Connecting to database %s:%s/postgres...", host, port)
		postgresDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
			host, port, user, password, dbsslmode)

		postgresDB, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to postgres database: %v", err)
		}
		log.Printf("Connected to database %s:%s/postgres.", host, port)

		log.Printf("Dropping all connections to database %s...", dbname)
		postgresDB.Exec("SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = '" + dbname + "';")
		log.Printf("All connections to database %s dropped.", dbname)

		log.Printf("Dropping database %s...", dbname)
		postgresDB.Exec("DROP DATABASE IF EXISTS " + dbname + ";")
		log.Printf("Database %s dropped.", dbname)

		log.Printf("Creating database %s...", dbname)
		postgresDB.Exec("CREATE DATABASE " + dbname + ";")
		log.Printf("Database %s created.", dbname)
	}

	log.Printf("Connecting to database %s:%s/%s...", host, port, dbname)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, dbsslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully!")

	log.Println("Migrating database...")

	err = db.AutoMigrate(&models.Website{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("Database migrated successfully!")

	return &DBManager{db: db}, nil
}

func (dm *DBManager) GetWebsite(domain string) (*models.Website, error) {
	var website models.Website
	result := dm.db.Where("domain = ?", domain).First(&website)
	if result.Error != nil {
		return nil, result.Error
	}
	return &website, nil
}

func (dm *DBManager) GetAllWebsites() ([]models.Website, error) {
	var websites []models.Website
	result := dm.db.Find(&websites)
	if result.Error != nil {
		return nil, result.Error
	}
	return websites, nil
}

func (dm *DBManager) CreateWebsite(config models.WebsiteConfig) error {
	website := models.Website{
		Domain:   config.Domain,
		Protocol: config.Protocol,
		Host:     config.Host,
		Port:     config.Port,
		SSL:      config.SSL,
		Active:   config.Active,
		LastSeen: time.Now(),
	}
	result := dm.db.Create(&website)
	return result.Error
}

func (dm *DBManager) UpdateWebsite(domain string, config models.WebsiteConfig) error {
	result := dm.db.Model(&models.Website{}).Where("domain = ?", domain).Updates(map[string]interface{}{
		"protocol":  config.Protocol,
		"host":      config.Host,
		"port":      config.Port,
		"ssl":       config.SSL,
		"active":    config.Active,
		"last_seen": time.Now(),
	})
	return result.Error
}

func (dm *DBManager) DeleteWebsite(domain string) error {
	result := dm.db.Where("domain = ?", domain).Delete(&models.Website{})
	return result.Error
}
