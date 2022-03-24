package usb

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
	"os/exec"
	"regexp"
	"strings"
	"test/serial"
	"time"
)

const (
	//vendor:0424.product:2807
	HubVendor  = 1060
	HubProduct = 10247

	IosVendor = 0x05AC
)

type OperationType string

const (
	PlugDevice     OperationType = "plugDevice"
	SetQuickCharge OperationType = "setQuickCharge"
	SetSlowCharge  OperationType = "setSlowCharge"
)

type HubController struct {
	OperationType OperationType
	UsbCtx        *gousb.Context
	Uuid          string
}

func (h *HubController) HandleControlDevice() error {
	devs, err := h.UsbCtx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return true
	})
	if err != nil {
		logger.Error("gousb, ctx open devices with error: ", err.Error())
		return err
	}

	isDeviceExist, devicePortPath := h.iSDeviceConnectedWithMachine(devs)
	if !isDeviceExist || len(devicePortPath) == 0 {
		logger.Warn("gousb, uuid not exist in gousb-devices", h.Uuid)
		return err
	}
	logger.Info("gousb, find given device [%s]", h.Uuid)

	isDeviceConnectedToHub, deviceHub := h.iSDeviceConnectedWithHub(devs, devicePortPath)
	if !isDeviceConnectedToHub || deviceHub == nil {
		logger.Warn("device not connected to hub")
		return err
	}
	logger.Info("gousb, find device hub: ", deviceHub.String())

	isFindHubSerial, hubSerialPath := h.GetHubSerial(deviceHub)
	if !isFindHubSerial || hubSerialPath == "" {
		logger.Warn("gousb, find no serial path for device: ", h.Uuid)
		return err
	}
	logger.Info("serial, find device hub serial: ", hubSerialPath)

	err = h.execControlDevice(hubSerialPath, devicePortPath[1])
	if err != nil {
		logger.Error("serial, control device with error: ", err.Error())
		return err
	}
	logger.Info("hubplug, succeed to plug device: ", h.Uuid)
	return nil
}

func (h *HubController) iSDeviceConnectedWithMachine(devs []*gousb.Device) (bool, []int) {
	for _, dev := range devs {
		if s, err := dev.SerialNumber(); err != nil {
			logger.Warn("get device serial with error: %s", err.Error())
			continue
		} else if s == "" {
			logger.Warn("get null serial")
			continue
		} else {
			uniformUuid := h.getDeviceUniformUuid(s, int(dev.Desc.Vendor))
			if uniformUuid == h.Uuid {
				return true, dev.Desc.Path
			}
		}
	}
	return false, nil
}

func (h *HubController) iSDeviceConnectedWithHub(devs []*gousb.Device, devicePortPath []int) (bool, *gousb.Device) {
	for _, hub := range devs {
		if hub.Desc.Vendor != gousb.ID(HubVendor) || hub.Desc.Product != gousb.ID(HubProduct) {
			continue
		}
		logger.Info("gousb, get bytedance hub: ", hub.String())
		isConnected := true
		for index, port := range hub.Desc.Path {
			if devicePortPath[index] != port {
				isConnected = false
				break
			}
		}
		if isConnected {
			return true, hub
		} else {
			continue
		}
	}
	return false, nil
}

func (h *HubController) GetHubSerial(hub *gousb.Device) (bool, string) {
	s := serial.HubSerial{}
	hubSerialPaths := s.GetSerialPathsByName(serial.HubLinuxPathPrefix)
	for _, path := range hubSerialPaths {
		usbInfo := h.GetSerialUsbInfo(path)
		if strings.Contains(usbInfo, fmt.Sprintf("ATTRS{busnum}==\"%d\"", hub.Desc.Bus)) &&
			strings.Contains(usbInfo, fmt.Sprintf("ATTRS{devnum}==\"%d\"", hub.Desc.Address)) {
			return true, path
		}
	}
	return false, ""
}

func (h *HubController) getDeviceUniformUuid(uuid string, vendor int) string {
	serialProcess := strings.Split(uuid, "\x00")
	if vendor == IosVendor {
		reg, _ := regexp.Compile("[A-Z]")
		if reg.MatchString(serialProcess[0]) {
			return serialProcess[0][:8] + "-" + serialProcess[0][8:]
		} else {
			return serialProcess[0]
		}
	} else {
		return serialProcess[0]
	}
}

func (h *HubController) GetSerialUsbInfo(path string) string {
	cmd := fmt.Sprintf("udevadm info -a  -n %s | grep -E  \"busnum|devnum\" | awk '{print $1}'", path)
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	if info, err := exec.CommandContext(ctx, "/bin/bash", "-c", cmd).Output(); err != nil {
		logger.Error("agent/check, get driver cmd with error: %s", err.Error())
		return string(info)
	} else {
		return string(info)
	}
}

func (h *HubController) getDevicesInfo(devs []*gousb.Device) {
	for _, dev := range devs {
		usbUuid, _ := dev.SerialNumber()
		logger.Info("device info: ",
			"usbUuid: ", usbUuid,
			"path: ", dev.Desc.Path,
			"vendor: ", dev.Desc.Vendor,
			"product: ", dev.Desc.Product,
			dev.Desc.String(),
			"desc: ", usbid.Describe(dev.Desc),
		)
	}
}

func (h *HubController) execControlDevice(hubSerialPath string, devicePort int) error {
	if h.OperationType == PlugDevice {
		return h.PlugDevice(hubSerialPath, devicePort)
	} else if h.OperationType == SetQuickCharge {
		return h.SetDeviceQuickCharge(hubSerialPath, devicePort)
	} else if h.OperationType == SetSlowCharge {
		return h.SetDeviceSlowCharge(hubSerialPath, devicePort)
	} else {
		return errors.New("undefined control type")
	}
}

func (h *HubController) PlugDevice(hubSerialPath string, devicePort int) error {
	hub := serial.HubSerial{}
	originCmd := []byte{0xF3, 0x02, 0xFF}
	portOrigin := 0x01
	controlPort := portOrigin << (devicePort - 1)
	port := 0xFF ^ controlPort
	originCmd = append(originCmd, byte(port))
	cmdCrc := hub.GenerateCrc16(originCmd, len(originCmd))
	originCmd = append(originCmd, cmdCrc...)
	logger.Info("control cmd: ", originCmd)
	fmt.Printf("cmd with crc: %X \n", originCmd)
	err := hub.CheckAndSendCmdToHub(hubSerialPath, originCmd)
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 100)
	originCmd = []byte{0xF3, 0x02, 0xFF, 0xFF}
	cmdCrc = hub.GenerateCrc16(originCmd, len(originCmd))
	originCmd = append(originCmd, cmdCrc...)
	logger.Info("control cmd: ", originCmd)
	fmt.Printf("cmd with crc: %X \n", originCmd)
	err = hub.CheckAndSendCmdToHub(hubSerialPath, originCmd)
	if err != nil {
		return err
	}
	return nil
}

func (h *HubController) SetDeviceQuickCharge(hubSerialPath string, devicePort int) error {
	hub := serial.HubSerial{}
	originCmd := []byte{0xF3, 0x06, 0x00}
	portOrigin := 0x01
	controlPort := portOrigin << (devicePort - 1)
	port := 0x00 ^ controlPort
	originCmd = append(originCmd, byte(port))
	cmdCrc := hub.GenerateCrc16(originCmd, len(originCmd))
	originCmd = append(originCmd, cmdCrc...)
	logger.Info("control cmd:  %X", originCmd)
	err := hub.CheckAndSendCmdToHub(hubSerialPath, originCmd)
	if err != nil {
		return err
	}
	return nil
}

func (h *HubController) SetDeviceSlowCharge(hubSerialPath string, devicePort int) error {
	hub := serial.HubSerial{}
	originCmd := []byte{0xF3, 0x04, 0x00}
	portOrigin := 0x01
	controlPort := portOrigin << (devicePort - 1)
	port := 0x00 ^ controlPort
	originCmd = append(originCmd, byte(port))
	cmdCrc := hub.GenerateCrc16(originCmd, len(originCmd))
	originCmd = append(originCmd, cmdCrc...)
	logger.Info("control cmd:  %X", originCmd)
	err := hub.CheckAndSendCmdToHub(hubSerialPath, originCmd)
	if err != nil {
		return err
	}
	return nil
}
