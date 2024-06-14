package Room

type Room struct {
	// 房间号
	RoomId string
	// 空调工作状态，
	Mode string
	// 房间当前温度
	Temperature float64
	// 房间目标温度
	TargetTemperature float64
	// 环境温度
	EnvTemperature float64
	// 房间风速，Low/Medium/High
	WindSpeed string
	// 当前是否被送风
	IfWind bool
	// 当前是否有请求
	IfReq bool
}

func NewRoom(roomId string, mode string, defaultTargetTemperature float64) *Room {
	return &Room{
		RoomId: roomId,
		Mode:   mode,
		Temperature: func() float64 {
			if mode == "Cool" {
				return 30
			} else {
				return 15
			}
		}(),
		TargetTemperature: defaultTargetTemperature,
		EnvTemperature: func() float64 {
			if mode == "Cool" {
				return 30
			} else {
				return 15
			}
		}(),
		WindSpeed: "Medium", // 默认风速为中风
		IfWind:    false,    // 默认不送风
		IfReq:     false,    // 默认无请求
	}
}

func (room *Room) CheckIfWind() bool {
	// 询问主控机，检查当前是否正在送风
	status, err := GetRequestState()
	// fmt.Println("Request status:", status)
	if err != nil {
		return false
	} else {
		if status == "Doing" {
			return true
		} else {
			return false
		}
	}
}

// 检查温度是否已超出目标温度范围
func (room *Room) CheckTemperatureOut() bool {
	diff := room.Temperature - room.TargetTemperature
	if (room.Mode == "Cool" && diff > 1) || (room.Mode == "Heat" && diff < -1) {
		return true
	} else {
		return false
	}
}

// 检查温度是否已达到目标温度
func (room *Room) CheckTemperatureCorrect() bool {
	if (room.Mode == "Cool" && room.Temperature <= room.TargetTemperature) || (room.Mode == "Heat" && room.Temperature >= room.TargetTemperature) {
		return true
	} else {
		return false
	}
}

// 根据当前IfWind和WindSpeed，更新温度
func (room *Room) UpdateTemperature() {
	if !room.IfWind {
		// 当前没有得到送风时，回归环境温度，每秒变化0.3度
		if room.Mode == "Cool" && room.Temperature < room.EnvTemperature {
			room.Temperature += 0.2
			if room.Temperature > room.EnvTemperature {
				room.Temperature = room.EnvTemperature
			}
		} else if room.Mode == "Warm" && room.Temperature > room.EnvTemperature {
			room.Temperature -= 0.2
			if room.Temperature < room.EnvTemperature {
				room.Temperature = room.EnvTemperature
			}
		}
	} else {
		// 当前正在送风时，按照温度变化曲线，模拟变化温度。在中风速下，每秒变化0.5度
		if room.WindSpeed == "Low" {
			if room.Mode == "Cool" {
				room.Temperature -= 0.3
			} else {
				room.Temperature += 0.3
			}
		} else if room.WindSpeed == "Medium" {
			if room.Mode == "Cool" {
				room.Temperature -= 0.5
			} else {
				room.Temperature += 0.5
			}
		} else if room.WindSpeed == "High" {
			if room.Mode == "Cool" {
				room.Temperature -= 0.7
			} else {
				room.Temperature += 0.7
			}
		}
	}
}
