package botlogic

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

// 处理所有接收到的信息
func ReceiveHandler(w http.ResponseWriter, r *http.Request) {
	body, errHttp := io.ReadAll(r.Body)
	if errHttp != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var data map[string]interface{}
	errUnmarshal := json.Unmarshal(body, &data)
	if errUnmarshal != nil {
		logrus.Warnf("Failed to unmarshal body: %v", errUnmarshal)
	}

	switch data["post_type"].(string) {
	case "meta_event":
		metaEventHandler(&body)
	case "message":
		messageHandler(&body)
	}

	w.Write([]byte("Received data successfully"))
}

func messageHandler(body *[]byte) {
	logrus.Debugf("Received data from middleware:%s", body)
	var message GroupMessageIn
	errMessage := json.Unmarshal(*body, &message)
	if errMessage != nil {
		logrus.Warnf("Failed to unmarshal message: %v", errMessage)
		return
	}
	logrus.Infof("Message received: %s", *body)
	switch message.MessageType {
	case "group":
		groupMessageHandler(&message)
	}
	// continue to handle message event as need
}

func metaEventHandler(body *[]byte) {
	var metaEvent MetaEvent
	errUnmarshal := json.Unmarshal(*body, &metaEvent)
	if errUnmarshal != nil {
		logrus.Warnf("Failed to unmarshal meta event: %v", errUnmarshal)
		return
	}
	switch metaEvent.MetaEventType {
	case "heartbeat":
		heartbeatHandler()
	}
	// continue to handle meta event as need
}
