package Room

import (
	"fmt"
	"time"
)

type Room struct {
	// 房间号
	RoomId string
	// 空调工作状态，有三种，warm/cold/standby
	WorkStatus string
	// 房间当前温度，开始时应该初始化为多少呢，这里我觉得应该是服务器先写好每个房间的初始温度，然后初始化的时候从服务器获取，和开机次数和关机次数一样，可能需要继承之前的值
	Temperature float64
	// 房间目标温度
	TargetTemperature float64
	// 房间风速，低中高分别对应1，2，3，关闭时为0
	WindSpeed string
}

func NewRoom(roomId string, mode string, targetTemperature float64) *Room {
	//新建房间时检查温度是否需要发起请求
	x := 28.0
	return &Room{
		RoomId:            roomId,
		WorkStatus:        mode,
		Temperature:       x, // TODO:可能要改为从服务器获取当前温度，这里先写死
		TargetTemperature: targetTemperature,
		WindSpeed:         "low",
	}
}

func CheckTemperature(room *Room) {
	//检查温度是否需要发起请求
	diff := room.Temperature - room.TargetTemperature
	if room.WorkStatus == "warm" && diff < -1 {
		//向服务器请求加热
		StartWind(room)
	} else if room.WorkStatus == "cool" && diff > 1 {
		//向服务器请求制冷
		StartWind(room)
	}
}

// 空调工作时的温度变化
func (room *Room) WorkingTemperatureChange(stop chan bool) {
	target := room.TargetTemperature
	var flag float64
	if room.WorkStatus == "warm" {
		flag = 1
	} else {
		flag = -1
	}
	var degreeLevel float64
	switch room.WindSpeed {
	case "low":
		degreeLevel = 0.5 * flag
		break
	case "medium":
		degreeLevel = 1 * flag
		break
	case "high":
		degreeLevel = 1.5
		break
	}

	// 温度每秒更新一次，每次变化speed度
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C:
				room.Temperature += degreeLevel
				fmt.Println("当前温度：", room.Temperature)
				// 检查是否达到目标温度
				if room.Temperature == target {
					fmt.Println("已达到目标温度")
					return
				}
			case <-stop:
				fmt.Println("主动停止温度变化")
				return
			}
		}
	}()
}
