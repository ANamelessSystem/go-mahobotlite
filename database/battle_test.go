package database

import (
	"fmt"
	"testing"
)

func TestGetProgress(t *testing.T) {
	db, err := connectToSqliteDB(testID)
	if err != nil {
		t.Fatalf("Unable to connect to test DB: %v", err)
	}
	err = db.Where("1 = 1").Delete(&AttackRecord{}).Error
	if err != nil {
		t.Fatalf("Unable to truncate table AttackRecord: %v", err)
	}
	// populate test data
	attackRecord := getAttackRecordsSample()
	if len(attackRecord) > 0 {
		err = db.Create(&attackRecord).Error
		if err != nil {
			t.Fatalf("Unable to insert sample data to test DB: %v", err)
		}
	}
	// Call the function
	results, err := GetProgress(testID)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	// Validate the results
	if len(results) != 5 {
		t.Fatalf("Expected 5 result, but got %d", len(results))
	}
	fmt.Printf("Boss, Round, Damage, HP, RoundMin, RoundMax")
	for _, v := range results {
		fmt.Printf("\n%v", v)
		if v.Damage >= v.HP {
			t.Fatalf("Expected Damage(%d) smaller than HP(%d), but got bigger or equal.", v.Damage, v.HP)
		}
	}
}
