package bdc

import (
	"context"
	"encoding/json"
	"github.com/bytedance/gopkg/util/logger"
	"os/exec"
)

type Bdc struct {

}

type SimulatorDevice struct {
	ProductType  string `json:"productType"`
	IosVersion   string `json:"ios_version"`
}

func (b *Bdc)GetIOSSimulatorProperties(ctx context.Context, serial string) {
	var args []string
	args = append(args, []string{
		"list-targets",
		"-u",
		serial,
	    "-o",
	    "json"}...)
	cmd := exec.CommandContext(ctx, "bdcstl-server", args...)
	output, err := cmd.Output()
	if err != nil {
		logger.CtxErrorf(ctx, "property, get device info with bdc error: %v", err)
		return
	}
	logger.Info("debug, get devices :%s", string(output))

	var SimulatorDevices []SimulatorDevice
	err = json.Unmarshal(output, &SimulatorDevices)
	if err != nil {
		logger.Error("ios_devices, Unmarshal exec return msg with error: %s", err.Error())
		return
	}
	if len(SimulatorDevices) != 1 {
		logger.Error("ios_devices, get wrong device length: %d", len(SimulatorDevices))
		return
	}
	logger.Info("succeed")
}