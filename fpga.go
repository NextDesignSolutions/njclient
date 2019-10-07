package main

import (
	"net/http"
)

type FpgaInfo struct {
	Name string `json:"name"`
	Dna string  `json:"dna"`
}

type Fpga struct {
	Name string
	Dna  string
}

type FpgaService struct {
	client *BoardService
	prefix string
	Fpgas map[string]Fpga
}

func NewFpgaService(client *BoardService) *FpgaService {
	b := &FpgaService{
		client: client,
		prefix:"/fpgas",
	}

	return b;
}

func (fs *FpgaService) NewRequest(urlstr string, method string, body interface{}) (*http.Request, error) {
	return fs.client.NewRequest(fs.prefix + urlstr, method, body);
}

func (fs *FpgaService) Do(req *http.Request, into interface{}) (*http.Response, error) {
	return fs.client.Do(req, into);
}

func (fs *FpgaService) QueryBoards() error {
	type GetFpgasResult struct {
		Fpgas []FpgaInfo `json:"fpgas"`
	}

	req, err := fs.NewRequest("", "GET", nil)
	if err != nil {
		return err
	}
	var result GetFpgasResult;
	_, err = fs.Do(req, &result)
	if err != nil {
		return err
	}
	// handle updating fpga objects

	return nil
}

func (fs *FpgaService) FpgaByDna(dna string) *Fpga {

	return nil;
}


func (fs *FpgaService) UpdatefromInfo(fi *FpgaInfo) error {
	return nil;
}
