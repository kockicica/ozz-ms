//go:build !windows
// +build !windows

package oto

func GetDeviceCount() (int, error) {
	return 0, nil
}

func GetDevices() ([]DeviceInfo, error) {
	return nil, nil
}
