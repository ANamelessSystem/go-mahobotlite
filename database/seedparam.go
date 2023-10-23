package database

import "gorm.io/gorm"

func getBossInfoSeeds() []SysBossInfo {
	return []SysBossInfo{
		{SortFlag: 11, Boss: 1, RoundMin: 1, RoundMax: 3, InitHP: 6000000},
		{SortFlag: 12, Boss: 2, RoundMin: 1, RoundMax: 3, InitHP: 8000000},
		{SortFlag: 13, Boss: 3, RoundMin: 1, RoundMax: 3, InitHP: 10000000},
		{SortFlag: 14, Boss: 4, RoundMin: 1, RoundMax: 3, InitHP: 12000000},
		{SortFlag: 15, Boss: 5, RoundMin: 1, RoundMax: 3, InitHP: 15000000},
		{SortFlag: 21, Boss: 1, RoundMin: 4, RoundMax: 9, InitHP: 8000000},
		{SortFlag: 22, Boss: 2, RoundMin: 4, RoundMax: 9, InitHP: 10000000},
		{SortFlag: 23, Boss: 3, RoundMin: 4, RoundMax: 9, InitHP: 13000000},
		{SortFlag: 24, Boss: 4, RoundMin: 4, RoundMax: 9, InitHP: 15000000},
		{SortFlag: 25, Boss: 5, RoundMin: 4, RoundMax: 9, InitHP: 20000000},
		{SortFlag: 31, Boss: 1, RoundMin: 10, RoundMax: 25, InitHP: 20000000},
		{SortFlag: 32, Boss: 2, RoundMin: 10, RoundMax: 25, InitHP: 22000000},
		{SortFlag: 33, Boss: 3, RoundMin: 10, RoundMax: 25, InitHP: 25000000},
		{SortFlag: 34, Boss: 4, RoundMin: 10, RoundMax: 25, InitHP: 28000000},
		{SortFlag: 35, Boss: 5, RoundMin: 10, RoundMax: 25, InitHP: 30000000},
		{SortFlag: 41, Boss: 1, RoundMin: 26, RoundMax: 99, InitHP: 200000000},
		{SortFlag: 42, Boss: 2, RoundMin: 26, RoundMax: 99, InitHP: 210000000},
		{SortFlag: 43, Boss: 3, RoundMin: 26, RoundMax: 99, InitHP: 230000000},
		{SortFlag: 44, Boss: 4, RoundMin: 26, RoundMax: 99, InitHP: 240000000},
		{SortFlag: 45, Boss: 5, RoundMin: 26, RoundMax: 99, InitHP: 250000000},
	}
}

func getInputLimitSeeds() []SysInputLimit {
	return []SysInputLimit{
		{InputType: "Damage", InputValue: 250000000},
		{InputType: "Round", InputValue: 99},
		{InputType: "Boss", InputValue: 5},
	}
}

func seedSysBossInfo(db *gorm.DB) error {
	// Set init HPs
	var count int64
	db.Model(&SysBossInfo{}).Count(&count)
	sets := getBossInfoSeeds()
	if count == 0 {
		return db.Create(&sets).Error
	} else {
		return ErrDataExists
	}
}

func seedSysInputLimits(db *gorm.DB) error {
	// Set input limits
	var count int64
	db.Model(&SysInputLimit{}).Count(&count)
	sets := getInputLimitSeeds()
	if count == 0 {
		return db.Create(&sets).Error
	} else {
		return ErrDataExists
	}
}
