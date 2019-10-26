package njclient

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type FpgaService struct {
	client *Board
	prefix string
	Fpgas  []*Fpga
}

func NewFpgaService(client *Board) *FpgaService {
	b := &FpgaService{
		client: client,
		prefix: "fpgas",
	}
	return b
}

func (fs *FpgaService) GetFpga(index int) (*Fpga, error) {
	if index >= len(fs.Fpgas) {
		return nil, errors.New(fmt.Sprintf("could not find FPGA with index %d", index))
	} else {
		return fs.Fpgas[index], nil
	}
}

func (fs *FpgaService) NewRequest(urlstr string, method string, body interface{}) (*http.Request, error) {
	return fs.client.NewRequest(fs.prefix+urlstr, method, body)
}

func (fs *FpgaService) Do(req *http.Request, into interface{}) (*http.Response, error) {
	return fs.client.Do(req, into)
}

func (fs *FpgaService) QueryFpgas() error {
	type GetFpgasResult struct {
		Fpgas []FpgaInfo `json:"fpgas"`
	}

	req, err := fs.NewRequest("", "GET", nil)
	if err != nil {
		return err
	}
	var result GetFpgasResult
	_, err = fs.Do(req, &result)
	if err != nil {
		return err
	}

	fs.Fpgas = nil
	for _, info := range result.Fpgas {
		err = fs.trackNewFpga(&info)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fs *FpgaService) trackNewFpga(fi *FpgaInfo) error {
	f, err := NewFpga(fs)
	if err != nil {
		return err
	}
	err = f.UpdatefromInfo(fi)
	if err != nil {
		return err
	}

	index := 0
	if len(fs.Fpgas) > 0 {
		index = fs.Fpgas[len(fs.Fpgas)-1].Index
		index += 1
	}
	f.Index = index
	fs.Fpgas = append(fs.Fpgas, f)
	return nil
}

func (fs *FpgaService) FpgaByDna(dna string) *Fpga {
	for _, fpga := range fs.Fpgas {
		if fpga.Dna == dna {
			return fpga
		}
	}
	return nil
}

func (fs *FpgaService) FpgaByIndex(index int) *Fpga {
	if index >= len(fs.Fpgas) {
		return nil
	}
	f := fs.Fpgas[index]
	if f.Index != index {
		log.Fatalln("Code bug in FpgaByIndex!")
	}
	return f
}
