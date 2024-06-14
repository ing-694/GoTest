// 用于开机时的认证登录
// 调用func Login()，会一直执行登录直至登录成功。返回房间号、空调模式、缺省温度

package Authentication

import (
	"GoTest/HttpRequest"
	"GoTest/Room"
)

type LoginRequestBody struct {
	RoomNumber string `json:"roomId"`
	Id         string `json:"identity"`
}

type LoginResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Mode               string  `json:"mode"`
		DefaultTemperature float64 `json:"defaultTemp"`
		RefreshRate        int     `json:"refreshRate"`
		Token              string  `json:"token"`
	} `json:"data"`
}

// 用于开机时的认证登录
// 使用得到的response信息，构造并返回新创建的room结构体
func Login(roomId string, identity string) *Room.Room {
	var room *Room.Room
	var data = LoginRequestBody{
		RoomNumber: roomId,
		Id:         identity,
	}
	var loginResp LoginResponse

	err, responseStatus := HttpRequest.SendPostRequestWithoutToken("/auth/login", data, &loginResp)
	if err != nil {
		return nil
	}

	if responseStatus == 200 {
		room = Room.NewRoom(roomId, loginResp.Data.Mode, loginResp.Data.DefaultTemperature)
		HttpRequest.Token = loginResp.Data.Token
	} else {
		return nil
	}
	return room
}
