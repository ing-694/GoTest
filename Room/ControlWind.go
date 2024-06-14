package Room

import (
	"GoTest/HttpRequest"
	"fmt"
)

type StartWindRequestBody struct {
	FanSpeed   string  `json:"fanSpeed"`
	TargetTemp float64 `json:"targetTemp"`
}

// 发送送风请求
type StartWindResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func (room *Room) StartWind() error {
	//发送送风请求
	requestBody := StartWindRequestBody{
		FanSpeed:   room.WindSpeed,
		TargetTemp: room.TargetTemperature,
	}
	var response StartWindResponse
	err, responseStatus := HttpRequest.SendPostRequestWithToken("/room/blowing/start", requestBody, &response)
	if err != nil {
		fmt.Println("送风请求发送错误：", err)
		return err
	}
	if responseStatus == 200 {
		fmt.Println("Send start wind request successfully")
		room.IfReq = true
	} else {
		fmt.Println("送风请求失败：", response.Message)
	}
	return nil
}

func (room *Room) StopWind() error {
	var requestBody map[string]interface{}
	var response map[string]interface{}
	err, responseStatus := HttpRequest.SendPostRequestWithToken("/room/blowing/stop", requestBody, &response)
	if err != nil {
		fmt.Println("停止送风请求发送错误：", err)
		return err
	}
	if responseStatus == 200 {
		fmt.Println("Send stop wind request successfully")
		room.IfReq = false
		room.IfWind = false
		return nil
	} else {
		fmt.Println("停止送风请求失败：", response["message"])
	}
	return nil
}
