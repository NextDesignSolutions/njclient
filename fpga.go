package main

import (
	"net/http"
)

type FpgaInfo struct {
	Name string `json:"name"`
	Dna  string `json:"dna"`
}

type Fpga struct {
	client *FpgaService
	Index  int
	Name   string
	Dna    string
}

func (f *Fpga) UpdatefromInfo(fi *FpgaInfo) error {
	f.Name = fi.Name
	f.Dna = fi.Dna
	return nil
}

func NewFpga(service *FpgaService) (*Fpga, error) {
	f := Fpga{
		client: service,
	}
	return &f, nil
}

func (f *Fpga) NewRequest(urlstr string, method string, body interface{}) (*http.Request, error) {
	return f.client.NewRequest("/"+string(f.Index)+"/"+urlstr, method, body)
}

func (f *Fpga) Do(req *http.Request, into interface{}) (*http.Response, error) {
	return f.client.Do(req, into)
}
