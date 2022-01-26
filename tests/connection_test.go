package tests_test

import (
	"fmt"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"
)

func TestWithSingleConnection(t *testing.T) {
	var expectedName = "test"
	var actualName string

	setSQL, getSQL := getSetSQL(DB.Dialector.Name())
	if len(setSQL) == 0 || len(getSQL) == 0 {
		return
	}

	err := DB.Connection(func(tx *gorm.DB) error {
		if err := tx.Exec(setSQL, expectedName).Error; err != nil {
			return err
		}

		if err := tx.Raw(getSQL).Scan(&actualName).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		t.Errorf(fmt.Sprintf("WithSingleConnection should work, but got err %v", err))
	}

	if actualName != expectedName {
		t.Errorf("WithSingleConnection() method should get correct value, expect: %v, got %v", expectedName, actualName)
	}
}

func getSetSQL(driverName string) (string, string) {
	switch driverName {
	case mysql.Dialector{}.Name():
		return "SET @testName := ?", "SELECT @testName"
	case postgres.Dialector{}.Name():
		return "SET test.test_name = ?", "SELECT current_setting('test.test_name')"
	default:
		return "", ""
	}
}

func TestConnectionNewSessionMode(t *testing.T) {
	err := DB.Connection(func(tx *gorm.DB) error {
		user := *GetUser("connection_new_session_mode", Config{Account: true})

		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		if err := tx.Preload("Account").Where("id = ?", user.ID).First(&user).Error; err != nil {
			return err
		}

		if err := tx.Where("number = ?", user.Account.Number).First(&tests.Account{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		t.Errorf(fmt.Sprintf("ConnectionNewSessionMode should work, but got err %v", err))
	}
}
