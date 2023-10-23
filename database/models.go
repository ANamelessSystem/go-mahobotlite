package database

import "time"

type MemberInfo struct {
	MemberID   int `gorm:"primaryKey"`
	MemberName string
}

type AttackRecord struct {
	RecID    int `gorm:"primaryKey"`
	MemberID int
	Damage   int
	RecRound int
	RecBoss  int
	RecType  int
	RecTime  time.Time
	Uploader int
}

type Queue struct {
	MemberID  int
	JoinBoss  int `gorm:"primaryKey"`
	JoinSeq   int `gorm:"primaryKey"`
	JoinType  int
	JoinRound int // JoinRound will only need when sos
	JoinTime  time.Time
}

type SysInputLimit struct {
	InputType  string `gorm:"primaryKey"`
	InputValue int
}

type SysBossInfo struct {
	Boss     int `gorm:"primaryKey"`
	InitHP   int
	RoundMin int `gorm:"primaryKey"`
	RoundMax int `gorm:"primaryKey"`
	SortFlag int
}

var AllModels = []interface{}{
	&MemberInfo{},
	&AttackRecord{},
	&Queue{},
	&SysInputLimit{},
	&SysBossInfo{},
}

// Progress data, not in database
type TempProgress struct {
	Boss     int
	Round    int
	Damage   int
	HP       int
	RoundMin int
	RoundMax int
}

type QueueWithUsername struct {
	MemberID   int
	JoinBoss   int
	JoinSeq    int
	JoinType   int
	JoinRound  int
	JoinTime   time.Time
	MemberName string
}
