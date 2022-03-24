package adb

type DeviceInfo struct {
	Ip     string
	Port   uint
	Serial string
	Cmd    string
}

//func PwdDevice(d *DeviceInfo) {
//	adbClient := adb.New(adb.Options{
//		Network: "tcp",
//		Address: fmt.Sprintf("%s:%d", d.Ip, d.Port),
//	})
//	output, err := adbClient.ShellWithOutput(context.Background(), d.Serial, d.Cmd)
//	if err != nil {
//		logger.Error("list devices with error: %s", err.Error())
//	}
//	logger.Info("device test, cmd: %s, result: %s", d.Cmd, string(output))
//}
