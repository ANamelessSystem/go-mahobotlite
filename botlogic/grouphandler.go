package botlogic

import (
	"strconv"
	"strings"
)

func groupMessageHandler(message *GroupMessageIn) {
	atMe := "[CQ:at,qq=" + strconv.Itoa(message.SelfID) + "]"
	if !strings.Contains(message.Message, atMe) {
		return
	}
	// get username for update database, queue and other...
	// username := getUserName(message)

	// create string slice for raw message
	args := strings.Split(strings.TrimSpace(strings.ReplaceAll(message.Message, atMe, "")), " ")
	bcp := getCmdParams(args)

	switch strings.ToLower(args[0]) {
	case "c1":
		enqueue(message, &bcp)
	case "c2":
		// c2
	case "c3":
		// c3
	case "c4":
		// c4
	case "clear":
		// clear
	case "dmg":
		// dmg
	case "help":
		// help
	case "mod":
		// mod
	case "show":
		// show
	case "f1":
		// f1
	case "f2":
		// f2
	case "nla":
		enroll(message)
	case "nls":
		// nls
	case "nld":
		// nld
	case "sos":
		// sos
	case "reset":
		reset(message, &bcp)
	case "test":
		// testFunction(bcp)
	default:
		// unknown
	}
}

// func testFunction(data interface{}) error {
// 	// v := reflect.ValueOf(data)
// 	// t := v.Type()

// 	// for i := 0; i < v.NumField(); i++ {
// 	// 	fieldValue := v.Field(i)
// 	// 	fieldName := t.Field(i).Name
// 	// 	logrus.Debugf("Field Name: %s, Field Value: %v\n", fieldName, fieldValue)
// 	// }

// 	m := &MessageBuilder{}
// 	m.AddPart(&TextPart{Text: "收到"})

// 	mp := m.Build()

// 	cm := GroupMessageOut{
// 		GroupID:    877184755,
// 		Message:    mp,
// 		AutoEscape: false,
// 	}
// 	logrus.Debugf("Group message:%v,%v,%v", cm.AutoEscape, cm.GroupID, cm.Message)

// 	err := sendGroupMessage(cm)
// 	if err != nil {
// 		return err
// 	}
// 	logrus.Debugf("Send message with:%v", err)
// 	return nil
// }
