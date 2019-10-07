package main

type BmcInfo struct {
	Is_setup    bool      `json:"is_setup"`
	Setup_error *NjtError `json:"setup_error"`
	Close_error *NjtError `json:"close_error"`
	Access_type string    `json:"access_type"`
	Sensor_list []string  `json:"sensor_list"`
}

type Bmc struct {
	Is_setup    bool
	Access_type string
}

func (b *Bmc) setCanChangeBitstream(allowed bool) error {
	return nil
}

func (b *Bmc) setVoltage(fpga_index int, voltage float32) error {
	return nil
}

func (b *Bmc) querySensors(filter *string) error {
	return nil
}
