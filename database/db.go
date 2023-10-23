package database

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Get sqlite connection by group ID
func connectToSqliteDB(grpid int) (*gorm.DB, error) {
	// get path by grpid
	path := getDBFilePath(grpid)

	// check if need to create a new db file
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := InitDB(path); err != nil {
			return nil, err
		}
	}

	// connect to db file and return a db instance
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// 本函数用于群组第一次使用时创建全新的数据库, 并设定基础参数
func InitDB(path string) error {
	// check and create dir
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		err = os.MkdirAll(dbDir, 0755)
		if err != nil {
			return err
		}
	}
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error when opening database file(%s): %v", path, err)
	}

	// 所有需要创建的表结构（模型）已注册到AllModels公共变量中
	err = db.AutoMigrate(AllModels...)
	if err != nil {
		return fmt.Errorf("error when migrating table structure: %v", err)
	}

	// Seed all sys-prefix tables
	err = seedSysBossInfo(db)
	if err != nil {
		return fmt.Errorf("error when seed boss info: %v", err)
	}
	err = seedSysInputLimits(db)
	if err != nil {
		return fmt.Errorf("error when seed input limit: %v", err)
	}
	return nil
}

// 本函数用于应对新一轮公会战时, 重置伤害记录和队列
func ResetRecords(grpid int) error {
	// connect to db file
	db, err := connectToSqliteDB(grpid)
	if err != nil {
		return err
	}

	// 重置前归档现有数据库
	err = archiveDBFile(grpid)
	if err != nil {
		return err
	}
	// 执行数据操作语句, 清理伤害记录表及队列表
	err = db.Where("1 = ?", "1").Delete(&AttackRecord{}).Error
	if err != nil {
		return err
	}
	err = db.Delete(&Queue{}).Error
	if err != nil {
		return err
	}
	return nil
}

// 本函数用于应对新一轮公会战时, 重置包括成员列表在内的所有数据
func ResetAll(grpid int) error {
	// connect to db file
	db, err := connectToSqliteDB(grpid)
	if err != nil {
		return err
	}

	// 删除记录
	err = ResetRecords(grpid)
	if err != nil {
		return err
	}
	// 删除成员名单
	err = db.Delete(&MemberInfo{}).Error
	if err != nil {
		return err
	}
	return nil
}
