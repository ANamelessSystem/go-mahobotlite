package database

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const testID = 1

// test connect to a DB on local disk
func TestConnectToSqliteDB(t *testing.T) {
	deleteTestDB(t)
	_, err := connectToSqliteDB(testID)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
}

// test init DB on local disk
func TestInitDB(t *testing.T) {
	deleteTestDB(t)
	// 1. Test file creation, dir would not being test here because it is same with prod env
	dbPath := getDBFilePath(testID)

	err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
	fmt.Printf("DB Created in %v", dbPath)

	// 2. Test database connection
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("Expected to connect to database, but got error: %v", err)
	}

	// 3. Test table migrations
	for _, model := range AllModels {
		if exists := db.Migrator().HasTable(model); !exists {
			t.Fatalf("Expected table for model %v to exist, but it does not", model)
		}
	}

	// 4. Test data seeding
	var actualBossInfo []SysBossInfo
	db.Find(&actualBossInfo)
	expectedBossInfo := getBossInfoSeeds()
	if !reflect.DeepEqual(expectedBossInfo, actualBossInfo) {
		t.Fatalf("Expected SysBossInfo to be seeded with %v, but got %v", expectedBossInfo, actualBossInfo)
	}

	var actualInputLimit []SysInputLimit
	db.Find(&actualInputLimit)
	expectedInputLimit := getInputLimitSeeds()
	if !reflect.DeepEqual(expectedInputLimit, actualInputLimit) {
		t.Fatalf("Expected SysInputLimit to be seeded with %v, but got %v", expectedInputLimit, actualInputLimit)
	}
}

func deleteTestDB(t *testing.T) {
	err := os.Remove(getDBFilePath(testID))
	if err != nil {
		if !os.IsNotExist(err) {
			t.Fatalf("Failed to start init DB: %v", err)
		}
	}
}
