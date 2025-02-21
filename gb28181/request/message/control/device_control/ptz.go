package devicecontrol

import (
	"context"
	"errors"
	"fmt"
	"goutil/gb28181/request"
	"goutil/gb28181/xml"
	"goutil/sip"
)

// 云台控制命令
// ptz 类型命令，0: 水平速度（0-255）；1: 垂直速度（0-255）；2: 变倍速度（0-15）
// fi 类型命令，0: 聚焦速度（0-15）；1: 光圈速度（0-15）；2: 无用
// preset 类型命令，0: 无用；1: 预置位编号（1-255）；2: 无用
// cruise 类型命令，0: 巡航组号（0-255）；1: 预置位编号（1-255）；2: 数值高 4 位（0-4095，速度和时间才有效）
// scan 类型命令，0: 扫描组号（0-255）；1: 无用；2: 无用
// assist 类型命令，0: 辅助开关编号（0-255）
const (
	// 镜头变倍，缩小
	CmdPTZZoomOut = "ptzZoomOut"
	// 镜头变倍，放大
	CmdPTZZoomIn = "ptzZoomIn"
	// 上
	CmdPTZUp = "ptzUp"
	// 下
	CmdPTZDownn = "ptzDown"
	// 左
	CmdPTZLeft = "ptzLeft"
	// 左上
	CmdPTZLeftUp = "ptzLeftUp"
	// 左下
	CmdPTZLeftDown = "ptzLeftDown"
	// 右
	CmdPTZRight = "ptzRight"
	// 右上
	CmdPTZRightUp = "ptzRightUp"
	// 右下
	CmdPTZRightDown = "ptzRightDown"
	// 停止
	CmdPTZStop = "ptzStop"
	// 聚焦，近
	CmdFIFocusNear = "fiFocusNear"
	// 聚焦，远
	CmdFIFocusFar = "fiFocusFar"
	// 光圈，缩小
	CmdFIIrisOut = "fiIrisOut"
	// 光圈，缩小，聚焦，远
	CmdFIIrisOutFocusFar = "fiIrisOutFocusFar"
	// 光圈，缩小，聚焦，近
	CmdFIIrisOutFocusNear = "fiIrisOutFocusNear"
	// 光圈，放大
	CmdFIIrisIn = "fiIrisIn"
	// 光圈，放大，聚焦，远
	CmdFIIrisInFocusFar = "fiIrisInFocusFar"
	// 光圈，放大，聚焦，近
	CmdFIIrisInFocusNear = "fiIrisInFocusNear"
	// 停止
	CmdFIStop = "fiStop"
	// 设置预置位
	CmdPresetSet = "presetSet"
	// 调用预置位
	CmdPresetCall = "presetCall"
	// 删除预置位
	CmdPresetDelete = "presetDelete"
	// 添加巡航点
	CmdCruiseAdd = "cruiseAdd"
	// 删除巡航点
	CmdCruiseDelete = "cruiseDelete"
	// 设置巡航速度
	CmdCruiseSpeed = "cruiseSpeed"
	// 设置巡航停留时间
	CmdCruiseStay = "cruiseStay"
	// 开始巡航
	CmdCruiseStart = "cruiseStart"
	// 停止巡航
	CmdCruiseStop = "cruiseStop"
	// 开始扫描
	CmdScanStart = "scanStart"
	// 扫描左边界
	CmdScanLeft = "scanLeft"
	// 扫描右边界
	CmdScanRight = "scanRight"
	// 扫描速度
	CmdScanSpeed = "scanSpeed"
	// 停止扫描
	CmdScanStop = "scanStop"
	// 开启辅助开关
	CmdAssistStart = "assistStart"
	// 关闭辅助开关
	CmdAssistStop = "assistStop"
)

var (
	errCmd = errors.New("error command")
)

// PTZScanSpeed 返回信令中的两个 byte
func PTZScanSpeed(speed int16) (byte, byte) {
	return byte(speed >> 4), byte(speed) & 0x0F
}

// PTZScanSpeedValue 返回从信令中的两个 byte 得到的数值
func PTZScanSpeedValue(v6, v7 byte) int16 {
	return int16(v6)<<4 | int16(v7&0x0F)
}

// PTZ 是 SendPTZ 的参数
type PTZ struct {
	// 命令
	Command string
	// 优先级
	ControlPriority string
	// 原始指令
	RawCmd string
	// 数值
	V1        byte
	V2        byte
	V3        byte
	Ser       *sip.Server
	Device    request.Request
	ChannelID string
	// 追踪标识
	TraceID string
}

// Cmd 返回 cmd 字符串
func (m *PTZ) Cmd() (string, error) {
	cmdCode := 0
	var v1, v2, v3 byte
	switch m.Command {
	case CmdPTZZoomOut:
		// 变倍，缩小
		cmdCode = 0x20
		v3 = m.V3 & 0x0F
	case CmdPTZZoomIn:
		// 变倍，放大
		cmdCode = 0x10
		v3 = m.V3 & 0x0F
	case CmdPTZUp:
		// 方向，上
		cmdCode = 0x08
		v2 = m.V2
	case CmdPTZDownn:
		// 方向，下
		cmdCode = 0x04
		v2 = m.V2
	case CmdPTZLeft:
		// 方向，左
		cmdCode = 0x02
		v1 = m.V1
	case CmdPTZRight:
		// 方向，右
		cmdCode = 0x01
		v1 = m.V1
	case CmdPTZLeftUp:
		// 方向，左上
		cmdCode = 0x0a
		v1 = m.V1
		v2 = m.V2
	case CmdPTZLeftDown:
		// 方向，左下
		cmdCode = 0x06
		v1 = m.V1
		v2 = m.V2
	case CmdPTZRightUp:
		// 方向，右上
		cmdCode = 0x09
		v1 = m.V1
		v2 = m.V2
	case CmdPTZRightDown:
		// 方向，右下
		cmdCode = 0x05
		v1 = m.V1
		v2 = m.V2
	case CmdPTZStop:
		// 停止 ptz
		cmdCode = 0x00
	case CmdFIIrisOut:
		// 光圈，缩小
		cmdCode = 0x48
		v2 = m.V2
	case CmdFIIrisIn:
		// 光圈，放大
		cmdCode = 0x44
		v2 = m.V2
	case CmdFIFocusNear:
		// 聚焦，近
		cmdCode = 0x42
		v1 = m.V1
	case CmdFIFocusFar:
		// 聚焦，远
		cmdCode = 0x41
		v1 = m.V1
	case CmdFIIrisOutFocusFar:
		// 光圈，缩小，聚焦，远
		cmdCode = 0x49
		v1 = m.V1
		v2 = m.V2
	case CmdFIIrisOutFocusNear:
		// 光圈，缩小，聚焦，近
		cmdCode = 0x4a
		v1 = m.V1
		v2 = m.V2
	case CmdFIIrisInFocusFar:
		// 光圈，放大，聚焦，远
		cmdCode = 0x45
		v1 = m.V1
		v2 = m.V2
	case CmdFIIrisInFocusNear:
		// 光圈，放大，聚焦，近
		cmdCode = 0x46
		v1 = m.V1
		v2 = m.V2
	case CmdFIStop:
		// 停止 fi
		cmdCode = 0x40
	case CmdPresetSet:
		// 设置预置位
		cmdCode = 0x81
		v2 = m.V2
	case CmdPresetCall:
		// 调用预置位
		cmdCode = 0x82
		v2 = m.V2
	case CmdPresetDelete:
		// 删除预置位
		cmdCode = 0x83
		v2 = m.V2
	case CmdCruiseAdd:
		// 添加巡航点
		cmdCode = 0x84
		v1 = m.V1
		v2 = m.V2
	case CmdCruiseDelete:
		// 删除巡航点
		cmdCode = 0x85
		v1 = m.V1
		v2 = m.V2
	case CmdCruiseSpeed:
		// 设置巡航速度
		cmdCode = 0x86
		v1 = m.V1
		v2 = m.V2
		v3 = m.V3 & 0x0F
	case CmdCruiseStay:
		// 设置巡航停留时间
		cmdCode = 0x87
		v1 = m.V1
		v2 = m.V2
		v3 = m.V3 & 0x0F
	case CmdCruiseStart:
		// 开始巡航
		cmdCode = 0x88
		v1 = m.V1
	case CmdCruiseStop:
		// 停止巡航
		cmdCode = 0x00
	case CmdScanStart:
		// 开始自动扫描
		cmdCode = 0x89
		v1 = m.V1
		v2 = 0x00
	case CmdScanLeft:
		// 设置自动扫描左边界
		cmdCode = 0x89
		v1 = m.V1
		v2 = 0x01
	case CmdScanRight:
		// 设置自动扫描右边界
		cmdCode = 0x89
		v1 = m.V1
		v2 = 0x02
	case CmdScanSpeed:
		// 设置自动扫描速度
		// v2 数据的低 8 位，v3 数据的高 4 位
		cmdCode = 0x8A
		v1 = m.V1
		v2 = m.V2 >> 4
		v3 = m.V3 & 0x0F
	case CmdScanStop:
		// 停止扫描
		cmdCode = 0x00
	case CmdAssistStart:
		// 开启辅助
		cmdCode = 0x8c
		v1 = m.V1
	case CmdAssistStop:
		// 关闭辅助
		cmdCode = 0x8d
		v1 = m.V1
	default:
		return "", errCmd
	}
	// 校验码
	code := 0xA5 + 0x0F + 0x01 + cmdCode + int(v1) + int(v2) + int(v3)
	code = code % 256
	// 最终指令
	return fmt.Sprintf("A50F01%.2X%.2X%.2X%.2X%.2X", cmdCode, int(v1), int(v2), int(v3), code), nil
}

// SendPTZ 云台控制
func SendPTZ(ctx context.Context, m *PTZ) error {
	// 命令
	cmd, err := m.Cmd()
	if err != nil {
		return err
	}
	m.RawCmd = cmd
	return SendPTZRaw(ctx, m)
}

// SendPTZRaw 云台控制
func SendPTZRaw(ctx context.Context, m *PTZ) error {
	// 消息
	var body xml.Message
	body.XMLName.Local = xml.TypeControl
	body.CmdType = xml.CmdDeviceControl
	// 这个应该是用的通道编号
	body.DeviceID = m.ChannelID
	body.SN = sip.GetSNString()
	body.PTZCmd = m.RawCmd
	if m.ControlPriority != "" {
		body.Info = new(xml.MessageInfo)
		body.Info.ControlPriority = m.ControlPriority
	}
	// 请求
	return request.SendMessage(ctx, m.TraceID, m.Ser, m.Device, &body, nil)
}
