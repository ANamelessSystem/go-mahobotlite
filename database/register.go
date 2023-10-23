package database

// Return row counts in member_infos
func CountMemberInfoByGroup(grpID int) (int64, error) {
	db, err := connectToSqliteDB(grpID)
	if err != nil {
		return 0, err
	}
	var totalCount int64
	err = db.Table("member_infos").Count(&totalCount).Error
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

// Get all rows in member_infos and return a struct array
func GetMemberInfoByGroup(grpID int) ([]MemberInfo, error) {
	db, err := connectToSqliteDB(grpID)
	if err != nil {
		return nil, err
	}

	var memberInfo []MemberInfo
	err = db.Table("member_infos").Find(&memberInfo).Error
	if err != nil {
		return nil, err
	}
	return memberInfo, nil
}

// Add a member to member_infos
func JoinMemberInfo(grpID int, memberID int, memberName string) (int64, error) {
	db, err := connectToSqliteDB(grpID)
	if err != nil {
		return 0, err
	}

	exists, err := IsMemberInfoExists(grpID, memberID)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, db.Table("member_infos").Where("member_id = ?", memberID).Update("member_name", memberName).Error
	}

	totalCount, err := CountMemberInfoByGroup(grpID)
	if err != nil {
		return 0, err
	}
	if totalCount >= 30 {
		return totalCount, ErrDataExceeded
	}
	memberInfo := MemberInfo{
		MemberID:   memberID,
		MemberName: memberName,
	}
	err = db.Table("member_infos").Create(&memberInfo).Error
	if err != nil {
		return 0, err
	}
	return totalCount + 1, nil
}

// Check if member in the MemberInfo table
func IsMemberInfoExists(grpID int, memberID int) (bool, error) {
	// connect to db file and return a db instance
	db, err := connectToSqliteDB(grpID)
	if err != nil {
		return false, err
	}

	var count int64
	db.Table("member_infos").Where("member_id = ?", memberID).Count(&count)
	switch count {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, ErrDataNotUnique
	}
}
