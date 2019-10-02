package main

import (
	"errors"
	"fmt"
	"net/http"
)

type BmcInfo struct {
	Is_setup bool `json:"is_setup"`
	Setup_error *NjtError`json:"setup_error"`
	Close_error *NjtError 'json:"close_error"'
	Access_type string `json:"access_type"`
	Sensor_list []string `json:"sensor_list"`
}


type Bmc struct {
	Is_setup bool
	Access_type string
}

func (b *Bmc) setCanChangeBitstream(bool allowed) error {

}

func (b *Bmc) setVoltage(int fpga_index, float32 voltage) error {

}

func (b *Bmc) querySensors(string *filter) error {

}

