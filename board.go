package main

import (
	"errors"
	"fmt"
	"net/http"
)


type NjtError struct {
	Msg string `json:"msg"`
	Code int `json:"code"`
}

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
	Key string
	Is_open bool
	Is_init bool
	Serial string
	Description string
	Vendor string
	Fpgas []Fpga
	Bmc   Bmc
	synced bool
}


func (b *Board) GetBmc() error {


}


func (b *Board) init() error {

}

func (b *Board) open() error {

}

type BoardService struct {
	client *Client
	prefix String
	Boards map[string]Board

}

func NewBoardService(client *Client)  *BoardService {
	b := &BoardService {
		client: client,
		prefix: "/boards",
	}

	// initialize services here

	return b;
}

func (bs *BoardService) NewRequest(urlstr string, method string, body interface{}) (*http.Request, error) {
	return bs.client.NewRequest(bs.prefix + urlStr, method, body);
}

func (bs *BoardService) Do(req *http.Request, into interface{}) {*http.Response, error) {
	return bs.client.Do(req, into)
}

func (bs *BoardService) QueryBoards() error {
	type GetBoardsResult struct {
		Boards []BoardInfo `json:"boards"`
	}

	req, err := bs.NewRequest("", "GET", nil)
	if err != nil {
		return err
	}
	var result GetBoardsResult;
	_, err = bs.Do(req, &result)
	if err != nil {
		return err
	}
	// as board keys are unique, we should probably use a map rather than a list


	// update board infos
	// this part is quite tricky
	//case 1. board obj has not been made yet
	//				- make new object
	//case 2. board obj has already been made
	//				- update existing object
	//case 3. board obj exists but BoardInfo is not found for it
	//				- warn that the board is no longer being detected and remove it

	// for each board, set synced to False.

	// for each bi, check if object exists
		// if it does not, make a new object and start tracking it
		// else, update its data
		// set board.synced = True

	// walk through list of boards one more time, if a board exists with synced = False,
	// then that board went missing on the server. Warn the user, and untrack that board.
}

func (bs *BoardService)  BoardBySerial(Serial string) *Board {
	for _, b := range bs.Boards:
		if b.Serial == Serial {
			return b
		}
	return nil
}

func (bs *BoardService) BoardByKey(Key string) *Board {
	for _, b := range bs.Boards:
		if b.Key == Key {
			return b
		}
	return nil
}

func (bs *BoardService) untrackBoard(Key string) error {
	return nil;
}

func (bs *BoardService) trackNewBoard(Key string) error {
	return nil;
}

func (b *Board) UpdatefromInfo(bi *BoardInfo) error {

	if bi.Open_error != nil {
		return errors.new(bi.Open_error.Description)
	}

	if bi.Init_error != nil {
		return errors.new(bi.Init_error.Description)
	}

	b.Key = bi.Key
	b.Is_open = bi.Is_open
	b.Is_init = bi.Is_init
	b.Serial = bi.Serial
	b.description = bi.Description
	b.Vendor = bi.Vendor

	// update BMC info TODO
	// update Fpga info TODO
	return nil
}
