package main

import (
	"fmt"
	modbusmodule "smart_streetlight/modbusmodel"
	"time"
)

func main() {
	// 创建 ModbusMaster 实例
	master, err := modbusmodule.NewModbusMaster("COM3", 9600)
	if err != nil {
		fmt.Println("创建 ModbusMaster 失败:", err)
		return
	}
	defer master.Close()

	// 创建温度传感器实例
	tempSensor, err := modbusmodule.CreateSensor(master, modbusmodule.TEMPERATURE)
	if err != nil {
		fmt.Println("创建温度传感器失败:", err)
		return
	}

	air, err := modbusmodule.CreateSensor(master, modbusmodule.AIR_QUALITY)
	if err != nil {
		fmt.Println("创建空气传感器实例失败", err)
		return
	}

	for {
		tempValue, err := tempSensor.Read(modbusmodule.Temperature)
		if err != nil {
			fmt.Println("读取温度失败:", err)
		} else {
			fmt.Println("温度值:", float64(tempValue)/10.0)
		}

		humidity, err1 := tempSensor.Read(modbusmodule.Humidity)
		if err1 != nil {
			fmt.Println("读取温度失败:", err)
		} else {
			fmt.Println("湿度值:", float64(humidity)/10.0)
		}

		pm25, err := air.Read(modbusmodule.PM2_5)
		if err != nil {
			fmt.Println("读取pm2.5失败:", err)
		} else {
			fmt.Println("pm2.5:", pm25)
		}

		pm10, err := air.Read(modbusmodule.PM10)
		if err != nil {
			fmt.Println("读取pm10失败:", err)
		} else {
			fmt.Println("pm10:", pm10)
		}

		time.Sleep(10000 * time.Millisecond)
	}
}
