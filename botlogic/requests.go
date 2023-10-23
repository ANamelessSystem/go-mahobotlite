package botlogic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

var cqHttpAddr string

func SetCQAddr(addr string) {
	cqHttpAddr = addr
}

func sendGroupMessage(grpID int, messageBuilder *MessageBuilder) {
	// handle all error inside this function
	mp := messageBuilder.Build()
	gmo := GroupMessageOut{
		GroupID:    grpID,
		Message:    mp,
		AutoEscape: false,
	}
	resp, err := sendRequest("send_group_msg", "POST", gmo)
	if err != nil {
		logrus.Errorf("failed to POST message to CQ Server: %v", err)
		return
	}
	var cqServerResponse CQServerResponse
	err = json.Unmarshal(resp, &cqServerResponse)
	if err != nil {
		logrus.Errorf("failed to unmarshal response from CQ Server: %v", err)
		return
	}
	if cqServerResponse.Status == "failed" {
		logrus.Errorf("cq server accepted message but response an error: %v", cqServerResponse)
		return
	}
}

func sendRequest(endpoint string, method string, payload interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", cqHttpAddr, endpoint)

	var body []byte
	if payload != nil {
		var errMarshal error
		body, errMarshal = json.Marshal(payload)
		if errMarshal != nil {
			return nil, fmt.Errorf("failed to marshal payload: %v", errMarshal)
		}
	}

	req, errReq := http.NewRequest(method, url, bytes.NewBuffer(body))
	if errReq != nil {
		return nil, fmt.Errorf("failed to create request: %v", errReq)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	logrus.Debugf("Sending JSON payload(unmarshal): %s", body)
	logrus.Debugf("Request sent %v", req)

	resp, errResp := http.DefaultClient.Do(req)
	if errResp != nil {
		return nil, fmt.Errorf("failed to send request: %v", errResp)
	}
	defer resp.Body.Close()

	respBody, errBody := io.ReadAll(resp.Body)
	if errBody != nil {
		return nil, fmt.Errorf("failed to read response body: %v", errBody)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d - %s", resp.StatusCode, respBody)
	}

	return respBody, nil
}

func GetLoginInfo() (*LoginInfoResponse, error) {
	body, err := sendRequest("get_login_info", "GET", nil)
	if err != nil {
		return nil, err
	}

	var loginInfo LoginInfoResponse
	err = json.Unmarshal(body, &loginInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return &loginInfo, nil
}
