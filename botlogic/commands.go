package botlogic

import (
	"fmt"
	"time"

	"github.com/ANamelessSystem/go-mahobotlite/database"
	"github.com/sirupsen/logrus"
)

// Add member into member list(register)
func enroll(messageIn *GroupMessageIn) {
	messageOut := &MessageBuilder{}
	memberName := getUserName(messageIn)
	memberCount, err := database.JoinMemberInfo(messageIn.GroupID, messageIn.Sender.UserID, memberName)
	if err != nil {
		if err == database.ErrDataExceeded {
			messageOut.AddPart(&TextPart{Text: "成员名单列表已到达上限, 无法加入"})
		} else if err == database.ErrDataNotUnique {
			messageOut.AddPart(&TextPart{Text: "该成员在名单列表中存在重复数据, 请手动处理"})
		} else {
			logrus.Errorf("Cannot enroll due to error:%v", err)
			messageOut.AddPart(&TextPart{Text: fmt.Sprintf("无法加入, 错误：%s", err)})
		}
	} else {
		messageOut.AddPart(&TextPart{Text: fmt.Sprintf("更新成员名单成功(%v/30)", memberCount)})
	}

	sendGroupMessage(messageIn.GroupID, messageOut)
}

// Add member to queue for prepare to attack
func enqueue(messageIn *GroupMessageIn, bcp *botCmdParams) {
	// all error lifecycle should end in command execute layer, and should not keep throwing
	messageOut := &MessageBuilder{}
	// Check member valid
	exists, err := database.IsMemberInfoExists(messageIn.GroupID, messageIn.Sender.UserID)
	if err != nil {
		logrus.Errorf("Cannot get member info due to error:%v", err)
		messageOut.AddPart(&TextPart{Text: fmt.Sprintf("无法加入, 错误：%s", err)})
		sendGroupMessage(messageIn.GroupID, messageOut)
		return
	}
	if !exists {
		messageOut.AddPart(&TextPart{Text: "未注册成员, 无法加入队列"})
		sendGroupMessage(messageIn.GroupID, messageOut)
		return
	}

	if bcp.AttackBoss == 0 {
		messageOut.AddPart(&TextPart{Text: "未指定boss, 无法加入队列"})
		sendGroupMessage(messageIn.GroupID, messageOut)
		return
	}
	err = database.Enqueue(messageIn.GroupID, bcp.UserID, bcp.AttackBoss, bcp.AttackType, bcp.AttackRound)
	if err != nil {
		logrus.Errorf("Cannot enqueue due to error:%v", err)
		messageOut.AddPart(&TextPart{Text: fmt.Sprintf("无法加入, 错误：%s", err)})
		sendGroupMessage(messageIn.GroupID, messageOut)
		return
	}
	messageOut.AddPart(&TextPart{Text: "已加入队列"})
	sendGroupMessage(messageIn.GroupID, messageOut)
	showQueue(messageIn, bcp)
}

// List members in queue, also used by enqueue and quitQueue
func showQueue(messageIn *GroupMessageIn, bcp *botCmdParams) {
	messageOut := &MessageBuilder{}
	queues, err := database.GetQueue(messageIn.GroupID, bcp.AttackBoss)
	if err != nil {
		logrus.Errorf("Cannot get queue due to error:%v", err)
		messageOut.AddPart(&TextPart{Text: fmt.Sprintf("无法查询队列, 错误：%s", err)})
		sendGroupMessage(messageIn.GroupID, messageOut)
		return
	}
	progress, err := database.GetProgress(messageIn.GroupID)
	if err != nil {
		logrus.Errorf("Cannot get progress due to error:%v", err)
		messageOut.AddPart(&TextPart{Text: fmt.Sprintf("无法查询进度, 错误：%s", err)})
		sendGroupMessage(messageIn.GroupID, messageOut)
		return
	}

	for boss := 1; boss <= 5; boss++ {
		countTotal := 0
		countExt := 0
		countSos := 0
		textExtime := ""
		textNormal := ""
		textProgress := ""
		for _, q := range queues {
			if q.JoinBoss == boss {
				countTotal += 1
				elapsedMinutes := time.Since(q.JoinTime).Minutes()
				switch q.JoinType {
				case 1:
					countExt += 1
					textExtime += fmt.Sprintf("【补时】%s(%d)", q.MemberName, q.MemberID)
					if elapsedMinutes > 10 {
						textExtime += fmt.Sprintf("\t(已等待%.f分钟)", elapsedMinutes)
					}
					textExtime += "\n"
				case 2:
					countSos += 1
				default:
					textNormal += fmt.Sprintf("【%d】%s(%d)", countTotal-countExt-countSos, q.MemberName, q.MemberID)
					if elapsedMinutes > 10 {
						textNormal += fmt.Sprintf("\t(已等待%.f分钟)", elapsedMinutes)
					}
					textNormal += "\n"
				}
			}
		}
		for _, p := range progress {
			remain := p.HP - p.Damage
			t := ""
			if remain >= 100000000 {
				yi := remain / 100000000
				wan := (remain % 100000000) / 10000
				if wan > 0 {
					t = fmt.Sprintf("%d亿%d万", yi, wan)
				} else {
					t = fmt.Sprintf("%d亿", yi)
				}
			} else if remain >= 10000 {
				t = fmt.Sprintf("%d万", remain/10000)
			} else {
				t = fmt.Sprintf("%d", remain)
			}
			if p.Round == p.RoundMin {
				textProgress += fmt.Sprintf(", %d周目(新难度!), 剩余：%s", p.Round, t)
			} else {
				textProgress += fmt.Sprintf(", %d周目, 剩余：%s", p.Round, t)
			}
		}
		messageOut.AddPart(&TextPart{Text: fmt.Sprintf("B%d%s", boss, textProgress)})
		if countTotal > 0 {
			switch {
			case countSos == countTotal:
				messageOut.AddPart(&TextPart{Text: fmt.Sprintf("【%d人等待救援】\n活动队列中无人。\n", countSos)})
			case countSos > 0:
				messageOut.AddPart(&TextPart{Text: fmt.Sprintf("【%d人等待救援】\n%s%s", countSos, textExtime, textNormal)})
			default:
				messageOut.AddPart(&TextPart{Text: fmt.Sprintf("\n%s%s", textExtime, textNormal)})
			}
		} else {
			messageOut.AddPart(&TextPart{Text: "\n活动队列中无人。\n"})
		}
		if bcp.IsAll {
			messageOut.AddPart(&TextPart{Text: "\n"})
		}
	}
	sendGroupMessage(messageIn.GroupID, messageOut)
}

func quitQueue(messageIn *GroupMessageIn, bcp *botCmdParams) {
	messageOut := &MessageBuilder{}

	sendErrorMessage := func(err error) {
		logrus.Errorf("Cannot quit queue due to error:%v", err)
		messageOut.AddPart(&TextPart{Text: fmt.Sprintf("无法退出队伍, 错误：%s", err)})
		sendGroupMessage(messageIn.GroupID, messageOut)
	}

	if bcp.IsAll {
		bcp.AttackBoss = 0
		err := database.QuitQueue(messageIn.GroupID, bcp.UserID, bcp.AttackBoss)
		if err != nil {
			sendErrorMessage(err)
			return
		}
	}

	queuedBosses, err := database.GetQueueCountByUserID(messageIn.GroupID, bcp.UserID)
	if err != nil {
		sendErrorMessage(err)
		return
	}
	if len(queuedBosses) == 0 {
		messageOut.AddPart(&TextPart{Text: "未加入任何Boss, 无法退出"})
		sendGroupMessage(messageIn.GroupID, messageOut)
		return
	}
	if !bcp.IsAll && bcp.AttackBoss == 0 {
		if len(queuedBosses) > 1 {
			var s string
			for _, b := range queuedBosses {
				s += fmt.Sprintf("%d ", b)
			}
			messageOut.AddPart(&TextPart{Text: fmt.Sprintf("你加入了不止一个Boss(%s), 请指定想退出的boss编号或使用all退出所有队列", s)})
			sendGroupMessage(messageIn.GroupID, messageOut)
			return
		} else {
			// Quit the only queue
			err := database.QuitQueue(messageIn.GroupID, bcp.UserID, bcp.AttackBoss)
			if err != nil {
				sendErrorMessage(err)
				return
			}
		}
	}
	if bcp.AttackBoss != 0 {
		exists := false
		var s string
		for _, b := range queuedBosses {
			if b == bcp.AttackBoss {
				exists = true
			} else {
				s += fmt.Sprintf("%d ", b)
			}
		}
		if exists {
			// Quit the specific queue
			err := database.QuitQueue(messageIn.GroupID, bcp.UserID, bcp.AttackBoss)
			if err != nil {
				sendErrorMessage(err)
				return
			}
		} else {
			messageOut.AddPart(&TextPart{Text: fmt.Sprintf("你并未加入B%d, 请指定想退出的Boss编号(已加入:%s)或使用all退出所有队列", bcp.AttackBoss, s)})
			sendGroupMessage(messageIn.GroupID, messageOut)
			return
		}
	}
	showQueue(messageIn, bcp)
}

func reset(messageIn *GroupMessageIn, bcp *botCmdParams) {
	messageOut := &MessageBuilder{}
	if messageIn.Sender.Role == "owner" || messageIn.Sender.Role == "admin" {
		if bcp.IsAll {
			err := database.ResetAll(messageIn.GroupID)
			if err != nil {
				logrus.Errorf("Cannot reset entire database due to error:%v", err)
				messageOut.AddPart(&TextPart{Text: fmt.Sprintf("无法重置数据库, 错误：%s", err)})
			} else {
				logrus.Warn("Successfully reset entire database!")
				messageOut.AddPart(&TextPart{Text: "已重置全部数据"})
			}
		} else {
			err := database.ResetRecords(messageIn.GroupID)
			if err != nil {
				logrus.Errorf("Cannot reset database for next run due to error:%v", err)
				messageOut.AddPart(&TextPart{Text: fmt.Sprintf("无法重置数据库, 错误：%s", err)})
			} else {
				logrus.Warn("Successfully reset database for next run!")
				messageOut.AddPart(&TextPart{Text: "已重置伤害数据"})
			}
		}
	} else {
		messageOut.AddPart(&TextPart{Text: "仅群主或管理员可以执行重置"})
	}
}
