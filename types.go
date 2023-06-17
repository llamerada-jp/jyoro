package main

type Config struct {
	Devices []DeviceConfig `json:"devices"`
}

type DeviceConfig struct {
	Path     string                `json:"path,omitempty"`
	Channels map[int]ChannelConfig `json:"channels"`
}

type ChannelConfig struct {
	Name *string `json:"name,omitempty"`
}
