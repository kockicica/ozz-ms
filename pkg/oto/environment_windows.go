package oto

import (
	"unicode/utf16"
)

func GetDeviceCount() (int, error) {
	var devices int
	if err := waveOutGetNumDevs(&devices); err != nil {
		return 0, err
	}
	return devices, nil
}

func GetDevices() ([]DeviceInfo, error) {

	deviceCount, err := GetDeviceCount()
	if err != nil {
		return nil, err
	}

	var res []DeviceInfo

	for i := 0; i < deviceCount; i++ {
		caps := new(waveoutcaps)
		err := waveOutGetDevCaps(i, caps)
		if err != nil {
			return nil, err
		}
		var li []uint16
		for i := 0; i < len(caps.szPname); i++ {
			li = append(li, caps.szPname[i])
		}
		runes := utf16.Decode((li))
		res = append(res, DeviceInfo{
			Id:   i,
			Name: string(runes),
		})
	}

	return res, nil
}
