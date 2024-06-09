package Room

import (
	"GoTest/HttpRequest"
	"fmt"
)

type ReportStatusRequestBody struct {
	Status      string  `json:"status"`
	Temperature float64 `json:"temperature"`
}
type ReportStatusResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    struct {
		WorkingStatus    string `json:"mode"`
		RefreshRate      int    `json:"refreshRate"`
		Daily_statistics struct {
			Energy float64 `json:"energy"`
			Cost   float64 `json:"cost"`
		} `json:"daily_statistics"`
	} `json:"data"`
}

// 向主控机汇报从控机的状态
func ReportStatus(status string, temperature float64) (error, string, int) {
	requestBody := ReportStatusRequestBody{
		Status:      status,
		Temperature: temperature,
	}
	var response ReportStatusResponse
	err, responseStatus := HttpRequest.SendPostRequestWithToken("/room/poll/room_status", requestBody, &response)
	if err != nil {
		fmt.Println("上报状态请求发送错误：", err)
		return err, "错误", -1
	} else if responseStatus == 200 {
		fmt.Printf("当前房间内温度:%.1f；当日使用金额：%.1f\n", temperature, response.Data.Daily_statistics.Cost)
	} else {
		fmt.Println("上报状态请求失败：", response.Message)
	}
	return err, response.Data.WorkingStatus, response.Data.RefreshRate
}
