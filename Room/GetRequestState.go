package Room

import (
	"GoTest/HttpRequest"
)

type getRequestStateRequestBody struct {
	FanSpeed   string  `json:"fanSpeed"`
	TargetTemp float64 `json:"targetTemp"`
}

type getRequestStateResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    struct {
		RequestStatus string `json:"requestStatus"`
	} `json:"data"`
}

func GetRequestState() (string, error) {
	var requestBody getRequestStateRequestBody
	var response getRequestStateResponse
	_, responseStatus := HttpRequest.SendPostRequestWithToken("/room/poll/request", requestBody, &response)
	if responseStatus == 200 {
		return response.Data.RequestStatus, nil
	} else {
		return "Pending", nil
	}
}
