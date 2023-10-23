package database

import (
	"time"
)

// Quit queue by userID and provided boss, if boss = 0 then quit all queue
func QuitQueue(grpID int, userID int, boss int) error {
	db, err := connectToSqliteDB(grpID)
	if err != nil {
		return err
	}

	if boss != 0 {
		err = db.Where("member_id = ? AND join_boss = ?", userID, boss).Delete(&Queue{}).Error
	} else {
		err = db.Where("member_id = ?", userID).Delete(&Queue{}).Error
	}
	if err != nil {
		return err
	}

	return nil
}

// Return rows by giving userID
func GetQueueCountByUserID(grpID int, userID int) ([]int, error) {
	db, err := connectToSqliteDB(grpID)
	if err != nil {
		return nil, err
	}

	var joinBosses []int
	err = db.Model(&Queue{}).Where("member_id = ?", userID).Distinct("join_boss").Pluck("join_boss", &joinBosses).Error
	if err != nil {
		return nil, err
	}
	return joinBosses, nil
}

// Return progress slice contains all 5 bosses, the progress should as close to in-game as possible.
func GetProgress(grpID int) ([]TempProgress, error) {
	db, err := connectToSqliteDB(grpID)
	if err != nil {
		return nil, err
	}

	var attackRecordsGrouped []struct {
		Boss      int
		MaxRound  int
		SumDamage int
	}

	err = db.Raw(`
					WITH round_grouped AS (
						SELECT rec_boss as boss, MAX(rec_round) AS max_round
						FROM attack_records
						GROUP BY rec_boss
					)
					SELECT 
						a.rec_boss AS boss, 
						m.max_round, 
						SUM(a.damage) AS sum_damage
					FROM 
						attack_records a
					JOIN 
					round_grouped m
					ON 
						a.rec_boss = m.boss AND a.rec_round = m.max_round
					GROUP BY 
						a.rec_boss, m.max_round;
				`).Scan(&attackRecordsGrouped).Error
	if err != nil {
		return nil, err
	}

	var results []TempProgress
	for _, record := range attackRecordsGrouped {
		var bossInfo SysBossInfo
		err := db.Model(&SysBossInfo{}).Where("boss = ? AND ? BETWEEN round_min AND round_max", record.Boss, record.MaxRound).Scan(&bossInfo).Error
		if err != nil {
			return nil, err
		}

		tempProgress := TempProgress{
			Boss:     record.Boss,
			Round:    record.MaxRound,
			Damage:   record.SumDamage,
			HP:       bossInfo.InitHP,
			RoundMin: bossInfo.RoundMin,
			RoundMax: bossInfo.RoundMax,
		}

		// Check if damage exceeds or is close to HP
		if tempProgress.Damage >= tempProgress.HP || (tempProgress.HP-tempProgress.Damage) < 10000 {
			// Increment round and reset damage
			tempProgress.Round++
			tempProgress.Damage = 0

			// Fetch HP for the next round
			err := db.Model(&SysBossInfo{}).
				Where("boss = ? AND ? BETWEEN round_min AND round_max", tempProgress.Boss, tempProgress.Round).
				First(&bossInfo).Error
			if err != nil {
				return nil, err
			}
			tempProgress.HP = bossInfo.InitHP
			tempProgress.RoundMin = bossInfo.RoundMin
			tempProgress.RoundMax = bossInfo.RoundMax
		}

		results = append(results, tempProgress)
	}

	// If there are no attack records, populate results with initial data for each boss
	if len(attackRecordsGrouped) == 0 {
		var allBosses []SysBossInfo
		err := db.Model(&SysBossInfo{}).Where("round_min = 1").Find(&allBosses).Error
		if err != nil {
			return nil, err
		}

		for _, bossInfo := range allBosses {
			tempProgress := TempProgress{
				Boss:     bossInfo.Boss,
				Round:    1,
				Damage:   0,
				HP:       bossInfo.InitHP,
				RoundMin: bossInfo.RoundMin,
				RoundMax: bossInfo.RoundMax,
			}
			results = append(results, tempProgress)
		}
	}
	return results, nil
}

// Get queue, while bossCode invalid, query all bosses
func GetQueue(grpID int, bossCode int) ([]QueueWithUsername, error) {
	db, err := connectToSqliteDB(grpID)
	if err != nil {
		return nil, err
	}
	var queueWithUsername []QueueWithUsername

	if bossCode > 5 || bossCode < 1 {
		// Query all bosses
		err = db.Table("queues").
			Joins("left join member_infos on member_infos.member_id = queues.member_id").
			Order("queues.join_boss ASC, queues.join_seq ASC").
			Select("queues.*, member_infos.member_name").
			Find(&queueWithUsername).Error
	} else {
		err = db.Table("queues").
			Joins("left join member_infos on member_infos.member_id = queues.member_id").
			Where("queues.join_boss = ?", bossCode).
			Order("queues.join_seq ASC").
			Select("queues.*, member_infos.member_name").
			Find(&queueWithUsername).Error
	}

	if err != nil {
		return nil, err
	}
	return queueWithUsername, nil
}

// Add queue
func Enqueue(grpID int, userID int, joinBoss int, joinType int, joinRound int) error {
	// connect to db file and return a db instance
	db, err := connectToSqliteDB(grpID)
	if err != nil {
		return err
	}

	// Get max sequence number on specific boss
	var maxSeq int
	err = db.Table("queues").Where("join_boss = ?", joinBoss).Select("IFNULL(MAX(join_seq), 0)").Row().Scan(&maxSeq)
	if err != nil {
		return err
	}

	queue := Queue{MemberID: userID, JoinBoss: joinBoss,
		JoinSeq: maxSeq + 1, JoinType: joinType,
		JoinRound: 0, JoinTime: time.Now()}
	return db.Create(&queue).Error
}
