package main

import (
	"net/http"
	"errors"
	"log"

)

type BoardInfo struct {
	Key         string `json:"key"`
	Is_open     bool `json:"is_open"`
	Open_error  *NjtError `json:"open_error"`
	Is_init     bool `json:"is_init"`
	Init_error  *NjtError `json:"init_error"`
	Serial      string `json:"serial"`
	Description string `json:"description"`
	Vendor      string `json:"vendor"`
	Bmc         *BmcInfo `json:"bmc"`
	Fpgas       []FpgaInfo `json:"fpgas"`
}

type Board struct {
	client *BoardService
	Key string
	Is_open bool
	Is_init bool
	Serial string
	Description string
	Vendor string
	Fpgas []Fpga
	FpgaService *FpgaService
	synced bool
}

func NewBoard(client *BoardService) (*Board, error) {
	n :=  Board {}
	err := n.initFpgaService()
	if err != nil {
		return  nil, err
	 }
	return &n, nil;
}


func (b *Board) GetBmc() error {
	return nil
}

func (b *Board) init() error {
	var result BoardInfo;
	req, err := b.NewRequest("init", "POST", nil)
	if err != nil {
		return err
	}
	_, err = b.Do(req, &result)
	if err != nil {
		return err
	}
	err = b.updateFromInfo(&result)
	return err
}

func (b *Board) open() error {
	var result BoardInfo;
	req, err := b.NewRequest("open", "POST", nil)
	if err != nil {
		return err
	}
	_, err = b.Do(req, &result)
	if err != nil {
		return err
	}
	err = b.updateFromInfo(&result)
	return err
}

func (b *Board)  initFpgaService() error {
	fs := NewFpgaService(b)
	if fs == nil {
		return errors.New("failed to init FpgaService for board " + b.Key)
	}
	b.FpgaService = fs
	return nil
}

func (b *Board) updateFromInfo(bi *BoardInfo) error {
	b.Key = bi.Key
	b.Is_open = bi.Is_open
	b.Is_init = bi.Is_init
	b.Serial = bi.Serial
	b.Description = bi.Description
	b.Vendor = bi.Vendor

	if bi.Open_error != nil {
		return errors.New(bi.Open_error.Msg)
	}

	if bi.Init_error != nil {
		return errors.New(bi.Init_error.Msg)
	}

	// update BMC info TODO
	// update Fpga info TODO
	return nil
}

func (b *Board) NewRequest(urlstr string, method string, body interface{}) (*http.Request, error) {
	log.Printf("b.Key = %s\n", b.Key)
	return b.client.NewRequest(b.Key + urlstr, method, body);
}

func (b *Board) Do(req *http.Request, into interface{}) (*http.Response, error) {
	return b.client.Do(req, into)
}
