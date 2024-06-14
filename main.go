package main

import (
	"GoTest/Authentication"
	"GoTest/Room"
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var RefreshSpeed = 1
var room *Room.Room
var uiUpdate chan func()

// 定义绑定变量
var roomId binding.String
var temperature binding.Float
var targetTemperature binding.Float
var windSpeed binding.String
var workStatus binding.String
var ifWind binding.Bool
var envTempetaure binding.Float

func main() {
	// 初始化Fyne应用和窗口
	a := app.New()
	w := a.NewWindow("Air Conditioner Controller")

	// 初始化UI更新通道
	uiUpdate = make(chan func())

	// 设置登录界面为初始内容
	loginScreen := buildLoginScreen(w)
	w.SetContent(loginScreen)
	w.Resize(fyne.NewSize(600, 400))

	// 开启一个goroutine用于处理UI更新
	go func() {
		for update := range uiUpdate {
			update()
		}
	}()

	// 显示窗口并运行事件循环
	w.ShowAndRun()
}

func buildLoginScreen(w fyne.Window) fyne.CanvasObject {
	roomIdEntry := widget.NewEntry()
	roomIdEntry.SetPlaceHolder("Enter Room ID")

	identityEntry := widget.NewEntry()
	identityEntry.SetPlaceHolder("Enter Identity")

	loginButton := widget.NewButton("Login", func() {
		// 假设 Authentication.Login 函数返回一个 *Room.Room
		room = Authentication.Login(roomIdEntry.Text, identityEntry.Text)
		if room != nil {
			// 初始化绑定变量
			roomId = binding.BindString(&room.RoomId)
			temperature = binding.BindFloat(&room.Temperature)
			targetTemperature = binding.BindFloat(&room.TargetTemperature)
			windSpeed = binding.BindString(&room.WindSpeed)
			workStatus = binding.BindString(&room.Mode)
			ifWind = binding.BindBool(&room.IfWind)
			envTempetaure = binding.BindFloat(&room.EnvTemperature)

			uiUpdate <- func() {
				w.SetContent(buildMainScreen(w))
			}
			quit := make(chan struct{})
			go reportStatusPeriodically(room, quit)
			go checkTemperaturePeriodically(room, quit)
		} else {
			uiUpdate <- func() {
				dialog.ShowInformation("Login Failed", "Invalid Room ID or Identity", w)
			}
		}
	})

	loginForm := container.NewVBox(
		widget.NewLabelWithStyle("Air Conditioner Login", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		roomIdEntry,
		identityEntry,
		loginButton,
	)

	return container.NewCenter(container.NewVBox(
		loginForm,
	))
}

func buildMainScreen(w fyne.Window) fyne.CanvasObject {
	roomIdLabel := widget.NewLabelWithData(roomId)
	workStatusLabel := widget.NewLabelWithData(workStatus)
	temperatureLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(temperature, "%.2f"))
	windSpeedLabel := widget.NewLabelWithData(windSpeed)
	targetTemperatureLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(targetTemperature, "%.2f"))
	envTempetaureLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(envTempetaure, "%.2f"))
	ifWindLabel := widget.NewLabelWithData(binding.BoolToString(ifWind))

	targetTempEntry := widget.NewEntry()
	targetTempEntry.SetPlaceHolder("Enter Target Temperature")

	setTempButton := widget.NewButton("Set", func() {
		temp, err := strconv.ParseFloat(targetTempEntry.Text, 64)
		fmt.Println(room.Mode)
		if err == nil {
			if (room.Mode == "Cool" && temp >= 18 && temp <= 25) || (room.Mode == "Warm" && temp <= 30 && temp >= 25) {
				uiUpdate <- func() {
					targetTemperature.Set(temp)
				}
			} else {
				uiUpdate <- func() {
					dialog.ShowInformation("Invalid Target Temperature", "Target temperature should be 18~25 when cooling, 25~30 when warming", w)
				}
			}
		} else {
			uiUpdate <- func() {
				dialog.ShowError(err, w)
			}
		}
	})
	targetTempBox := container.NewHBox(widget.NewLabel("Set Target Temperature: "), container.New(layout.NewGridWrapLayout(fyne.NewSize(200, targetTempEntry.MinSize().Height)), targetTempEntry), setTempButton)

	windSpeedSelect := widget.NewSelect([]string{"Low", "Medium", "High"}, func(value string) {
		uiUpdate <- func() {
			windSpeed.Set(value)
		}
	})
	windSpeedBox := container.NewHBox(widget.NewLabel("Set Wind Speed: "), container.New(layout.NewGridWrapLayout(fyne.NewSize(200, windSpeedSelect.MinSize().Height)), windSpeedSelect))

	// 静态数据部分
	staticData := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Room ID", roomIdLabel),
			widget.NewFormItem("Work Status", workStatusLabel),
			widget.NewFormItem("Wind Speed", windSpeedLabel),
			widget.NewFormItem("If Wind", ifWindLabel),
			widget.NewFormItem("Env Temperature", envTempetaureLabel),
		),
	)
	staticDataBox := container.NewVBox(
		widget.NewCard("", "", staticData),
	)

	// 动态数据部分
	dynamicData := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Current Temperature", temperatureLabel),
			widget.NewFormItem("Target Temperature", targetTemperatureLabel),
		),
	)
	dynamicDataBox := container.NewVBox(
		widget.NewCard("", "", dynamicData),
	)

	logoutContainer := container.NewHBox(
		layout.NewSpacer(),
		widget.NewButtonWithIcon("", theme.LogoutIcon(), func() {
			Authentication.Logout()
			uiUpdate <- func() {
				w.SetContent(buildLoginScreen(w))
			}
		}),
	)

	controlPanel := container.NewVBox(
		targetTempBox,
		windSpeedBox,
	)

	return container.NewBorder(nil, logoutContainer, nil, nil,
		container.NewVBox(
			staticDataBox,
			dynamicDataBox,
			controlPanel,
		))
}

func reportStatusPeriodically(room *Room.Room, quit chan struct{}) {
	ticker1 := time.NewTicker(3 * time.Second / time.Duration(RefreshSpeed))
	defer ticker1.Stop()
	for {
		select {
		case <-ticker1.C:
			err, mode, refreshSpeed := Room.ReportStatus(room.Mode, room.Temperature)
			if err != nil {
				fmt.Println("ReportStatus error:", err)
			} else {
				if refreshSpeed != RefreshSpeed {
					RefreshSpeed = refreshSpeed
					ticker1.Stop()
					ticker1 = time.NewTicker(3 * time.Second / time.Duration(RefreshSpeed))
				}
				if room.Mode != mode {
					room.Mode = mode
				}
			}
		case <-quit:
			return
		}
	}
}

func checkTemperaturePeriodically(room *Room.Room, quit chan struct{}) {
	ticker2 := time.NewTicker(1 * time.Second)
	defer ticker2.Stop()
	for {
		select {
		case <-ticker2.C:
			// 检查当前有无请求
			if room.IfReq {
				// 当前有请求
				// 询问主控机，检查当前请求的送风情况
				room.IfWind = room.CheckIfWind()
				fmt.Println("Requesting. Now IfWind:", room.IfWind)

				if !room.IfWind {
					// 未送风
					// 等待主控机……
					fmt.Println("Remain request. Waiting for wind...")
				} else {
					// 正在送风
					// 检查当前温度达到目标登录。若是，则发送停止送风请求
					if room.CheckTemperatureCorrect() {
						fmt.Println("Temperature correct, stop wind")
						room.StopWind()
					}
				}
			} else {
				// 当前无请求
				fmt.Println("No request now")
				if room.CheckTemperatureOut() {
					// 一旦温度超出范围，就发送开始送风请求
					fmt.Println("Temperature out of range, start wind")
					room.StartWind()
				}
			}

			// 根据IfWind和WindSpeed，更新房间温度
			room.UpdateTemperature()

			// 更新UI
			uiUpdate <- func() {
				temperature.Set(room.Temperature)
				ifWind.Set(room.IfWind)
				windSpeed.Set(room.WindSpeed)
				workStatus.Set(room.Mode)
			}
		case <-quit:
			return
		}
	}
}
