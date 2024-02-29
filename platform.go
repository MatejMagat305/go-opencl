package opencl

import (
	"strings"
)

type Platform uint

func GetPlatforms() ([]Platform, error) {
	numPlatforms := uint32(0)
	st := getPlatformIDs(0, nil, &numPlatforms)
	if st != CL_SUCCESS {
		return nil, st.getError()
	}

	platformIDs := make([]Platform, numPlatforms)
	st = getPlatformIDs(numPlatforms, platformIDs, nil)
	return platformIDs, st.getError()
}

type platformInfo uint

const (
	platformInfoProfile    platformInfo = 0x0900
	platformInfoVersion    platformInfo = 0x0901
	platformInfoName       platformInfo = 0x0902
	platformInfoVendor     platformInfo = 0x0903
	platformInfoExtensions platformInfo = 0x0904
)

func (p Platform) getInfo(name platformInfo) (string, error) {
	size := clSize(0)
	st := getPlatformInfo(p, name, clSize(0), nil, &size)
	if st != CL_SUCCESS {
		return "", st.getError()
	}

	info := make([]byte, size)
	st = getPlatformInfo(p, name, size, info, nil)
	if st != CL_SUCCESS {
		return "", st.getError()
	}

	return string(info), nil
}

func (p Platform) GetProfile() (string, error) {
	return p.getInfo(platformInfoProfile)
}

func (p Platform) GetVersion() (string, error) {
	return p.getInfo(platformInfoVersion)
}

func (p Platform) GetName() (string, error) {
	return p.getInfo(platformInfoName)
}

func (p Platform) GetVendor() (string, error) {
	return p.getInfo(platformInfoVendor)
}

func (p Platform) GetExtensions() ([]Extension, error) {
	extensions, err := p.getInfo(platformInfoExtensions)
	if err != nil {
		return nil, err
	}
	return strings.Split(extensions, " "), nil
}

func (p Platform) GetDevices(deviceType DeviceType) ([]Device, error) {
	numDevices := uint32(0)
	st := getDeviceIDs(p, deviceType, 0, nil, &numDevices)
	if st != CL_SUCCESS {
		return nil, st.getError()
	}

	deviceIDs := make([]Device, numDevices)
	st = getDeviceIDs(p, deviceType, numDevices, deviceIDs, nil)
	if st != CL_SUCCESS {
		return nil, st.getError()
	}

	return deviceIDs, nil
}
