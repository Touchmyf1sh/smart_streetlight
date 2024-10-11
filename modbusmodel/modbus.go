package modbusmodule

import (
	"encoding/binary"
	"fmt"
	"github.com/goburrow/modbus"
	"sync"
	"time"
)

const (
	PM2_5       = 0x0000 // pm2.5，实际值 = 读取值， 单位ug/m3
	PM10        = 0x0001 // pm10，实际值 = 读取值， 单位ug/m3
	Light       = 0x0002 // 光照强度，实际值 = 读取值 / 1000，单位Lux
	Temperature = 0x0000 // 温度，实际值 = 读取值/10，单位摄氏度
	Humidity    = 0x0001 // 空气湿度,实际值 = 读取值/10，单位为%
	Rain        = 0x0000 // 雨量，实际值 = 读取值/10，单位mm
	Wind        = 0x0000 // 风速，实际值 = 读取值/10，单位m/s
	Ultra_light = 0x0000 // 紫外线，实际值 = 读取值/10，单位uw/cm2
)

// 设备名常量
const (
	RAIN        = 0x01 // 雨量传感器
	ULTRA_LIGHT = 0x02 // 紫外线传感器
	WIND        = 0x03 // 风速传感器
	AIR_QUALITY = 0x04 // 空气质量传感器
	LIGHT       = 0x05 // 光照传感器
	TEMPERATURE = 0x06 // 温、湿度传感器
)

type ModbusDevice interface {
	Read(address uint16) (int32, error)
	Write(address uint16, value int16) error
}

// 通用的 Modbus 基础结构体
type ModbusMaster struct {
	handler *modbus.RTUClientHandler
	client  modbus.Client
	mu      sync.Mutex // 新增的互斥锁
}

// 温、湿度传感器
type TemperatureSensor struct {
	master     *ModbusMaster
	slaveID    byte
	maxRetries int
	retryDelay time.Duration
}

func (t *TemperatureSensor) Read(address uint16) (int32, error) {
	var results []byte
	var err error

	t.master.mu.Lock()
	t.master.handler.SlaveId = t.slaveID
	t.master.mu.Unlock()
	for attempt := 1; attempt <= t.maxRetries; attempt++ {

		// 功能码 0x03
		results, err = t.master.client.ReadHoldingRegisters(address, 1)
		if err == nil {
			if len(results) < 2 {
				return 0, fmt.Errorf("读取的数据长度不足")
			}
			data := binary.BigEndian.Uint16(results)
			value := int32(data & 0x7FFF)
			if (data & 0x8000) == 0x8000 {
				value = -value
			}
			return value, nil
		}
		if attempt == t.maxRetries {
			return 0, fmt.Errorf("读取失败，已重试 %d 次：%v", t.maxRetries, err)
		}
		time.Sleep(t.retryDelay)
	}
	return 0, fmt.Errorf("未知错误")
}

func (t *TemperatureSensor) Write(address uint16, value int16) error {
	var err error
	data := uint16(value)

	t.master.mu.Lock()
	t.master.handler.SlaveId = t.slaveID
	t.master.mu.Unlock()
	for attempt := 1; attempt <= t.maxRetries; attempt++ {

		// 功能码 0x06
		_, err = t.master.client.WriteSingleRegister(address, data)
		if err == nil {
			return nil
		}
		if attempt == t.maxRetries {
			return fmt.Errorf("写入失败，已重试 %d 次：%v", t.maxRetries, err)
		}
		time.Sleep(t.retryDelay)
	}
	return fmt.Errorf("未知错误")
}

//雨量传感器

type RainSensor struct {
	master     *ModbusMaster
	slaveID    byte
	maxRetries int
	retryDelay time.Duration
}

func (t *RainSensor) Read(address uint16) (int32, error) {
	var results []byte
	var err error

	t.master.mu.Lock()
	t.master.handler.SlaveId = t.slaveID
	t.master.mu.Unlock()
	for attempt := 1; attempt <= t.maxRetries; attempt++ {

		// 功能码 0x03
		results, err = t.master.client.ReadHoldingRegisters(address, 1)
		if err == nil {
			if len(results) < 2 {
				return 0, fmt.Errorf("读取的数据长度不足")
			}
			data := binary.BigEndian.Uint16(results)
			value := int32(data)
			return value, nil
		}
		if attempt == t.maxRetries {
			return 0, fmt.Errorf("读取失败，已重试 %d 次：%v", t.maxRetries, err)
		}
		time.Sleep(t.retryDelay)
	}
	return 0, fmt.Errorf("未知错误")
}

func (t *RainSensor) Write(address uint16, value int16) error {
	var err error
	data := uint16(value)

	t.master.mu.Lock()
	t.master.handler.SlaveId = t.slaveID
	t.master.mu.Unlock()
	for attempt := 1; attempt <= t.maxRetries; attempt++ {

		// 功能码 0x06
		_, err = t.master.client.WriteSingleRegister(address, data)
		if err == nil {
			return nil
		}
		if attempt == t.maxRetries {
			return fmt.Errorf("写入失败，已重试 %d 次：%v", t.maxRetries, err)
		}
		time.Sleep(t.retryDelay)
	}
	return fmt.Errorf("未知错误")
}

// 光照传感器
type LightSensor struct {
	master     *ModbusMaster
	slaveID    byte
	maxRetries int
	retryDelay time.Duration
}

func (t *LightSensor) Read(address uint16) (int32, error) {
	var results []byte
	var err error

	t.master.mu.Lock()
	t.master.handler.SlaveId = t.slaveID
	t.master.mu.Unlock()
	for attempt := 1; attempt <= t.maxRetries; attempt++ {

		// 功能码 0x03
		results, err = t.master.client.ReadHoldingRegisters(address, 2)
		if err == nil {
			if len(results) < 2 {
				return 0, fmt.Errorf("读取的数据长度不足")
			}
			data := binary.BigEndian.Uint32(results)
			value := int32(data)
			return value, nil
		}
		if attempt == t.maxRetries {
			return 0, fmt.Errorf("读取失败，已重试 %d 次：%v", t.maxRetries, err)
		}
		time.Sleep(t.retryDelay)
	}
	return 0, fmt.Errorf("未知错误")
}

func (t *LightSensor) Write(address uint16, value int16) error {
	var err error
	data := uint16(value)

	t.master.mu.Lock()
	t.master.handler.SlaveId = t.slaveID
	t.master.mu.Unlock()
	for attempt := 1; attempt <= t.maxRetries; attempt++ {

		// 功能码 0x06
		_, err = t.master.client.WriteSingleRegister(address, data)
		if err == nil {
			return nil
		}
		if attempt == t.maxRetries {
			return fmt.Errorf("写入失败，已重试 %d 次：%v", t.maxRetries, err)
		}
		time.Sleep(t.retryDelay)
	}
	return fmt.Errorf("未知错误")
}

//风速传感器

type WindCensor struct {
	master     *ModbusMaster
	slaveID    byte
	maxRetries int
	retryDelay time.Duration
}

func (t *WindCensor) Read(address uint16) (int32, error) {
	var results []byte
	var err error

	t.master.mu.Lock()
	t.master.handler.SlaveId = t.slaveID
	t.master.mu.Unlock()

	for attempt := 1; attempt <= t.maxRetries; attempt++ {

		// 功能码 0x03
		results, err = t.master.client.ReadHoldingRegisters(address, 1)
		if err == nil {
			if len(results) < 2 {
				return 0, fmt.Errorf("读取的数据长度不足")
			}
			data := binary.BigEndian.Uint16(results)
			value := int32(data)
			return value, nil
		}
		if attempt == t.maxRetries {
			return 0, fmt.Errorf("读取失败，已重试 %d 次：%v", t.maxRetries, err)
		}
		time.Sleep(t.retryDelay)
	}
	return 0, fmt.Errorf("未知错误")
}

func (t *WindCensor) Write(address uint16, value int16) error {
	var err error
	data := uint16(value)

	t.master.mu.Lock()
	t.master.handler.SlaveId = t.slaveID
	t.master.mu.Unlock()
	for attempt := 1; attempt <= t.maxRetries; attempt++ {

		// 功能码 0x06
		_, err = t.master.client.WriteSingleRegister(address, data)
		if err == nil {
			return nil
		}
		if attempt == t.maxRetries {
			return fmt.Errorf("写入失败，已重试 %d 次：%v", t.maxRetries, err)
		}
		time.Sleep(t.retryDelay)
	}
	return fmt.Errorf("未知错误")
}

// 紫外线传感器
type UltralightSensor struct {
	master     *ModbusMaster
	slaveID    byte
	maxRetries int
	retryDelay time.Duration
}

func (t *UltralightSensor) Read(address uint16) (int32, error) {
	var results []byte
	var err error

	t.master.mu.Lock()
	t.master.handler.SlaveId = t.slaveID
	t.master.mu.Unlock()
	for attempt := 1; attempt <= t.maxRetries; attempt++ {

		// 功能码 0x03
		results, err = t.master.client.ReadHoldingRegisters(address, 1)
		if err == nil {
			if len(results) < 2 {
				return 0, fmt.Errorf("读取的数据长度不足")
			}
			data := binary.BigEndian.Uint16(results)
			value := int32(data)
			return value, nil
		}
		if attempt == t.maxRetries {
			return 0, fmt.Errorf("读取失败，已重试 %d 次：%v", t.maxRetries, err)
		}
		time.Sleep(t.retryDelay)
	}
	return 0, fmt.Errorf("未知错误")
}

func (t *UltralightSensor) Write(address uint16, value int16) error {
	var err error
	data := uint16(value)

	t.master.mu.Lock()
	t.master.handler.SlaveId = t.slaveID
	t.master.mu.Unlock()
	for attempt := 1; attempt <= t.maxRetries; attempt++ {

		// 功能码 0x06
		_, err = t.master.client.WriteSingleRegister(address, data)
		if err == nil {
			return nil
		}
		if attempt == t.maxRetries {
			return fmt.Errorf("写入失败，已重试 %d 次：%v", t.maxRetries, err)
		}
		time.Sleep(t.retryDelay)
	}
	return fmt.Errorf("未知错误")
}

// 空气传感器
type AirQualitySensor struct {
	master     *ModbusMaster
	slaveID    byte
	maxRetries int
	retryDelay time.Duration
}

func (a *AirQualitySensor) Read(address uint16) (int32, error) {
	var results []byte
	var err error

	a.master.mu.Lock()
	a.master.handler.SlaveId = a.slaveID
	a.master.mu.Unlock()
	for attempt := 1; attempt <= a.maxRetries; attempt++ {

		// 功能码 0x03
		results, err = a.master.client.ReadInputRegisters(address, 1)
		if err == nil {
			if len(results) < 2 {
				return 0, fmt.Errorf("读取的数据长度不足")
			}
			data := binary.BigEndian.Uint16(results)
			value := int32(data)
			return value, nil
		}
		if attempt == a.maxRetries {
			return 0, fmt.Errorf("读取失败，已重试 %d 次：%v", a.maxRetries, err)
		}
		time.Sleep(a.retryDelay)
	}
	return 0, fmt.Errorf("未知错误")
}

func (t *AirQualitySensor) Write(address uint16, value int16) error {
	var err error
	data := uint16(value)

	t.master.mu.Lock()
	t.master.handler.SlaveId = t.slaveID
	t.master.mu.Unlock()
	for attempt := 1; attempt <= t.maxRetries; attempt++ {

		// 功能码 0x06
		_, err = t.master.client.WriteSingleRegister(address, data)
		if err == nil {
			return nil
		}
		if attempt == t.maxRetries {
			return fmt.Errorf("写入失败，已重试 %d 次：%v", t.maxRetries, err)
		}
		time.Sleep(t.retryDelay)
	}
	return fmt.Errorf("未知错误")
}

// 创建 Modbus 设备实例
func NewModbusMaster(port string, baudRate int) (*ModbusMaster, error) {
	handler := modbus.NewRTUClientHandler(port)
	handler.BaudRate = baudRate
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.Timeout = 1 * time.Second

	err := handler.Connect()
	if err != nil {
		return nil, fmt.Errorf("无法连接到串口: %v", err)
	}

	client := modbus.NewClient(handler)

	return &ModbusMaster{
		handler: handler,
		client:  client,
	}, nil
}

func (m *ModbusMaster) Close() error {
	return m.handler.Close()
}

func CreateSensor(master *ModbusMaster, deviceName byte) (ModbusDevice, error) {
	switch deviceName {
	case RAIN:
		return &RainSensor{
			master:     master,
			slaveID:    RAIN,
			maxRetries: 3,
			retryDelay: 100 * time.Millisecond,
		}, nil
	case ULTRA_LIGHT:
		return &UltralightSensor{
			master:     master,
			slaveID:    ULTRA_LIGHT,
			maxRetries: 3,
			retryDelay: 100 * time.Millisecond,
		}, nil
	case WIND:
		return &WindCensor{
			master:     master,
			slaveID:    WIND,
			maxRetries: 3,
			retryDelay: 100 * time.Millisecond,
		}, nil
	case AIR_QUALITY:
		return &AirQualitySensor{
			master:     master,
			slaveID:    AIR_QUALITY,
			maxRetries: 3,
			retryDelay: 100 * time.Millisecond,
		}, nil
	case LIGHT:
		return &LightSensor{
			master:     master,
			slaveID:    LIGHT,
			maxRetries: 3,
			retryDelay: 100 * time.Millisecond,
		}, nil
	case TEMPERATURE:
		return &TemperatureSensor{
			master:     master,
			slaveID:    TEMPERATURE,
			maxRetries: 3,
			retryDelay: 100 * time.Millisecond,
		}, nil
	default:
		return nil, fmt.Errorf("未知的设备名：0x%X", deviceName)
	}
}
