package main

import (
	"net/http"
	"log"

)

type BoardService struct {
	client *Client
	prefix string
	Boards map[string]*Board

}

func NewBoardService(client *Client)  *BoardService {
	b := &BoardService {
		client: client,
		prefix: "boards",
	}
	b.Boards = make(map[string]*Board)

	// initialize services here

	return b;
}

func (bs *BoardService) NewRequest(urlstr string, method string, body interface{}) (*http.Request, error) {
	log.Printf("bs.prefix = %s\n", bs.prefix)
	return bs.client.NewRequest(bs.prefix + urlstr, method, body);
}

func (bs *BoardService) Do(req *http.Request, into interface{}) (*http.Response, error) {
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

	// update all boards that are already known
	for key, value := range bs.Boards {
		for i := len(result.Boards)-1; i >=0; i-- {
			bi := result.Boards[i]
			if key == bi.Key {
				err = value.updateFromInfo(&bi)
				if err != nil {
					return err
				}
				// remove board info 
				result.Boards = append(result.Boards[:i], result.Boards[i+1:]...)
				break
			}
		}

		// if wee get here, then a board has been dropped
		// FIXME: warn user
		err = bs.untrackBoard(key)
		if err != nil {
			return err
		}
	}
	// board infos left in list are new and previously untracked
	for _, bi := range result.Boards {
		err := bs.trackNewBoard(&bi)
		if err != nil {
			return err;
		}
	}
	return nil
}

func (bs *BoardService)  BoardBySerial(serial string) *Board {
	for _, b := range bs.Boards {
		if b.Serial == serial {
			return b
		}
	}
	return nil
}

func (bs *BoardService) BoardByKey(key string) *Board {
	i, ok := bs.Boards[key]
	if ok == true {
		return i;
	}
	return nil;
}

func (bs *BoardService) untrackBoard(key string) error {
	delete(bs.Boards, key)
	return nil;
}

func (bs *BoardService) trackNewBoard(bi *BoardInfo) error {
	n, err := NewBoard(bs)
	if err != nil {
		return err
	}
	err = n.updateFromInfo(bi)
	if err != nil {
		return err
	}
	bs.Boards[bi.Key] = n
	return nil
}