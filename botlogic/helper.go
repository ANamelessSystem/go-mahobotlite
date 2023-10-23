package botlogic

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

func getCmdParams(args []string) botCmdParams {
	var bcp botCmdParams
	reEventID := regexp.MustCompile(`^(?:[eE](\d+)|(\d+)[eE])$`)
	reUserID1 := regexp.MustCompile(`^(?:[uU](\d+)|(\d+)[uU])$`)
	reUserID2 := regexp.MustCompile(`^\[CQ:at,qq=(\d+)\]$`)
	reAttackBoss := regexp.MustCompile(`^(?:[bB](\d+)|(\d+)[bB])$`)
	reAttackRound := regexp.MustCompile(`^(?:[rR](\d+)|(\d+)[rR]|周目(\d+)|(\d+)周目)$`)
	reAttackDamage := regexp.MustCompile(`^(\d+)(\.\d+)?([kKmMwW万千亿])?$`)

	for _, arg := range args {
		lowerArg := strings.ToLower(arg)
		switch {
		case lowerArg == "all":
			bcp.IsAll = true
		case reEventID.MatchString(arg):
			matches := reEventID.FindStringSubmatch(arg)
			for _, match := range matches[1:] {
				if match != "" {
					bcp.EventID, _ = strconv.Atoi(match)
					break
				}
			}
		case reUserID1.MatchString(arg):
			matches := reUserID1.FindStringSubmatch(arg)
			for _, match := range matches[1:] {
				if match != "" {
					bcp.UserID, _ = strconv.Atoi(match)
					break
				}
			}
		case reUserID2.MatchString(arg):
			matches := reUserID2.FindStringSubmatch(arg)
			bcp.UserID, _ = strconv.Atoi(matches[1])
		case reAttackBoss.MatchString(arg):
			matches := reAttackBoss.FindStringSubmatch(arg)
			for _, match := range matches[1:] {
				if match != "" {
					bcp.AttackBoss, _ = strconv.Atoi(match)
					break
				}
			}
		case reAttackRound.MatchString(arg):
			matches := reAttackRound.FindStringSubmatch(arg)
			for _, match := range matches[1:] {
				if match != "" {
					bcp.AttackRound, _ = strconv.Atoi(match)
					break
				}
			}
		case reAttackDamage.MatchString(arg):
			bcp.AttackDamage = getDamageNumber(arg)
		case lowerArg == "last", arg == "尾刀":
			bcp.AttackType = 2
		case lowerArg == "ext", arg == "补时":
			bcp.AttackType = 1
		case lowerArg == "normal", arg == "通常":
			bcp.AttackType = 0
		case lowerArg == "timeout", arg == "掉线":
			bcp.AttackLost = true
		}
	}
	return bcp
}

func getDamageNumber(s string) int {
	re := regexp.MustCompile(`^(\d+)(\.\d+)?([kKmMwW万千亿])?$`)
	matches := re.FindStringSubmatch(s)

	num, err := strconv.ParseFloat(matches[1]+matches[2], 64)
	if err != nil {
		return 0
	}

	switch matches[3] {
	case "k", "K", "千":
		num *= 1_000
	case "m", "M":
		num *= 1_000_000
	case "w", "万":
		num *= 10_000
	case "亿":
		num *= 100_000_000
	}

	return int(num)
}

func getUserName(message *GroupMessageIn) string {
	username := ""
	if message.Sender.Card != "" {
		username = replaceInvalidRunes(message.Sender.Card)
	} else if message.Sender.Nickname != "" {
		username = replaceInvalidRunes(message.Sender.Nickname)
	} else {
		username = strconv.Itoa(message.Sender.UserID)
	}
	return username
}

func replaceInvalidRunes(s string) string {
	result := ""
	for _, r := range s {
		if utf8.ValidRune(r) {
			result += string(r)
		} else {
			result += "*"
		}
	}
	return result
}
