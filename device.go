package hass

import (
	"errors"
	"strings"
)

// Device is a generic interface for interacting with devices
type Device interface {
	On() error
	Off() error
	Toggle() error
	EntityID() string
	Domain() string
}

// GetDevice returns a Device object from an State object
func (a *Access) GetDevice(state State) (Device, error) {
	dom := strings.TrimSuffix(strings.SplitAfter(state.EntityID, ".")[0], ".")
	switch dom {
	case "light":
		return a.NewLight(state.EntityID), nil
	case "switch":
		return a.NewSwitch(state.EntityID), nil
	case "lock":
		return a.NewLock(state.EntityID), nil
	}
	return nil, errors.New("Device type not supported yet")
}

// SupportedDeviceTypes returns a list of supported device types
func (a *Access) SupportedDeviceTypes() []string {
	return []string{"light", "switch", "lock"}
}

// IsSupportedDevice returns true if an entityID is a supported device
func (a *Access) IsSupportedDevice(id string) bool {
	dom := strings.TrimSuffix(strings.SplitAfter(id, ".")[0], ".")
	for _, d := range a.SupportedDeviceTypes() {
		if dom == d {
			return true
		}
	}
	return false
}
