package dbase

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
  
func db() string {
  dsn := "host=localhost user=wolf-dev-test password=edummypword dbname=dev-test port=5433 sslmode=disable TimeZone=America/Los_Angeles"
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

  if err != nil {
	  fmt.Println("Error")
	  return "False"
  }
  return "True"
}