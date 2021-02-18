package main

import (
	"fmt"
	"log"
	"os"

	tcm "github.com/gurbos/tcmodels"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	HOST = "127.0.0.1"
	PORT = "8000"
)

// DataSourceName attributes hold informaition about a specific database.
// The information is used to connect to said database.
type DataSourceName struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// DSNString returns a data source identifier string
func (dsn *DataSourceName) DSNString() string {
	format := "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	return fmt.Sprintf(format, dsn.User, dsn.Password, dsn.Host, dsn.Port, dsn.Database)
}

const (
	ProductLineNotFoundErr = "Product line not found"
)

type APIError struct {
	ErrorMsg string `json:"errorMsg"`
}

// GetDataSource returns a DataSourceName object with that specifies
// the data source used during production.
func GetDataSource() DataSourceName {
	godotenv.Load(".env")
	dataSource := DataSourceName{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWD"),
		Database: os.Getenv("DB_NAME"),
	}
	return dataSource
}

func DBConnection(dsn string, logLevel logger.LogLevel) *gorm.DB {
	conn, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{Logger: logger.Default.LogMode(logLevel)})
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func MakeCardSetRepList(list []tcm.SetInfo, pli tcm.ProductLine) []CardSetRep {
	cardSetRepList := make([]CardSetRep, len(list), len(list))
	for i := 0; i < len(list); i++ {
		cardSetRepList[i].Set(list[i])
	}
	return cardSetRepList
}

func SetsLink(productLine string) string {
	setsLink := HOST + ":" + PORT + "/" + productLine + "/sets"
	return setsLink
}
